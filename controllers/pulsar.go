package controllers

import (
	"github.com/beego/beego/v2/server/web"
)

type PulsarController struct {
	web.Controller
}

func RegisterPulsarControllers() {
	web.Router("/api/pulsar/instances", &PulsarInstancesController{}, "get:Get")
	pulsarClearInactiveTopicsController := &PulsarClearInactiveTopicsController{}
	web.Router("/api/pulsar/instances/:instance/clear_inactive_topics", pulsarClearInactiveTopicsController, "get:ClearInactiveTopics")
	web.Router("/api/pulsar/instances/:instance/clear_inactive_topics/refresh_status/:taskId", pulsarClearInactiveTopicsController, "get:GetStatus")
}
