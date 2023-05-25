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

func pulsarConsumerAbnormalCheck(instance *module.PulsarInstance) {
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
	if interval == 0 {
		interval = 300
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for range ticker.C {
		logs.Info("pulsar instance:[%s] consumer abnormal check start.", instance.Name)
		checkConsumerAbnormalInterval(admins)
		logs.Info("pulsar instance:[%s] consumer abnormal check done.", instance.Name)
	}
}

func checkConsumerAbnormalInterval(admins []*padmin.PulsarAdmin) {
	for _, admin := range admins {
		tenantList, err := admin.Tenants.List()
		if err != nil {
			logs.Error("[check consumer abnormal]get tenant list failed: %v", err)
			return
		}
		for _, tenant := range tenantList {
			namespaceList, err := admin.Namespaces.List(tenant)
			if err != nil {
				logs.Error("[check consumer abnormal]get tenant %s namespace failed: %v", tenant, err)
				continue
			}
			for _, namespace := range namespaceList {
				ret := strings.Split(namespace, "/")
				ns := ret[len(ret)-1]

				topicList, err := admin.PersistentTopics.ListNamespaceTopics(tenant, ns)
				if err != nil {
					logs.Error("[check consumer abnormal]get tenant %s namespace %s topic list fail: %v", tenant, ns, err)
					continue
				}

				for _, v := range topicList {
					ret := strings.Split(v, "/")
					topic := ret[len(ret)-1]

					topicStatistics, err := admin.PersistentTopics.GetStats(tenant, ns, topic)
					if err != nil {
						logs.Error("[check consumer abnormal]get stats failed:%v", err)
						continue
					}
					for _, sub := range topicStatistics.Subscriptions {
						if sub.MsgBacklog <= 0 {
							continue
						}
						sub := time.Now().UnixMilli() - int64(sub.LastConsumedFlowTimestamp)
						if sub > int64(5*60*1000*time.Millisecond) {
							logs.Warn("[check consumer abnormal] consumer abnormal! send flow %d ago,it seems to be broken consumer!")
						}
					}
				}
			}
		}
	}
}
