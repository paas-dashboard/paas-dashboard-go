package controllers

import (
	"context"
	"github.com/beego/beego/v2/core/logs"
	"github.com/protocol-laboratory/pulsar-admin-go/padmin"
	"paas-dashboard-go/task"
	"strings"
)

type PulsarClearInactiveTopicsController struct {
	PulsarController
}

func (c *PulsarClearInactiveTopicsController) ClearInactiveTopics() {
	instance := c.GetString(":instance")
	ins, loaded := PulsarInstanceMap[instance]
	if !loaded {
		logs.Error("instance not found")
		return
	}
	admin, err := padmin.NewPulsarAdmin(padmin.Config{
		Host:              ins.Host,
		Port:              ins.WebPort,
		TlsEnable:         false,
		TlsConfig:         nil,
		ConnectionTimeout: 0,
	})
	if err != nil {
		logs.Error("newPulsarAdmin failed : %v", err)
	}
	t, _ := task.New("pulsarClearInactiveTopic", func(ctx context.Context) error {
		return doClearInactiveTopic(admin)
	})

	c.Data["json"] = map[string]string{"taskId": t.Id, "status": t.Status.String()}
	_ = c.ServeJSON()
}

func doClearInactiveTopic(admin *padmin.PulsarAdmin) error {
	tenantList, err := admin.Tenants.List()

	if err != nil {
		return err
	}

	logs.Info("got tenant : %v", tenantList)

	for _, tenantName := range tenantList {
		processNamespace(admin, tenantName)
	}

	return nil
}

func processNamespace(admin *padmin.PulsarAdmin, tenant string) {
	logs.Info("processing tenant : %s", tenant)
	tenantNamespaceList, _ := admin.Namespaces.List(tenant)
	for _, tenantNamespaceName := range tenantNamespaceList {
		list := strings.Split(tenantNamespaceName, "/")
		tenantName, namespaceName := list[0], list[1]
		processTopics(admin, tenantName, namespaceName)
	}
}

func processTopics(admin *padmin.PulsarAdmin, tenant, namespace string) {
	logs.Info("processing namespace %s/%s", tenant, namespace)
	tenantNamespaceTopicList, _ := admin.PersistentTopics.ListPartitioned(tenant, namespace)
	for _, tenantNamespaceTopicName := range tenantNamespaceTopicList {
		list := strings.Split(tenantNamespaceTopicName, "/")
		tenantName, namespaceName, topicName := list[2], list[3], list[4]
		isClearInactiveTopic(admin, tenantName, namespaceName, topicName)
	}
}

func isClearInactiveTopic(admin *padmin.PulsarAdmin, tenant, namespace, topic string) {
	logs.Info("processing topic %s/%s/%s", tenant, namespace, topic)
	err := admin.PersistentTopics.DeletePartitioned(tenant, namespace, topic)
	if err != nil {
		logs.Error("delete partition failed : %v", err)
		return
	}
}

func (c *PulsarClearInactiveTopicsController) GetStatus() {
	taskId := c.GetString(":taskId")

	c.Data["json"] = map[string]string{"taskId": taskId, "status": task.GetTaskStatus(taskId)}
	_ = c.ServeJSON()
}
