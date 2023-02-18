package k8s

import (
	"context"
	"fmt"
	cm "github.com/boanlab/kargos/common"
	"gopkg.in/mgo.v2"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"log"
	"math"
	"time"
)

// K8s
type K8sHandler struct {
	K8sClient       *kubernetes.Clientset
	MetricK8sClient *versioned.Clientset
	session         *mgo.Session
	// TODO DB Manager
}

func NewK8sHandler() *K8sHandler {

	//In Cluster
	kh := &K8sHandler{
		K8sClient:       cm.InitK8sClient(),
		MetricK8sClient: cm.InitMetricK8sClient(),
		session:         GetDBSession(),
	}

	////// Out of Cluster
	//kh := &K8sHandler{
	//	K8sClient:       cm.ClientSetOutofCluster(),
	//	MetricK8sClient: cm.MetricClientSetOutofCluster(),
	//	session:         GetDBSession(),
	//}

	return kh
}

//// generateDescribeString generates string that represent kubernetes resource like "kubectl describe"
//// The code originated from kubectl source code's kubectl/pkg/cmd/cmd.go
//func generateDescribeString(name string, namespace string, resourceType string) string {
//	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag().WithDiscoveryBurst(300).WithDiscoveryQPS(50.0)
//	cmdutil.NewMatchVersionFlags(kubeConfigFlags)
//	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)
//	f := cmdutil.NewFactory(matchVersionKubeConfigFlags)
//	flags := kubectl.NewDescribeFlags(f, genericclioptions.IOStreams{})
//	o, _ := flags.ToOptions("kubectl", []string{resourceType, name, "namespace", namespace})
//	ret := o.Run()
//	return ret
//}

//

func (kh K8sHandler) WatchEvents() {

	var result cm.Event

	watcher, err := kh.K8sClient.CoreV1().Events(metav1.NamespaceAll).Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println("Failed to create event watcher: %v", err)
	}
	defer watcher.Stop()

	for watch := range watcher.ResultChan() {
		event, ok := watch.Object.(*v1.Event)
		if !ok {
			log.Println("Received non-Event object")
			continue
		}

		result.Created = event.LastTimestamp.Time.Format("2006-01-02 15:04")
		result.Name = event.InvolvedObject.Name
		result.Type = event.InvolvedObject.Kind
		result.Status = event.Reason
		result.Message = event.Message
		result.EventLevel = event.Type

		kh.StoreEventInDB(result)
	}
}

// overview
func (kh K8sHandler) GetOverviewStatus() (cm.Overview, error) {
	var result cm.Overview

	ready, notReady, err := kh.nodeStatus()
	if err != nil {
		return result, err
	}
	running, pending, errorStatus, err := kh.podStatus()
	if err != nil {
		return result, err
	}

	result = cm.Overview{
		NodeStatus: cm.NodeStatus{
			NotReady: notReady,
			Ready:    ready,
		},
		PodStatus: cm.PodStatus{
			Error:   errorStatus,
			Pending: pending,
			Running: running,
		},
	}

	return result, nil
}
func (kh K8sHandler) PodOverview() ([]cm.Pod, error) {

	var result []cm.Pod
	var containerStats []v1.ContainerStatus
	var podName, namespace, controller string
	var cpuUsage, ramUsage int64

	podList, err := kh.GetPodList()
	if err != nil {
		log.Println(err)
		return result, err
	}

	for _, pod := range podList.Items {

		// Find Container's name
		containerNames := []string{}
		containerStats = pod.Status.ContainerStatuses
		for _, containerStat := range containerStats {
			containerNames = append(containerNames, containerStat.Name)
		}

		podName = pod.GetName()
		namespace = pod.GetNamespace()
		cpuUsage, ramUsage, err = kh.calculatePodUsage(podName, namespace)
		if err != nil {
			log.Println(err)
		}

		volumes := []string{}
		for _, volume := range pod.Spec.Volumes {
			volumes = append(volumes, volume.Name)
		}

		// Find controller details
		if pod.ObjectMeta.OwnerReferences != nil && len(pod.ObjectMeta.OwnerReferences) > 0 {
			controller = pod.ObjectMeta.OwnerReferences[0].Name
		}

		result = append(result, cm.Pod{
			Name:           podName,
			Namespace:      namespace,
			CpuUsage:       cpuUsage,
			RamUsage:       ramUsage,
			Restarts:       GetRestartCount(pod),
			PodIP:          pod.Status.PodIP,
			Node:           pod.Spec.NodeName,
			Volumes:        volumes,
			Controller:     controller,
			Status:         string(pod.Status.Phase),
			Image:          CheckContainerOfPod(pod).Image,
			ContainerNames: containerNames,
			Timestamp:      time.Now().Format("2006-01-02 15:04"),
		})
	}

	return result, nil
}

