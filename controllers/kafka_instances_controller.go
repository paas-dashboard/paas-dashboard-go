package controllers

import (
	"os"
	"paas-dashboard-go/module"
	"strings"
)

var kafkaInstanceMap = make(map[string]*module.KafkaInstance)

func init() {
	prefix := "PD_KAFKA_"
	prefixLen := len(prefix)

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		key, value := pair[0], pair[1]

		if strings.HasPrefix(key, prefix) {
			key = key[prefixLen:]
			index := strings.Index(key, "_")
			name := strings.ToLower(key[:index])
			confProperty := key[index+1:]

			if _, ok := kafkaInstanceMap[name]; !ok {
				kafkaInstanceMap[name] = &module.KafkaInstance{Name: name}
			}
			kafkaInstance := kafkaInstanceMap[name]

			if confProperty == "BOOTSTRAP_SERVERS" {
				kafkaInstance.BootstrapServers = value
			}
		}
	}
}

type KafkaInstancesController struct {
	KafkaController
}

func (k *KafkaInstancesController) Get() {
	instances := make([]*module.KafkaInstance, 0, len(kafkaInstanceMap))
	for _, instance := range kafkaInstanceMap {
		instances = append(instances, instance)
	}
	k.Data["json"] = instances
	_ = k.ServeJSON()
}
