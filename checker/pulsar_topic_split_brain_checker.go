package checker

import (
	"fmt"
	"github.com/beego/beego/v2/adapter/toolbox"
	"github.com/protocol-laboratory/pulsar-admin-go/padmin"
	"github.com/sirupsen/logrus"
	"os"
	"paas-dashboard-go/module"
	"strconv"
	"strings"
	"time"
)

func pulsarTopicSplitBrainCheck(instance *module.PulsarInstance) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("panic:%v", err)
			return
		}
	}()
	hosts, err := GetHosts()
	if err != nil {
		panic(err)
	}
	logrus.Infof("check hosts list:%v", hosts)
	var srcAdmin *padmin.PulsarAdmin
	var destAdmin []*padmin.PulsarAdmin
	for _, host := range hosts {
		admin, err := padmin.NewPulsarAdmin(padmin.Config{
			Host: host,
			Port: instance.WebPort,
		})
		if err != nil {
			panic(err)
		}
		if srcAdmin == nil {
			srcAdmin = admin
		} else {
			destAdmin = append(destAdmin, admin)
		}
	}

	if srcAdmin == nil || len(destAdmin) == 0 {
		return
	}

	interval, _ := strconv.Atoi(os.Getenv("PD_PULSAR_TOPIC_SPLIT_BRAIN_CHECK_INTERVAL"))
	if interval == 0 || interval < 60 {
		interval = 300
	}

	var status string

	min := interval / 60

	cronExpress := fmt.Sprintf("0 */%d * * * *", min)

	tk := toolbox.NewTask("splitBrainCheck", cronExpress, func() error {
		if status == "start" {
			return nil
		}
		startTime := time.Now().Unix()
		logrus.Infof("pulsar instance:[%s] topic split brain check start.", instance.Name)
		status = "start"
		totalCount, splitCount := checkInterval(srcAdmin, destAdmin)
		logrus.Infof("pulsar instance:[%s] topic split brain check done. spend %ds. scan result:%d|%d",
			instance.Name, time.Now().Unix()-startTime, totalCount, splitCount)
		status = "done"
		return nil
	})
	err = tk.Run()
	if err != nil {
		panic(err)
	}
	toolbox.AddTask("splitBrainCheck", tk)
	toolbox.StartTask()
}

func checkInterval(srcAdmin *padmin.PulsarAdmin, destAdmin []*padmin.PulsarAdmin) (int, int) {
	var (
		totalCount int
		splitCount int
	)
	tenantList, err := srcAdmin.Tenants.List()
	if err != nil {
		logrus.Errorf("check split brain| src get tenant list failed: %v", err)
		return totalCount, splitCount
	}
	for _, tenant := range tenantList {
		namespaceList, err := srcAdmin.Namespaces.List(tenant)
		if err != nil {
			logrus.Errorf("check split brain|tenant:%s| src get namespace failed: %v", tenant, err)
			return totalCount, splitCount
		}
		for _, namespace := range namespaceList {
			ret := strings.Split(namespace, "/")
			ns := ret[len(ret)-1]

			topicList, err := srcAdmin.PersistentTopics.ListNamespaceTopics(tenant, ns)
			if err != nil {
				logrus.Errorf("check split brain|tenant:%s|namespace:%s| src get topic list fail: %v", tenant, ns, err)
				return totalCount, splitCount
			}
			for _, v := range topicList {
				ret := strings.Split(v, "/")
				topic := ret[len(ret)-1]

				srcData, err := srcAdmin.Lookup.GetOwner(padmin.TopicDomainPersistent, tenant, ns, topic)
				if err != nil {
					logrus.Errorf("check split brain|tenant:%s|namespace:%s|topic:%s| get owner fail: %v", tenant, ns, topic, err)
					continue
				}
				logrus.Debugf("check split brain|%s| scan tpoic start ... ", v)
				var destData *padmin.LookupData
				for _, admin := range destAdmin {
					destData, err = admin.Lookup.GetOwner(padmin.TopicDomainPersistent, tenant, ns, topic)
					if err != nil {
						logrus.Errorf("check split brain|tenant:%s|namespace:%s|topic:%s| get owner fail: %v", tenant, ns, topic, err)
						continue
					}
					if !eq(srcData, destData, v) {
						splitCount++
						break
					}
				}
				totalCount++
				logrus.Debugf("check split brain|%s| scan tpoic done. ", v)
			}
		}
	}
	return totalCount, splitCount
}

func eq(src, dest *padmin.LookupData, topic string) bool {
	if src.BrokerUrl != dest.BrokerUrl {
		logrus.Errorf("%s|srcBrokerUrl:%s|destBrokerUrl:%s| pulsar topic split brain.", topic, src.BrokerUrl, dest.BrokerUrl)
		return false
	}
	if src.HttpUrl != dest.HttpUrl {
		logrus.Errorf("%s|srcHttpUrl:%s|destHttpUrl:%s| pulsar topic split brain.", topic, src.HttpUrl, dest.HttpUrl)
		return false
	}
	if src.NativeUrl != dest.NativeUrl {
		logrus.Errorf("%s|srcNativeUrl:%s|destNativeUrl:%s| pulsar topic split brain.", topic, src.NativeUrl, dest.NativeUrl)
		return false
	}
	if src.BrokerUrlTls != dest.BrokerUrlTls {
		logrus.Errorf("%s|srcBrokerUrlTls:%s|destBrokerUrlTls:%s| pulsar topic split brain.", topic, src.BrokerUrlTls, dest.BrokerUrlTls)
		return false
	}
	if src.HttpUrlTls != dest.HttpUrlTls {
		logrus.Errorf("%s|srcHttpUrlTls:%s|destHttpUrlTls:%s| pulsar topic split brain.", topic, src.HttpUrlTls, dest.HttpUrlTls)
		return false
	}
	return true
}