func (kh K8sHandler) GetPodUsage() ([]cm.PodUsage, error) {
	var result []cm.PodUsage
	var containerStats []v1.ContainerStatus

	podList, err := kh.GetPodList()
	if err != nil {
		log.Println(err)
		return result, err
	}

	var podName, namespace string
	var cpuUsage, ramUsage int64

	for _, pod := range podList.Items {

		// Find Container's name
		containerNames := []string{}
		containerStats = pod.Status.ContainerStatuses
		for _, containerStat := range containerStats {
			containerNames = append(containerNames, containerStat.Name)
		}

		podName = pod.GetName()
		namespace = pod.GetNamespace()

		cpuUsage, ramUsage, err = kh.calculatePodUsage(podName, namespace)
		result = append(result, cm.PodUsage{
			Name:     podName,
			CpuUsage: cpuUsage,
			RamUsage: ramUsage,
			// TODO Network , Disk Usage
			Timestamp: time.Now().Format("2006-01-02 15:04"),
		})
	}
	return result, nil
}

func (kh K8sHandler) GetController() ([]cm.Controller, error) {

	var result []cm.Controller
	var volumes []v1.Volume

	deployList, err := kh.GetDeploymentList()
	if err != nil {
		return []cm.Controller{}, err
	}

	for _, deploy := range deployList.Items {

		podList, err := kh.K8sClient.CoreV1().Pods(deploy.Namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: metav1.FormatLabelSelector(deploy.Spec.Selector),
		})
		if err != nil {
			log.Println(err)
		}
		var pods []string
		for _, pod := range podList.Items {
			pods = append(pods, pod.GetName())
		}
		var volumeList []string
		volumes = deploy.Spec.Template.Spec.Volumes
		for _, volume := range volumes {
			volumeList = append(volumeList, volume.Name)
		}

		result = append(result, cm.Controller{
			Name:      deploy.GetName(),
			Type:      "Deployment",
			Namespace: deploy.GetNamespace(),
			Pods:      pods,
			Volumes:   volumeList,
		})
	}

	statefulSetList, err := kh.GetStatefulSetList()
	if err != nil {
		return []cm.Controller{}, err
	}

	for _, statefulSet := range statefulSetList.Items {

		podList, err := kh.K8sClient.CoreV1().Pods(statefulSet.Namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: metav1.FormatLabelSelector(statefulSet.Spec.Selector),
		})
		if err != nil {
			log.Println(err)
		}
		var pods []string
		for _, pod := range podList.Items {
			pods = append(pods, pod.Name)
		}

		var volumeOfState []string
		volumes = statefulSet.Spec.Template.Spec.Volumes
		for _, volume := range volumes {
			volumeOfState = append(volumeOfState, volume.Name)
		}

		result = append(result, cm.Controller{
			Name:      statefulSet.GetName(),
			Type:      "StatefulSet",
			Namespace: statefulSet.GetNamespace(),
			Volumes:   volumeOfState,
		})
	}

	daemonSetList, err := kh.GetDaemonSetList()
	if err != nil {
		return []cm.Controller{}, err
	}

	for _, daemonSet := range daemonSetList.Items {

		podList, err := kh.K8sClient.CoreV1().Pods(daemonSet.Namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: metav1.FormatLabelSelector(daemonSet.Spec.Selector),
		})
		if err != nil {
			log.Println(err)
		}
		var pods []string
		for _, pod := range podList.Items {
			pods = append(pods, pod.Name)
		}
		var volumeOfDaemon []string
		volumes = daemonSet.Spec.Template.Spec.Volumes
		for _, volume := range volumes {
			volumeOfDaemon = append(volumeOfDaemon, volume.Name)
		}

		result = append(result, cm.Controller{
			Name:      daemonSet.GetName(),
			Type:      "DaemonSet",
			Namespace: daemonSet.GetNamespace(),
			Volumes:   volumeOfDaemon,
		})
	}

	return result, nil
}

