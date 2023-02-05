package k8s

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"math"
	"time"
)

// TODO Filtering (파라미터에 namespace등을 받는 걸로 고려 중)

// Get Kubernetes Version of Node

func (kh K8sHandler) GetKubeVersion() string {
	nodeList, err := kh.GetNodeList()

	if err != nil {
		log.Printf("failed to get node list %s", err)
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
	nodeList, _ := kh.GetNodeList()

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

func (kh K8sHandler) GetNodeList() (*corev1.NodeList, error) {
	nodes, err := kh.K8sClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get node list %s", err)
	}
	return nodes, nil
}

func (kh K8sHandler) GetMetricUsage(node corev1.Node) (cpuUsage float64, memoryUsage float64, diskUsage float64) {
	metrics, err := kh.MetricK8sClient.MetricsV1beta1().NodeMetricses().Get(context.TODO(), node.GetName(), metav1.GetOptions{})
	if err != nil {
		fmt.Errorf("failed to get node metrics: %s", err)
		return 0, 0, 0
	}

	allocatableCpu := float64(node.Status.Allocatable.Cpu().MilliValue())
	allocatableRam := float64(node.Status.Allocatable.Memory().MilliValue())
	allocatableDisk := float64(node.Status.Capacity.StorageEphemeral().MilliValue())

	usingCpu := float64(metrics.Usage.Cpu().MilliValue())
	usingRam := float64(metrics.Usage.Memory().MilliValue())
	usingDisk := float64(metrics.Usage.StorageEphemeral().MilliValue())

	usageCpu := ToPercentage(usingCpu, allocatableCpu)
	usageMemory := ToPercentage(usingRam, allocatableRam)
	usageDisk := ToPercentage(usingDisk, allocatableDisk)

	return usageCpu, usageMemory, usageDisk
}

// -- Namespace -- //

// TODO
func (kh K8sHandler) GetTopNamespaces() []string {

	return nil
}

func (kh K8sHandler) GetTotalNamespaces() int {
	namespaceList, _ := kh.GetNamespaceList()

	return len(namespaceList.Items)
}

func (kh K8sHandler) GetNamespaceList() (*corev1.NamespaceList, error) {
	fmt.Println("namespace")
	namespaces, err := kh.K8sClient.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to get namespace list")
	}
	return namespaces, nil
}

// -- Deployment -- //

func (kh K8sHandler) GetTotalDeploy() int {
	deployList, _ := kh.GetDeploymentList()

	return len(deployList.Items)
}

func (kh K8sHandler) GetDeploymentList() (*appsv1.DeploymentList, error) {
	deploys, err := kh.K8sClient.AppsV1().Deployments(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment list %s", err)
	}
	return deploys, nil
}

// -- Pod -- //

func (kh K8sHandler) GetTotalPods() int {
	podList, _ := kh.GetPodList()

	return len(podList.Items)
}

func (kh K8sHandler) GetPodList() (*corev1.PodList, error) {
	pods, err := kh.K8sClient.CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod list %s", err)
	}
	return pods, nil
}

// -- Ingress -- //

func (kh K8sHandler) GetTotalIngresses() int {
	ingressList, _ := kh.GetIngressList()

	return len(ingressList.Items)
}

func (kh K8sHandler) GetIngressList() (*networkv1.IngressList, error) {
	ingresses, err := kh.K8sClient.NetworkingV1().Ingresses(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get ingress list %s", err)
	}

	return ingresses, nil
}

// -- Service -- //

func (kh K8sHandler) GetTotalServices() int {
	serviceList, _ := kh.GetServiceList()

	return len(serviceList.Items)
}

func (kh K8sHandler) GetServiceList() (*corev1.ServiceList, error) {
	services, err := kh.K8sClient.CoreV1().Services(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, fmt.Errorf("failed to get service list %s", err)
	}

	return services, nil
}

// -- Persistent Volume -- //
func (kh K8sHandler) GetTotalPersistentVolumes() int {
	pvList, _ := kh.GetPersistentVolumeList()

	return len(pvList.Items)
}

func (kh K8sHandler) GetPersistentVolumeList() (*corev1.PersistentVolumeList, error) {
	pvs, err := kh.K8sClient.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, fmt.Errorf("failed to get persistent volume list %s", err)
	}

	return pvs, nil
}

// -- Jobs -- //

func (kh K8sHandler) GetTotalJobs() int {
	jobList, _ := kh.GetJobsList()

	return len(jobList.Items)
}

func (kh K8sHandler) GetJobsList() (*v1.JobList, error) {
	jobs, err := kh.K8sClient.BatchV1().Jobs(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, fmt.Errorf("failed to get job list %s", err)
	}

	return jobs, nil
}

// -- Daemon sets -- //

func (kh K8sHandler) GetTotalDaemonSets() int {
	dsList, _ := kh.GetDaemonSetList()

	return len(dsList.Items)
}

func (kh K8sHandler) GetDaemonSetList() (*appsv1.DaemonSetList, error) {
	daemonsets, err := kh.K8sClient.AppsV1().DaemonSets(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, fmt.Errorf("failed to get daemonset list %s", err)
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
