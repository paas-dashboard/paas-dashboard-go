package controllers

import (
	"github.com/beego/beego/v2/server/web"
)

type KafkaController struct {
	web.Controller
}

func RegisterKafkaControllers() {
	web.Router("/api/kafka/instances", &KafkaInstancesController{}, "get:Get")
}