func (kh K8sHandler) GetLogsOfPod(namespace string, podName string) ([]string, error) {
	var result []string

	// create a time range for the logs
	now := time.Now()
	before := now.Add(-24 * time.Hour) // get logs from the last 24 hours

	// create options for retrieving the logs
	options := &v1.PodLogOptions{
		Timestamps: true,
		SinceTime:  &metav1.Time{Time: before},
	}

	// get the logs for the specified pod
	req := kh.K8sClient.CoreV1().Pods(namespace).GetLogs(podName, options)
	logs, err := req.Stream(context.Background())
	if err != nil {
		return result, err
	}
	defer logs.Close()

	// read the logs and append them to the result slice
	buf := make([]byte, 1024)
	for {
		numBytes, err := logs.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return result, err
		}
		line := string(buf[0:numBytes])
		result = append(result, line)
	}
	// return the result slice
	return result, nil
}

//func (kh K8sHandler) GetLogsOfNode(nodeName string) ([]string, error) {
//	var result []string
//
//	// create a time range for the logs
//	now := time.Now()
//	before := now.Add(-24 * time.Hour) // get logs from the last 24 hours
//
//	// create the REST client for the nodes API
//	restConfig, err := kh.GetRestConfig()
//	if err != nil {
//		return result, err
//	}
//	restClient, err := rest.RESTClientFor(restConfig)
//	if err != nil {
//		return result, err
//	}
//
//	// create the URL for the node logs
//	nodeLogURL := restClient.Post().
//		Resource("nodes").
//		Name(nodeName).
//		SubResource("log").
//		VersionedParams(&v1.PodLogOptions{
//			Timestamps: true,
//			SinceTime:  &metav1.Time{Time: before},
//		}, scheme.ParameterCodec).URL()
//
//	// create the request for the node logs
//	req := restClient.Get().
//		AbsPath(nodeLogURL.String())
//
//	// start the request
//	readCloser, err := req.Stream(context.Background())
//	if err != nil {
//		return result, err
//	}
//	defer readCloser.Close()
//
//	// read the logs and append them to the result slice
//	buf := make([]byte, 1024)
//	for {
//		numBytes, err := readCloser.Read(buf)
//		if err != nil {
//			if err == io.EOF {
//				break
//			}
//			return result, err
//		}
//		line := string(buf[0:numBytes])
//		result = append(result, line)
//	}
//
//	// return the result slice
//	return result, nil
//}

func (kh K8sHandler) GetNodeInfo(nodeName string) (cm.NodeInfo, error) {

	var result cm.NodeInfo

	// get the node with the specified name
	node, err := kh.K8sClient.CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return result, err
	}

	// get the hostname and IP address of the node
	result.HostName = node.ObjectMeta.Name
	result.IP = node.Status.Addresses[0].Address

	// get the Kubernetes version and containerd version of the node
	result.KubeletVersion = node.Status.NodeInfo.KubeletVersion
	result.ContainerRuntimeVersion = node.Status.NodeInfo.ContainerRuntimeVersion

	// get the number of running containers on the node
	pods, err := kh.K8sClient.CoreV1().Pods(metav1.NamespaceAll).List(context.Background(), metav1.ListOptions{FieldSelector: "spec.nodeName=" + result.HostName})
	if err != nil {
		log.Println(err)
		return result, err
	}
	numContainers := 0
	for _, pod := range pods.Items {
		numContainers += len(pod.Spec.Containers)
	}
	result.NumContainers = numContainers
	// get the CPU and RAM capacity of the node
	capacity := node.Status.Capacity
	result.CpuCores = capacity.Cpu().Value()
	ramBytes := capacity.Memory().Value()
	result.Ram = math.Round(float64(ramBytes) / float64(1024*1024*1024))

	return result, nil

}

func (kh K8sHandler) GetEventsByController(namespace string, controllerName string) ([]string, error) {

	var result []string

	// Set up a field selector to only retrieve events for the deployment
	fieldSelector := fmt.Sprintf("involvedObject.name=%s,involvedObject.namespace=%s", controllerName, namespace)

	// Retrieve the events for the deployment
	eventList, err := kh.K8sClient.CoreV1().Events(namespace).List(context.Background(), metav1.ListOptions{
		FieldSelector: fieldSelector,
		Limit:         10,
	})
	if err != nil {
		log.Println(err)
		return result, err
	}

	// Create a string array to hold the event messages
	eventMessages := make([]string, len(eventList.Items))

	// Add each event message to the string array
	for i, event := range eventList.Items {
		eventMessages[i] = fmt.Sprintf("%s: %s", event.LastTimestamp.Format("2023-01-02 15:04:05"), event.Message)
	}

	// Print out the string array of event messagesç
	result = eventMessages
	return result, nil

}

func (kh K8sHandler) NumberOfNodes() (cm.Count, error) {
	var result cm.Count

	nodes, err := kh.K8sClient.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return result, err
	}

	result.Count = len(nodes.Items)
	return result, err
}
