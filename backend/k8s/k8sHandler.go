package k8s

import (
	cm "github.com/boanlab/kargos/backend/common"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

// K8s
type K8sHandler struct {
	K8sClient       *kubernetes.Clientset
	MetricK8sClient *versioned.Clientset
	// TODO DB Manager
}

func NewK8sHandler() *K8sHandler {
	kh := &K8sHandler{}

	kh.K8sClient = cm.InitK8sClient()
	kh.MetricK8sClient = cm.InitMetricK8sClient()

	return kh
}

// for Overview/main
func (kh K8sHandler) GetHome() cm.Home {
	var result cm.Home

	// TODO 프론트에서 TotalResources 는 처리 불가능할까?
	namespaces := kh.GetTotalNamespaces()
	deployments := kh.GetTotalDeploy()
	pods := kh.GetTotalPods()
	ingresses := kh.GetTotalIngresses()
	services := kh.GetTotalServices()
	persistentVolumes := kh.GetTotalPersistentVolumes()
	jobs := kh.GetTotalJobs()
	daemonSets := kh.GetTotalDaemonSets()

	result = cm.Home{
		Version:    kh.GetVersion(),
		TotalNodes: kh.GetTotalNodes(),
		Created:    kh.GetCreatedOfCluster(),
		Tabs: map[string]int{
			"TotalResources":   namespaces + deployments + pods + ingresses + services + persistentVolumes + jobs + daemonSets,
			"Namespaces":       kh.GetTotalNamespaces(),
			"Deployments":      kh.GetTotalDeploy(),
			"Pods":             kh.GetTotalPods(),
			"Ingresses":        kh.GetTotalIngresses(),
			"Services":         kh.GetTotalServices(),
			"PersistentVolume": kh.GetTotalPersistentVolumes(),
			"Jobs":             kh.GetTotalJobs(),
			"DaemonSets":       kh.GetTotalDaemonSets(),
		},

		// TODO
		TopNamespaces: kh.GetTopNamespaces(),
		AlertCount:    kh.GetAlertsCount(),
	}

	return result
}
