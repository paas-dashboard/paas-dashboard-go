package main

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"paas-dashboard-go/controllers"
)

func main() {
	client := controllers.KubeClient
	deploymentsList, err := client.AppsV1().Deployments("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, deployment := range deploymentsList.Items {
		println(deployment.Name)
	}
}
