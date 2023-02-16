package k8s

import (
	"context"
	"fmt"
	cm "github.com/boanlab/kargos/common"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"log"
	"math"
	"math/rand"
	"sort"
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

		cpuUsage, ramUsage, diskAllocated := kh.GetMetricUsage(node)

		result = append(result, cm.Node{
			Name:          node.GetName(),
			CpuUsage:      cpuUsage,
			RamUsage:      ramUsage,
			DiskAllocated: diskAllocated,
			NetworkUsage:  rand.Intn(99) + 1, // TODO
			IP:            node.Status.Addresses[0].Address,
			Status:        NodeStatus(&node),
			Timestamp:     time.Now().String(),
		})
	}
	return result, nil
}

//// To Store Metrics of node in DB
//func (kh K8sHandler) GetNodeMetric() ([]cm.RecordOfNode, error) {
//	var result []cm.RecordOfNode
//	nodeList, err := kh.K8sClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
//
//	if err != nil {
//		log.Println(err)
//		return result, err
//	}
//
//	for _, node := range nodeList.Items {
//
//		// TODO fix diskUsage (only return zero)
//		cpuUsage, ramUsage, diskAllocated := kh.GetMetricUsage(node)
//
//		result = append(result, cm.RecordOfNode{
//			Name:          node.GetName(),
//			CpuUsage:      cpuUsage,
//			RamUsage:      ramUsage,
//			DiskAllocated: diskAllocated,
//			Timestamp:     time.Now(),
//		})
//	}
//	return result, nil
//}

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

func (kh K8sHandler) GetTopMetric() (nodeCpu map[string]float64, nodeMemory map[string]float64) {
	nodeList, err := kh.K8sClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return nodeCpu, nodeMemory
	}

	for _, node := range nodeList.Items {
		metrics, err := kh.MetricK8sClient.MetricsV1beta1().NodeMetricses().Get(context.TODO(), node.GetName(), metav1.GetOptions{})
		if err != nil {
			log.Println(err)
			return nodeCpu, nodeMemory
		}

		usageCpu := ToPercentage(float64(metrics.Usage.Cpu().MilliValue()), float64(node.Status.Allocatable.Cpu().MilliValue()))
		usageMemory := ToPercentage(float64(metrics.Usage.Memory().MilliValue()), float64(node.Status.Allocatable.Memory().MilliValue()))

		nodeName := node.GetName()

		nodeCpu[nodeName] = usageCpu
		nodeMemory[nodeName] = usageMemory
	}

	// sort in descending order
	// go의 map은 그 자체로 정렬할 수 없고, slice를 사용해야 하기 때문에 조금 비효율적임
	// 이렇게 정렬할 바에.. db에 저장해서 꺼내는 게 나을듯?..
	//
	var keys []string
	for k := range nodeCpu {
		keys = append(keys, k)
	}

	// Sort the slice of key-value pairs by the values in the map in descending order
	sort.Slice(keys, func(i, j int) bool {
		return nodeCpu[keys[i]] > nodeCpu[keys[j]]
	})

	// Print the sorted map
	for _, k := range keys {
		fmt.Printf("%s: %d\n", k, nodeCpu[k])
	}
	return nodeCpu, nodeMemory
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

func (kh K8sHandler) GetNamespaceName() ([]string, error) {
	var result []string

	namespaces, err := kh.K8sClient.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, namespace := range namespaces.Items {
		result = append(result, namespace.GetName())
	}

	return result, nil
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
			Image:  CheckContainerOfPod(pod).Image,
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
//func (kh K8sHandler) GetJobList() (*v1.JobList, error) {
//	jobs, err := kh.K8sClient.BatchV1().Jobs(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
//	if err != nil {
//		log.Println(err)
//		return nil, err
//	}
//
//	return jobs, nil
//}

// -- StatefulSets -- //

func (kh K8sHandler) GetStatefulSetList() (*appsv1.StatefulSetList, error) {
	staefulSetList, err := kh.K8sClient.AppsV1().StatefulSets(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return staefulSetList, nil
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

// Check if the container has been created //

func CheckContainerOfPod(pod corev1.Pod) corev1.Container {
	if len(pod.Spec.Containers) > 0 {
		return pod.Spec.Containers[0]

	} else {
		return corev1.Container{}
	}
}

func CheckContainerOfPodMetrics(metrics *v1beta1.PodMetrics) *v1beta1.ContainerMetrics {
	if len(metrics.Containers) > 0 {
		return &metrics.Containers[0]

	} else {
		return &v1beta1.ContainerMetrics{}
	}
}

func CheckOwnerOfPod(pod corev1.Pod) metav1.OwnerReference {
	if len(pod.ObjectMeta.OwnerReferences) > 0 {
		return pod.ObjectMeta.OwnerReferences[0]
	} else {
		return metav1.OwnerReference{}
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

func (kh K8sHandler) nodeStatus() (ready []string, notReady []string) {

	nodeList, err := kh.K8sClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return ready, notReady
	}

	for _, node := range nodeList.Items {
		if isNodeReady(&node) {
			ready = append(ready, node.GetName())
		} else {
			notReady = append(ready, node.GetName())
		}
	}
	return ready, notReady
}

func (kh K8sHandler) podStatus() (running int, pending []string, error []string) {

	running = 0

	podList, err := kh.K8sClient.CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return running, pending, error
	}

	for _, pod := range podList.Items {
		switch pod.Status.Phase {
		case corev1.PodPending:
			pending = append(pending, pod.Name)
		case corev1.PodRunning:
			running++
		case corev1.PodSucceeded:
			running++
		case corev1.PodFailed:
			error = append(error, pod.Name)
		default:
			error = append(error, pod.Name)
		}
	}
	return running, pending, error
}

func isNodeReady(node *corev1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			return condition.Status == corev1.ConditionTrue
		}
	}
	return false
}

func NodeStatus(node *corev1.Node) string {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			return "Ready"
		}
	}
	return "Not Ready"
}
func (kh K8sHandler) calculatePodUsage(podName string, namespace string) (cpuPercent float64, memPercent float64, err error) {
	// Get the current CPU and memory usage of the pod
	podMetrics, err := kh.MetricK8sClient.MetricsV1beta1().PodMetricses(namespace).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		return 0.0, 0.0, err
	}

	// Get the CPU and memory limits for the pod
	pod, err := kh.K8sClient.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		return 0.0, 0.0, err
	}
	if len(pod.Spec.Containers) < 1 {
		return 0.0, 0.0, nil
	}
	cpuLimit := pod.Spec.Containers[0].Resources.Limits.Cpu().MilliValue()
	memLimit := pod.Spec.Containers[0].Resources.Limits.Memory().Value()

	// Convert memory usage to bytes
	memUsage := podMetrics.Containers[0].Usage.Memory().Value()

	// Calculate the percentage CPU and memory usage
	cpuPercent = float64(podMetrics.Containers[0].Usage.Cpu().MilliValue()) / float64(cpuLimit) * 100.0
	memPercent = float64(memUsage) / float64(memLimit) * 100.0

	return cpuPercent, memPercent, nil
}
