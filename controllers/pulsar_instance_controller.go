package controllers

import (
	"os"
	"paas-dashboard-go/module"
	"strconv"
	"strings"
)

var PulsarInstanceMap = make(map[string]*module.PulsarInstance)

func init() {
	prefix := "PD_PULSAR_"
	prefixLen := len(prefix)

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		key, value := pair[0], pair[1]

		if strings.HasPrefix(key, prefix) {
			key = key[prefixLen:]
			index := strings.Index(key, "_")
			name := strings.ToLower(key[:index])
			confProperty := key[index+1:]

			if _, ok := PulsarInstanceMap[name]; !ok {
				PulsarInstanceMap[name] = &module.PulsarInstance{Name: name}
			}
			pulsarInstance := PulsarInstanceMap[name]

			switch confProperty {
			case "HOST":
				pulsarInstance.Host = value
			case "WEB_PORT":
				pulsarInstance.WebPort, _ = strconv.Atoi(value)
			case "TCP_PORT":
				pulsarInstance.TcpPort, _ = strconv.Atoi(value)
			}
		}
	}
}

type PulsarInstancesController struct {
	PulsarController
}

func (p *PulsarInstancesController) Get() {
	instances := make([]*module.PulsarInstance, 0, len(PulsarInstanceMap))
	for _, instance := range PulsarInstanceMap {
		instances = append(instances, instance)
	}
	p.Data["json"] = instances
	_ = p.ServeJSON()
}
