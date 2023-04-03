package controllers

import (
	"github.com/beego/beego/v2/server/web"
)

type PulsarController struct {
	web.Controller
}

func RegisterPulsarControllers() {
	web.Router("/api/pulsar/instances", &PulsarInstancesController{}, "get:Get")
	web.Router("/api/pulsar/instances/clear_inactive_topics", &PulsarClearInactiveTopicsController{}, "get:Get")
}
