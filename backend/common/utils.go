package common

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

// Init Kubernetes Client
func InitK8sClient() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the Kubernetes Client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return client
}

// Init Kubernetes Metric Client
func InitMetricK8sClient() *versioned.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// create the Kubernetes Metric Client
	client, err := metrics.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return client
}
