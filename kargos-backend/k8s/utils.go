package k8s

import (
	"context"
	"errors"
	"fmt"
	cm "github.com/boanlab/kargos/common"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"log"
	"math"
	"math/rand"
	"sort"
	"time"
)

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

		cpuUsage, ramUsage, diskAllocated := kh.GetMetrics(node)

		result = append(result, cm.Node{
			Name:          node.GetName(),
			CpuUsage:      cpuUsage,
			RamUsage:      ramUsage,
			DiskAllocated: diskAllocated,
			NetworkUsage:  float64(rand.Intn(99) + 1), // TODO
			IP:            node.Status.Addresses[0].Address,
			Status:        NodeStatus(&node),
			Timestamp:     time.Now().Format("2006-01-02 15:04"),
		})
	}
	return result, nil
}

func (kh K8sHandler) GetPodInfoList() ([]cm.PodInfo, error) {
	var result []cm.PodInfo
	var controller string
	podList, err := kh.K8sClient.CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		log.Println(err)
		return result, err
	}

	for _, pod := range podList.Items {

		volumes := []string{}
		for _, volume := range pod.Spec.Volumes {
			volumes = append(volumes, volume.Name)
		}

		// Find controller details
		if pod.ObjectMeta.OwnerReferences != nil && len(pod.ObjectMeta.OwnerReferences) > 0 {
			controller = pod.ObjectMeta.OwnerReferences[0].Name
		}

		result = append(result, cm.PodInfo{
			Name:       pod.GetName(),
			Namespace:  pod.GetNamespace(),
			Restarts:   GetRestartCount(pod),
			PodIP:      pod.Status.PodIP,
			Node:       pod.Spec.NodeName,
			Volumes:    volumes,
			Controller: controller,
			Status:     string(pod.Status.Phase),
			Image:      CheckContainerOfPod(pod).Image,
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

func (kh K8sHandler) GetMetrics(node corev1.Node) (cpuUsage float64, memoryUsage float64, diskAllocated float64) {
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

func (kh K8sHandler) GetTopUsage() (nodeCpu map[string]float64, nodeMemory map[string]float64) {
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

// -- Deployment -- //

func (kh K8sHandler) GetDeploymentList() (*appsv1.DeploymentList, error) {
	deploys, err := kh.K8sClient.AppsV1().Deployments(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return deploys, nil
}

func (kh K8sHandler) GetPodList() (*corev1.PodList, error) {
	pods, err := kh.K8sClient.CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to get pod list %s", err)
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

// -- StatefulSets -- //

func (kh K8sHandler) GetStatefulSetList() (*appsv1.StatefulSetList, error) {
	statefulSetList, err := kh.K8sClient.AppsV1().StatefulSets(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return statefulSetList, nil
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

func (kh K8sHandler) nodeStatus() (ready []string, notReady []string, err error) {

	nodeList, err := kh.K8sClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return ready, notReady, err
	}

	for _, node := range nodeList.Items {
		if isNodeReady(&node) {
			ready = append(ready, node.GetName())
		} else {
			notReady = append(ready, node.GetName())
		}
	}
	return ready, notReady, nil
}

func (kh K8sHandler) podStatus() (running int, pending []string, errorStatus []string, err error) {

	running = 0

	podList, err := kh.K8sClient.CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return running, pending, errorStatus, err
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
			errorStatus = append(errorStatus, pod.Name)
		default:
			errorStatus = append(errorStatus, pod.Name)
		}
	}
	return running, pending, errorStatus, nil
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
func (kh K8sHandler) calculatePodUsage(podName string, namespace string) (int64, int64, error) {
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
		return 0.0, 0.0, errors.New("no containers fount in the pod")
	}
	if len(podMetrics.Containers) < 1 {
		return 0.0, 0.0, errors.New("no containers fount in the pod")

	}

	cpuUsage := podMetrics.Containers[0].Usage.Cpu().MilliValue()
	memUsage := podMetrics.Containers[0].Usage.Memory().Value() / 1048576 // Convert bytes to mebibytes

	return cpuUsage, memUsage, nil
}

func (kh K8sHandler) GetReplicaSetList() (*appsv1.ReplicaSetList, error) {
	replicaSetList, err := kh.K8sClient.AppsV1().ReplicaSets(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return replicaSetList, nil
}

func (kh K8sHandler) GetJobList() (*batchv1.JobList, error) {
	jobList, err := kh.K8sClient.BatchV1().Jobs(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return jobList, nil
}

func (kh K8sHandler) GetCronJobList() (*batchv1.CronJobList, error) {
	cronjobList, err := kh.K8sClient.BatchV1().CronJobs(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return cronjobList, nil
}

func (kh K8sHandler) GetPVCInfo(namespace, pvcName string) (string, string, bool, error) {
	pvc, err := kh.K8sClient.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), pvcName, metav1.GetOptions{})
	if err != nil {
		return "", "", false, err
	}

	claimName := pvc.ObjectMeta.GetName()
	accessMode := pvc.Spec.AccessModes[0]
	readonly := pvc.Spec.VolumeMode != nil && *pvc.Spec.VolumeMode == corev1.PersistentVolumeFilesystem

	return string(accessMode), claimName, readonly, nil
}

//func (kh K8sHandler) GetVolumeInfo(volumeName string, namespace string) (string, string, bool, error) {
//	pv, err := kh.K8sClient.CoreV1().PersistentVolumes().Get(context.TODO(), volumeName, metav1.GetOptions{})
//	if err != nil {
//		return "", "", false, err
//	}
//
//	if pv.Spec.ClaimRef == nil {
//		return "", "", false, fmt.Errorf("No claim found for volume %s", volumeName)
//	}
//
//	pvcName := pv.Spec.ClaimRef.Name
//	pvc, err := kh.K8sClient.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), pvcName, metav1.GetOptions{})
//	if err != nil {
//		return "", "", false, err
//	}
//
//	return pvc.Spec.VolumeMode, pvcName, pvc.Spec.ReadOnly, nil
//}
