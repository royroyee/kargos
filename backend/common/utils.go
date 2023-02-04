package common

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

// Init Kubernetes Client (In Cluster)
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

// Init Kubernetes Metric Client (In Cluster)
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

func ClientSetOutofCluster() *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("serverhwan.shop:8001", "/Users/kyh-macbook/Kubernetes/kube-config")
	if err != nil {
		panic(err)
	}
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return cs
}

func MetricClientSetOutofCluster() *versioned.Clientset {

	config, err := clientcmd.BuildConfigFromFlags("serverhwan.shop:8001", "/Users/kyh-macbook/Kubernetes/kube-config")
	if err != nil {
		panic(err)
	}

	metricClientset, err := metrics.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return metricClientset
}
