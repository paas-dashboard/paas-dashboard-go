package checker

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/protocol-laboratory/pulsar-admin-go/padmin"
	"os"
	"paas-dashboard-go/module"
	"strconv"
	"strings"
	"time"
)

func pulsarTopicSplitBrainCheck(instance *module.PulsarInstance) {
	defer func() {
		if err := recover(); err != nil {
			logs.Error("panic:", err)
			return
		}
	}()
	hosts, err := GetHosts()
	if err != nil {
		panic(err)
	}
	logs.Info("check hosts list:%v", hosts)
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
	if interval == 0 {
		interval = 300
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for range ticker.C {
		logs.Info("pulsar instance:[%s] topic split brain check start.", instance.Name)
		checkInterval(srcAdmin, destAdmin)
		logs.Info("pulsar instance:[%s] topic split brain check done.", instance.Name)
	}
}

func checkInterval(srcAdmin *padmin.PulsarAdmin, destAdmin []*padmin.PulsarAdmin) {
	tenantList, err := srcAdmin.Tenants.List()
	if err != nil {
		logs.Error("[check split brain]src get tenant list failed: %v", err)
		return
	}
	for _, tenant := range tenantList {
		namespaceList, err := srcAdmin.Namespaces.List(tenant)
		if err != nil {
			logs.Error("[check split brain]src get tenant %s namespace failed: %v", tenant, err)
			return
		}
		for _, namespace := range namespaceList {
			ret := strings.Split(namespace, "/")
			ns := ret[len(ret)-1]

			topicList, err := srcAdmin.PersistentTopics.ListNamespaceTopics(tenant, ns)
			if err != nil {
				logs.Error("[check split brain]src get tenant %s namespace %s topic list fail: %v", tenant, ns, err)
				return
			}
			for _, v := range topicList {
				ret := strings.Split(v, "/")
				topic := ret[len(ret)-1]

				srcData, err := srcAdmin.Lookup.GetOwner(padmin.TopicDomainPersistent, tenant, ns, topic)
				if err != nil {
					logs.Error("[check split brain][tenant:%s namespace:%s topic:%s] get owner fail: %v", tenant, ns, topic, err)
					continue
				}
				var destData *padmin.LookupData
				for _, admin := range destAdmin {
					destData, err = admin.Lookup.GetOwner(padmin.TopicDomainPersistent, tenant, ns, topic)
					if err != nil {
						logs.Error("[check split brain][tenant:%s namespace:%s topic:%s] get owner fail: %v", tenant, ns, topic, err)
						continue
					}
				}
				eq(srcData, destData, v)
			}
		}
	}
}

func eq(src, dest *padmin.LookupData, topic string) bool {
	if src.BrokerUrl != dest.BrokerUrl {
		logs.Error("BrokerUrl [%s] pulsar topic split brain. src topic owner: %s", topic, src.BrokerUrl)
		logs.Error("BrokerUrl [%s] pulsar topic split brain. dest topic owner: %s", topic, dest.BrokerUrl)
		return false
	}
	if src.HttpUrl != dest.HttpUrl {
		logs.Error("HttpUrl [%s] pulsar topic split brain. src topic owner: %s", topic, src.HttpUrl)
		logs.Error("HttpUrl [%s] pulsar topic split brain. dest topic owner: %s", topic, dest.HttpUrl)
		return false
	}
	if src.NativeUrl != dest.NativeUrl {
		logs.Error("NativeUrl [%s] pulsar topic split brain. src topic owner: %s", topic, src.NativeUrl)
		logs.Error("NativeUrl [%s] pulsar topic split brain. dest topic owner: %s", topic, dest.NativeUrl)
		return false
	}
	if src.BrokerUrlTls != dest.BrokerUrlTls {
		logs.Error("BrokerUrlTls [%s] pulsar topic split brain. src topic owner: %s", topic, src.BrokerUrlTls)
		logs.Error("BrokerUrlTls [%s] pulsar topic split brain. dest topic owner: %s", topic, dest.BrokerUrlTls)
		return false
	}
	if src.HttpUrlTls != dest.HttpUrlTls {
		logs.Error("HttpUrlTls [%s] pulsar topic split brain. src topic owner: %s", topic, src.HttpUrlTls)
		logs.Error("HttpUrlTls [%s] pulsar topic split brain. dest topic owner: %s", topic, dest.HttpUrlTls)
		return false
	}
	return true
}
