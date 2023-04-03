package controllers

type KubernetesInstancesController struct {
	KubernetesController
}

func (k *KubernetesInstancesController) Get() {
	defaultInstance := map[string]string{"name": "default"}
	instanceList := []map[string]string{defaultInstance}
	k.Data["json"] = instanceList
	_ = k.ServeJSON()
}
