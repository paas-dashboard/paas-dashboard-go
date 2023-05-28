package checker

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/go-zookeeper/zk"
	"os"
	"paas-dashboard-go/controllers"
	"strings"
	"time"
)

func Start() {
	if os.Getenv("PD_PULSAR_CONSUMER_ABNORMAL_CHECK_ENABLE") == "true" {
		for _, instance := range controllers.PulsarInstanceMap {
			if instance.Host == "" {
				continue
			}
			go pulsarConsumerAbnormalCheck(instance)
		}
	}
	if os.Getenv("PD_PULSAR_TOPIC_SPLIT_BRAIN_CHECK_ENABLE") == "true" {
		for _, instance := range controllers.PulsarInstanceMap {
			if instance.Host == "" {
				continue
			}
			go pulsarTopicSplitBrainCheck(instance)
		}
	}
}

func GetHosts() ([]string, error) {
	zkServer := strings.Split(os.Getenv("ZOOKEEPER_SERVICE"), ",")
	conn, _, err := zk.Connect(zkServer, time.Second)
	if err != nil {
		logs.Error("connect zk failed:%v", err)
		return []string{}, err
	}
	data, _, err := conn.Children("/loadbalance/brokers")
	if err != nil {
		logs.Error("get zk data failed:%v", err)
		return []string{}, err
	}
	if len(data) == 0 {
		logs.Error("get zk data is none.")
		return []string{}, err
	}
	var hosts []string
	for _, v := range data {
		hosts = append(hosts, strings.Split(v, ":")[0])
	}
	return hosts, nil
}
