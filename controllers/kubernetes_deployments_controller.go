package controllers

import (
	"context"
	"github.com/beego/beego/v2/server/web"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type KubernetesDeployController struct {
	web.Controller
}

func (k *KubernetesDeployController) GetDeployments() {
	namespace := k.GetString(":namespace")

	var deploymentsList *v1.DeploymentList
	var err error
	if namespace != "" {
		deploymentsList, err = KubeClient.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	} else {
		deploymentsList, err = KubeClient.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	}

	if err != nil {
		k.CustomAbort(500, err.Error())
		return
	}

	k.Data["json"] = deploymentsList.Items
	_ = k.ServeJSON()
}

func (k *KubernetesDeployController) PatchDeployment() {
	namespace := k.GetString(":namespace")
	deploymentName := k.GetString(":deployment_name")
	patchData := k.Ctx.Input.RequestBody

	_, err := KubeClient.AppsV1().Deployments(namespace).Patch(context.TODO(), deploymentName, types.MergePatchType, patchData, metav1.PatchOptions{})
	if err != nil {
		k.CustomAbort(500, err.Error())
		return
	}

	k.Data["json"] = map[string]string{"message": "Deployment patched successfully"}
	_ = k.ServeJSON()
}

func (k *KubernetesDeployController) DeploymentReadyCheck() {
	namespace := k.GetString(":namespace")
	deploymentName := k.GetString(":deployment_name")
	image := k.GetString("image")

	deploymentList, err := KubeClient.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		k.CustomAbort(500, err.Error())
		return
	}

	var deployment *v1.Deployment
	for _, d := range deploymentList.Items {
		if d.Name == deploymentName {
			deployment = &d
			break
		}
	}

	if deployment == nil {
		k.CustomAbort(404, "Deployment not found")
		return
	}

	if image != "" && deployment.Spec.Template.Spec.Containers[0].Image != image {
		k.CustomAbort(406, "Image not matching")
		return
	}

	if deployment.Status.AvailableReplicas != deployment.Status.Replicas {
		k.Data["json"] = map[string]bool{"ready": false}
		_ = k.ServeJSON()
		return
	}

	k.Data["json"] = map[string]bool{"ready": true}
	_ = k.ServeJSON()
}
