package controllers

import (
	"os"
	"paas-dashboard-go/module"
	"strings"
)

var redisInstanceMap = make(map[string]*module.RedisInstance)

func init() {
	prefix := "PD_REDIS_"
	prefixLen := len(prefix)

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		key, value := pair[0], pair[1]

		if strings.HasPrefix(key, prefix) {
			key = key[prefixLen:]
			index := strings.Index(key, "_")
			name := strings.ToLower(key[:index])
			confProperty := key[index+1:]

			if _, ok := redisInstanceMap[name]; !ok {
				redisInstanceMap[name] = &module.RedisInstance{Name: name}
			}
			redisInstance := redisInstanceMap[name]

			switch confProperty {
			case "URL":
				redisInstance.Url = value
			case "CLUSTER_URL":
				redisInstance.ClusterUrl = strings.Split(value, ";")
			case "DEPLOY_TYPE":
				redisInstance.RedisType = value
			case "PASSWORD":
				redisInstance.Password = value
			}
		}
	}
}

type RedisInstanceController struct {
	RedisController
}

func (r *RedisInstanceController) Get() {
	instances := make([]*module.RedisInstance, 0, len(redisInstanceMap))
	for _, instance := range redisInstanceMap {
		instances = append(instances, instance)
	}
	r.Data["json"] = instances
	_ = r.ServeJSON()
}
