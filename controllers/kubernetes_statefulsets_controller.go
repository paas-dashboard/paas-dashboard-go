package controllers

import (
	"context"
	"github.com/beego/beego/v2/server/web"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type KubernetesStatefulSetController struct {
	web.Controller
}

func (k *KubernetesStatefulSetController) GetStatefulSets() {
	namespace := k.GetString(":namespace")

	var statefulSetsList *v1.StatefulSetList
	var err error
	if namespace != "" {
		statefulSetsList, err = KubeClient.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	} else {
		statefulSetsList, err = KubeClient.AppsV1().StatefulSets("").List(context.TODO(), metav1.ListOptions{})
	}

	if err != nil {
		k.CustomAbort(500, err.Error())
		return
	}

	k.Data["json"] = statefulSetsList.Items
	_ = k.ServeJSON()
}

func (k *KubernetesStatefulSetController) PatchStatefulSet() {
	namespace := k.GetString(":namespace")
	statefulSetName := k.GetString(":statefulset_name")
	patchData := k.Ctx.Input.RequestBody

	_, err := KubeClient.AppsV1().StatefulSets(namespace).Patch(context.TODO(), statefulSetName, types.MergePatchType, patchData, metav1.PatchOptions{})
	if err != nil {
		k.CustomAbort(500, err.Error())
		return
	}

	k.Data["json"] = map[string]string{"message": "Stateful set patched successfully"}
	_ = k.ServeJSON()
}

func (k *KubernetesStatefulSetController) StatefulSetReadyCheck() {
	namespace := k.GetString(":namespace")
	statefulSetName := k.GetString(":statefulset_name")
	image := k.GetString("image")

	statefulSetList, err := KubeClient.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		k.CustomAbort(500, err.Error())
		return
	}

	var statefulSet *v1.StatefulSet
	for _, s := range statefulSetList.Items {
		if s.Name == statefulSetName {
			statefulSet = &s
			break
		}
	}

	if statefulSet == nil {
		k.CustomAbort(404, "Stateful set not found")
		return
	}

	if image != "" && statefulSet.Spec.Template.Spec.Containers[0].Image != image {
		k.CustomAbort(406, "Image not matching")
		return
	}

	if statefulSet.Status.AvailableReplicas != statefulSet.Status.Replicas {
		k.Data["json"] = map[string]bool{"ready": false}
		_ = k.ServeJSON()
		return
	}

	k.Data["json"] = map[string]bool{"ready": true}
	_ = k.ServeJSON()
}
