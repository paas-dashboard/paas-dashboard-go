package checker

import (
	"fmt"
	"github.com/beego/beego/v2/adapter/toolbox"
	"github.com/protocol-laboratory/pulsar-admin-go/padmin"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/rand"
	"os"
	"paas-dashboard-go/module"
	"strconv"
	"strings"
	"time"
)

func pulsarConsumerAbnormalCheck(instance *module.PulsarInstance) {
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
	var admins []*padmin.PulsarAdmin
	for _, host := range hosts {
		admin, err := padmin.NewPulsarAdmin(padmin.Config{
			Host: host,
			Port: instance.WebPort,
		})
		if err != nil {
			panic(err)
		}
		admins = append(admins, admin)
	}
	interval, _ := strconv.Atoi(os.Getenv("PD_PULSAR_CONSUMER_ABNORMAL_CHECK_INTERVAL"))
	if interval == 0 || interval < 60 {
		interval = 300
	}

	var status string

	min := interval / 60

	cronExpress := fmt.Sprintf("0 */%d * * * *", min)

	tk := toolbox.NewTask("abnormalCheck", cronExpress, func() error {
		if status == "start" {
			return nil
		}
		logrus.Infof("pulsar instance:[%s] consumer abnormal check start.", instance.Name)
		startTime := time.Now().Unix()
		status = "start"
		totalCount, abnormalCount := checkConsumerAbnormalInterval(admins)
		logrus.Infof("pulsar instance:[%s] consumer abnormal check done. spend %ds. scan result:%d|%d",
			instance.Name, time.Now().Unix()-startTime, totalCount, abnormalCount)
		status = "done"
		return nil
	})
	err = tk.Run()
	if err != nil {
		panic(err)
	}
	toolbox.AddTask("abnormalCheck", tk)
	toolbox.StartTask()
}

func checkConsumerAbnormalInterval(admins []*padmin.PulsarAdmin) (int, int) {
	var (
		totalCount    int
		abnormalCount int
	)
	admin := admins[rand.IntnRange(0, len(admins))]
	if admin == nil {
		return 0, 0
	}
	tenantList, err := admin.Tenants.List()
	if err != nil {
		logrus.Errorf("check consumer abnormal|get tenant list failed: %v", err)
		return totalCount, abnormalCount
	}
	for _, tenant := range tenantList {
		namespaceList, err := admin.Namespaces.List(tenant)
		if err != nil {
			logrus.Errorf("check consumer abnormal|get tenant %s namespace failed: %v", tenant, err)
			continue
		}
		for _, namespace := range namespaceList {
			ret := strings.Split(namespace, "/")
			ns := ret[len(ret)-1]

			topicList, err := admin.PersistentTopics.ListNamespaceTopics(tenant, ns)
			if err != nil {
				logrus.Errorf("check consumer abnormal|get tenant %s namespace %s topic list fail: %v", tenant, ns, err)
				continue
			}

			for _, v := range topicList {
				ret := strings.Split(v, "/")
				topic := ret[len(ret)-1]

				topicStatistics, err := admin.PersistentTopics.GetStats(tenant, ns, topic)
				if err != nil {
					logrus.Errorf("check consumer abnormal| get stats failed:%v", err)
					continue
				}
				logrus.Debugf("check consumer abnormal|%s| scan tpoic start ... ", v)
				abNormal := false
				for _, sub := range topicStatistics.Subscriptions {
					if sub.MsgBacklog <= 0 {
						continue
					}
					dev := time.Now().UnixMilli() - int64(sub.LastConsumedFlowTimestamp)
					if dev > int64(5*60*1000*time.Millisecond) {
						abNormal = true
						logrus.Warnf("check consumer abnormal|%s|%+v it seems to be broken consumer!", v, sub)
					}
				}
				if abNormal {
					abnormalCount++
				}
				totalCount++
				logrus.Debugf("check consumer abnormal|%s| scan tpoic done. ", v)
			}
		}
	}
	return totalCount, abnormalCount
}
