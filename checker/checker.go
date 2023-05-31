package checker

import (
	"github.com/go-zookeeper/zk"
	"github.com/sirupsen/logrus"
	"net/http"
	_ "net/http/pprof"
	"os"
	"paas-dashboard-go/controllers"
	"strings"
	"time"
)

func Start() {
	go func() {
		http.ListenAndServe(":6060", nil)
	}()
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
		logrus.Errorf("connect zk failed:%v", err)
		return []string{}, err
	}
	data, _, err := conn.Children("/loadbalance/brokers")
	if err != nil {
		logrus.Errorf("get zk data failed:%v", err)
		return []string{}, err
	}
	if len(data) == 0 {
		logrus.Errorf("get zk data is none.")
		return []string{}, err
	}
	var hosts []string
	for _, v := range data {
		hosts = append(hosts, strings.Split(v, ":")[0])
	}
	return hosts, nil
}
