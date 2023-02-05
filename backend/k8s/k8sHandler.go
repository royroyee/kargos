package k8s

import (
	cm "github.com/boanlab/kargos/common"
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

	// In Cluster
	//kh := &K8sHandler{
	//  K8sClient = cm.InitK8sClient()
	//	MetricK8sClient = cm.InitMetricK8sClient()
	//}

	// Out of Cluster
	kh := &K8sHandler{

		K8sClient:       cm.ClientSetOutofCluster(),
		MetricK8sClient: cm.MetricClientSetOutofCluster(),
	}

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
		Version:    kh.GetKubeVersion(),
		TotalNodes: kh.GetTotalNodes(),
		Created:    kh.GetCreatedOfCluster(),
		Tabs: map[string]int{
			"TotalResources":   namespaces + deployments + pods + ingresses + services + persistentVolumes + jobs + daemonSets,
			"Namespaces":       namespaces,
			"Deployments":      deployments,
			"Pods":             pods,
			"Ingresses":        ingresses,
			"Services":         services,
			"PersistentVolume": persistentVolumes,
			"Jobs":             jobs,
			"DaemonSets":       daemonSets,
		},

		// TODO
		TopNamespaces: kh.GetTopNamespaces(),
		AlertCount:    kh.GetAlertsCount(),
	}

	return result
}

// for nodes/overview
func (kh K8sHandler) GetNodeOverview() []cm.Node {
	var result []cm.Node

	var cpuUsage, ramUsage, diskUsage float64
	nodeList, _ := kh.GetNodeList()
	for _, node := range nodeList.Items {

		// TODO fix diskUsage (only return zero)
		cpuUsage, ramUsage, diskUsage = kh.GetMetricUsage(node)

		result = append(result, cm.Node{
			Name:      node.GetName(),
			CpuUsage:  cpuUsage,
			RamUsage:  ramUsage,
			DiskUsage: diskUsage,
			IP:        node.Status.Addresses[0].Address,
		})
	}
	return result
}
