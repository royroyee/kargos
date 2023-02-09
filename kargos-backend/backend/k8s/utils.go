package k8s

import (
	"context"
	"fmt"
	cm "github.com/boanlab/kargos/common"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"math"
	"time"
)

// TODO Filtering

// Get Kubernetes Version of Node

func (kh K8sHandler) GetKubeVersion() string {
	nodeList, err := kh.K8sClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return ""
	}

	var result string
	for _, node := range nodeList.Items {
		result = node.Status.NodeInfo.KubeletVersion
	}

	return result
}

// Get Total Resources in Cluster (for overview/main)
func (kh K8sHandler) GetTotalResources() (totalResources int, totalNamespaces int, totalDeployments int, totalPods int, totalIngresses int, totalServices int, totalPersistentVolumes int, totalJobs int, totalDaemonsets int, err error) {
	namespaceList, err := kh.GetNamespaceList()
	if err != nil {
		log.Println(err)
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, err
	}

	totalNamespaces = len(namespaceList.Items)

	deploymentList, err := kh.GetDeploymentList()
	if err != nil {
		log.Println(err)
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, err
	}
	totalDeployments = len(deploymentList.Items)

	podList, err := kh.GetPodList()
	if err != nil {
		log.Println(err)
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, err
	}
	totalPods = len(podList.Items)

	ingressList, err := kh.GetIngressList()
	if err != nil {
		log.Println(err)
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, err
	}
	totalIngresses = len(ingressList.Items)

	serviceList, err := kh.GetServiceList()
	if err != nil {
		log.Println(err)
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, err
	}
	totalServices = len(serviceList.Items)

	persistentVolumeList, err := kh.GetPersistentVolumeList()
	if err != nil {
		log.Println(err)
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, err
	}
	totalPersistentVolumes = len(persistentVolumeList.Items)

	jobList, err := kh.GetJobsList()
	if err != nil {
		log.Println(err)
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, err
	}
	totalJobs = len(jobList.Items)

	daemonsetList, err := kh.GetDaemonSetList()
	if err != nil {
		log.Println(err)
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, err
	}
	totalDaemonsets = len(daemonsetList.Items)

	totalResources = totalNamespaces + totalDeployments + totalPods + totalIngresses + totalServices + totalPersistentVolumes + totalJobs + totalDaemonsets

	return totalResources, totalNamespaces, totalDeployments, totalPods, totalIngresses, totalServices, totalPersistentVolumes, totalJobs, totalDaemonsets, nil
}

// Get Number of Nodes in Cluster
func (kh K8sHandler) GetTotalNodes() int {
	nodeList, err := kh.K8sClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return 0
	}
	return len(nodeList.Items)
}

// Get Date Created of Cluster
func (kh K8sHandler) GetCreatedOfCluster() string {
	// Get the Info of Cluster
	cluster, _ := kh.GetClusterInfo()

	// Get the creation timestamp of the cluster
	creationTimestamp := cluster.CreationTimestamp.Time

	// Format the creation timestamp as a string
	creationTime := creationTimestamp.Format(time.RFC3339)

	return creationTime
}

// TODO
func (kh K8sHandler) GetAlertsCount() int {
	return 0
}

// -- Cluster -- //

func (kh K8sHandler) GetClusterInfo() (*corev1.ComponentStatus, error) {
	cluster, err := kh.K8sClient.CoreV1().ComponentStatuses().Get(context.TODO(), "controller-manager", metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster information %s", err)
	}
	return cluster, nil
}

// --- Node -- //

// Get NodeList
func (kh K8sHandler) GetNodeList() ([]cm.Node, error) {
	var result []cm.Node
	nodeList, err := kh.K8sClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		log.Println(err)
		return []cm.Node{}, err
	}

	for _, node := range nodeList.Items {

		// TODO fix diskUsage (only return zero)
		cpuUsage, ramUsage, diskAllocated := kh.GetMetricUsage(node)

		result = append(result, cm.Node{
			Name:          node.GetName(),
			CpuUsage:      cpuUsage,
			RamUsage:      ramUsage,
			DiskAllocated: diskAllocated,
			IP:            node.Status.Addresses[0].Address,
		})
	}
	return result, nil
}

// To Store Metrics of node in DB
func (kh K8sHandler) GetNodeMetric() ([]cm.RecordOfNode, error) {
	var result []cm.RecordOfNode
	nodeList, err := kh.K8sClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		log.Println(err)
		return result, err
	}

	for _, node := range nodeList.Items {

		// TODO fix diskUsage (only return zero)
		cpuUsage, ramUsage, diskAllocated := kh.GetMetricUsage(node)

		result = append(result, cm.RecordOfNode{
			Name:          node.GetName(),
			CpuUsage:      cpuUsage,
			RamUsage:      ramUsage,
			DiskAllocated: diskAllocated,
			Timestamp:     time.Now(),
		})
	}
	return result, nil
}

// Get Node (name)
func (kh K8sHandler) GetNode(nodeName string) (*corev1.Node, error) {
	node, err := kh.K8sClient.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return node, nil
}

func (kh K8sHandler) GetMetricUsage(node corev1.Node) (cpuUsage float64, memoryUsage float64, diskAllocated float64) {
	metrics, err := kh.MetricK8sClient.MetricsV1beta1().NodeMetricses().Get(context.TODO(), node.GetName(), metav1.GetOptions{})
	if err != nil {
		fmt.Errorf("failed to get node metrics: %s", err)
		return 0, 0, 0
	}

	allocatableCpu := float64(node.Status.Allocatable.Cpu().MilliValue())
	allocatableRam := float64(node.Status.Allocatable.Memory().MilliValue())
	diskAllocated = float64(node.Status.Capacity.StorageEphemeral().MilliValue())

	diskAllocated = math.Round((diskAllocated / (1024 * 1024 * 1024)) / 1000)

	usingCpu := float64(metrics.Usage.Cpu().MilliValue())
	usingRam := float64(metrics.Usage.Memory().MilliValue())
	//	usingDisk := float64(metrics.Usage.Storage().Value())

	usageCpu := ToPercentage(usingCpu, allocatableCpu)
	usageMemory := ToPercentage(usingRam, allocatableRam)
	//	usageDisk := ToPercentage(usingDisk, allocatableDisk)

	return usageCpu, usageMemory, diskAllocated
}

