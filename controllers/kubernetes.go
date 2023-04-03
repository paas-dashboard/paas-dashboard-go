package controllers

import (
	"github.com/beego/beego/v2/server/web"
)

type KubernetesController struct {
	web.Controller
}

func RegisterKubernetesControllers() {
	web.Router("/api/kubernetes/instances", &KubernetesInstancesController{}, "get:Get")

	kubernetesDeployController := &KubernetesDeployController{}

	web.Router("/api/kubernetes/deployments", kubernetesDeployController, "get:GetDeployments")
	web.Router("/api/kubernetes/namespaces/:namespace/deployments", kubernetesDeployController, "get:GetDeployments")
	web.Router("/api/kubernetes/namespaces/:namespace/deployments/:deployment_name", kubernetesDeployController, "post:PatchDeployment")
	web.Router("/api/kubernetes/namespaces/:namespace/deployments/:deployment_name/ready-check", kubernetesDeployController, "post:DeploymentReadyCheck")

	kubernetesStatefulSetController := &KubernetesStatefulSetController{}

	web.Router("/api/kubernetes/stateful-sets", kubernetesStatefulSetController, "get:GetStatefulSets")
	web.Router("/api/kubernetes/namespaces/:namespace/stateful-sets", kubernetesStatefulSetController, "get:GetStatefulSets")
	web.Router("/api/kubernetes/namespaces/:namespace/stateful-sets/:statefulset_name", kubernetesStatefulSetController, "post:PatchStatefulSet")
	web.Router("/api/kubernetes/namespaces/:namespace/stateful-sets/:statefulset_name/ready-check", kubernetesStatefulSetController, "post:StatefulSetReadyCheck")
}
