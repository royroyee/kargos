package common

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

// Init Kubernetes Client (In Cluster)
// TODO AUTOMATIC
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

//
//func ClientSetOutofCluster() *kubernetes.Clientset {
//	// Check environment variables for K8s API
//	k8sAPIAddr := os.Getenv("K8S_API_LISTEN_ADDR")
//	k8sAPIPort := os.Getenv("K8S_API_LISTEN_PORT")
//	k8sConfig := os.Getenv("K8S_API_CONFIG")
//
//	if len(k8sAPIAddr) == 0 || len(k8sAPIPort) == 0 || len(k8sConfig) == 0 {
//		log.Fatalf("k8s API was set invalid: %s:%s for API server and %s for config", k8sAPIAddr, k8sAPIPort, k8sConfig)
//	}
//
//	config, err := clientcmd.BuildConfigFromFlags("https://"+k8sAPIAddr+":"+k8sAPIPort, k8sConfig)
//	if err != nil {
//		panic(err)
//	}
//	cs, err := kubernetes.NewForConfig(config)
//	if err != nil {
//		panic(err)
//	}
//	return cs
//}
//
//func MetricClientSetOutofCluster() *versioned.Clientset {
//	// Check environment variables for K8s API
//	k8sAPIAddr := os.Getenv("K8S_API_LISTEN_ADDR")
//	k8sAPIPort := os.Getenv("K8S_API_LISTEN_PORT")
//	k8sConfig := os.Getenv("K8S_API_CONFIG")
//
//	if len(k8sAPIAddr) == 0 || len(k8sAPIPort) == 0 || len(k8sConfig) == 0 {
//		log.Fatalf("k8s API was set invalid: %s:%s for API server and %s for config", k8sAPIAddr, k8sAPIPort, k8sConfig)
//	}
//
//	config, err := clientcmd.BuildConfigFromFlags("https://"+k8sAPIAddr+":"+k8sAPIPort, k8sConfig)
//	if err != nil {
//		panic(err)
//	}
//
//	metricClientset, err := metrics.NewForConfig(config)
//	if err != nil {
//		panic(err)
//	}
//
//	return metricClientset
//}