// -- Namespace -- //

// to get 5 namespaces with status "Active"
func (kh K8sHandler) GetTopNamespaces() []string {
	var result []string
	namespaceList, err := kh.GetNamespaceList()
	if err != nil {
		log.Println(err)
		return result
	}

	// Create a map of namespace usage
	usage := make(map[string]bool)
	count := 0
	for _, namespace := range namespaceList.Items {
		usage[namespace.Name] = namespace.Status.Phase == "Active"
		count += 1
		if count >= 5 {
			break
		}
	}

	// Sort the namespaces by usage

	for name := range usage {
		result = append(result, name)
	}

	return result
}

func (kh K8sHandler) GetNamespaceList() (*corev1.NamespaceList, error) {
	namespaces, err := kh.K8sClient.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return namespaces, nil
}

func (kh K8sHandler) GetNamespaceByName(name string) (*corev1.Namespace, error) {
	namespace, err := kh.K8sClient.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return namespace, nil
}

// -- Deployment -- //

func (kh K8sHandler) GetDeploymentList() (*appsv1.DeploymentList, error) {
	deploys, err := kh.K8sClient.AppsV1().Deployments(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return deploys, nil
}

// -- Pod -- //

// Get Pod (name)
func (kh K8sHandler) GetPodByName(namespace string, podName string) (*corev1.Pod, error) {
	pod, err := kh.K8sClient.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return pod, nil
}

func (kh K8sHandler) GetPodList() (*corev1.PodList, error) {
	pods, err := kh.K8sClient.CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to get pod list %s", err)
	}
	return pods, nil
}

func (kh K8sHandler) GetPodsByNode(nodeName string) (*corev1.PodList, error) {
	pods, err := kh.K8sClient.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
	})
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to get deployment list %s", err)
	}

	return pods, nil
}

func (kh K8sHandler) TransferPod(podList *corev1.PodList) []cm.Pod {
	var result []cm.Pod

	for _, pod := range podList.Items {
		result = append(result, cm.Pod{
			Name:   pod.Name,
			Status: string(pod.Status.Phase),
			Image:  pod.Spec.Containers[0].Image,
			Age:    pod.CreationTimestamp.Time.Format(time.RFC3339),
		})
	}

	return result

}

// -- Ingress -- //

func (kh K8sHandler) GetIngressList() (*networkv1.IngressList, error) {
	ingresses, err := kh.K8sClient.NetworkingV1().Ingresses(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to get ingress list %s", err)
	}

	return ingresses, nil
}

// -- Service -- //

func (kh K8sHandler) GetServiceList() (*corev1.ServiceList, error) {
	services, err := kh.K8sClient.CoreV1().Services(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to get service list %s", err)
	}

	return services, nil
}

func (kh K8sHandler) GetServiceByName(namespace string, serviceName string) (*corev1.Service, error) {
	service, err := kh.K8sClient.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return service, nil
}

// -- Persistent Volume -- //

func (kh K8sHandler) GetPersistentVolumeList() (*corev1.PersistentVolumeList, error) {
	pvs, err := kh.K8sClient.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return pvs, nil
}

func (kh K8sHandler) GetPersistentVolumeByName(name string) (*corev1.PersistentVolume, error) {
	pv, err := kh.K8sClient.CoreV1().PersistentVolumes().Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return pv, nil
}

// -- Jobs -- //

func (kh K8sHandler) GetJobsList() (*v1.JobList, error) {
	jobs, err := kh.K8sClient.BatchV1().Jobs(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return jobs, nil
}

// -- Daemon sets -- //

func (kh K8sHandler) GetDaemonSetList() (*appsv1.DaemonSetList, error) {
	daemonsets, err := kh.K8sClient.AppsV1().DaemonSets(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return daemonsets, nil
}
func GetPersistentVolumeClaim(pv *corev1.PersistentVolume) string {
	var claim string

	if pv.Spec.ClaimRef != nil {
		claim = pv.Spec.ClaimRef.Namespace + "/" + pv.Spec.ClaimRef.Name
	}
	return claim
}

// -- util -- //

// Round float to 2 demical places & percentage
func ToPercentage(val1 float64, val2 float64) float64 {
	result := (val1 / val2) * 100
	result = math.Round(result)
	return result
}

// counting restart of pod
func GetRestartCount(pod corev1.Pod) int32 {
	var restartCount int32 = 0
	for _, containerStatus := range pod.Status.ContainerStatuses {
		restartCount += containerStatus.RestartCount
	}
	return restartCount
}

func CheckContainerOfPod(pod corev1.Pod) string {
	if len(pod.Spec.Containers) > 0 {
		return pod.Spec.Containers[0].Image

	} else {
		return "unknown"
	}
}

func CheckContainerOfDeploy(deployment appsv1.Deployment) (status string, image string) {

	status = "unknown"
	image = "unknown"

	if len(deployment.Status.Conditions) > 0 {
		status = string(deployment.Status.Conditions[0].Status)
	}

	if len(deployment.Spec.Template.Spec.Containers) > 0 {
		image = deployment.Spec.Template.Spec.Containers[0].Image
	}

	return status, image
}
