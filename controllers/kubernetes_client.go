package controllers

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"paas-dashboard-go/conf"
)

var KubeClient *kubernetes.Clientset

func init() {
	if !conf.KubernetesDisable {
		var config *rest.Config
		var err error

		if os.Getenv("KUBERNETES_DEFAULT_CONFIG_TYPE") == "cluster" {
			config, err = rest.InClusterConfig()
		} else {
			kubeconfig := clientcmd.RecommendedHomeFile
			config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		}

		if err != nil {
			panic(err.Error())
		}

		KubeClient, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
	}
}
