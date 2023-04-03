package controllers

import (
	"k8s.io/client-go/kubernetes"
)

var KubeClient *kubernetes.Clientset

func init() {
	//var config *rest.Config
	//var err error
	//
	//if os.Getenv("KUBERNETES_DEFAULT_CONFIG_TYPE") == "cluster" {
	//	config, err = rest.InClusterConfig()
	//} else {
	//	kubeconfig := clientcmd.RecommendedHomeFile
	//	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	//}
	//
	//if err != nil {
	//	panic(err.Error())
	//}
	//
	//KubeClient, err = kubernetes.NewForConfig(config)
	//if err != nil {
	//	panic(err.Error())
	//}
}
