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

//// for nodes/overview
//func (kh K8sHandler) GetNodeOverview() ([]cm.Node, error) {
//	result, err := kh.GetNodeList()
//	return result, err
//}
//
//// for node/:name
//func (kh K8sHandler) GetNodeDetail(nodeName string) (cm.Node, error) {
//	var result cm.Node
//
//	node, err := kh.GetNode(nodeName)
//	if err != nil {
//		return cm.Node{}, err
//	}
//
//	cpuUsage, ramUsage, diskAllocated := kh.GetMetricUsage(*node)
//
//	hours24, hours12, hours6 := kh.GetRecordOfNode(nodeName)
//
//	podList, err := kh.GetPodsByNode(nodeName)
//
//	if err != nil {
//		return cm.Node{}, err
//	}
//
//	result = cm.Node{
//		Name:          nodeName,
//		CpuUsage:      cpuUsage,
//		RamUsage:      ramUsage,
//		DiskAllocated: diskAllocated,
//		IP:            node.Status.Addresses[0].Address,
//		Ready:         string(node.Status.Conditions[0].Status),
//		OsImage:       node.Status.NodeInfo.OSImage,
//		Pods:          kh.TransferPod(podList),
//		Record: map[string]cm.RecordOfNode{
//			"24hours": hours24,
//			"12hours": hours12,
//			"6hours":  hours6,
//		},
//	}
//	return result, nil
//}

//// controllers/deployments/overview
//func (kh K8sHandler) GetDeploymentOverview() ([]cm.Deployment, error) {
//
//	var result []cm.Deployment
//
//	deployList, err := kh.GetDeploymentList()
//	if err != nil {
//		return []cm.Deployment{}, err
//	}
//
//	for _, deploy := range deployList.Items {
//		status, image := CheckContainerOfDeploy(deploy)
//		result = append(result, cm.Deployment{
//			Name:      deploy.GetName(),
//			Namespace: deploy.GetNamespace(),
//			Image:     image,
//			Status:    status,
//			Labels:    deploy.Labels,
//			Created:   deploy.GetCreationTimestamp().String(),
//		})
//	}
//	return result, nil
//}

//// GetDeploymentSpecific retrieves information of a deployment. This will also get details as well.
//func (kh K8sHandler) GetDeploymentSpecific(namespace string, name string) (cm.Deployment, error) {
//	ret := cm.Deployment{}
//
//	deployment, err := kh.K8sClient.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
//	if err != nil {
//		return ret, err
//	}
//
//	ret.Name = deployment.GetName()
//	ret.Namespace = deployment.GetNamespace()
//	//ret.Image = deployment.Spec.Template.Spec.Containers[0].Image
//	//ret.Status = string(deployment.Status.Conditions[0].Status)
//	//ret.Labels = deployment.Labels
//	//ret.Created = deployment.GetCreationTimestamp().String()
//	ret.Details = generateDescribeString(name, ret.Namespace, "deployment")
//
//	return ret, nil
//}

//// controlelrs/ingresses/overview
//func (kh K8sHandler) GetIngressOverview() ([]cm.Ingress, error) {
//	var result []cm.Ingress
//
//	ingressList, err := kh.GetIngressList()
//	if err != nil {
//		return []cm.Ingress{}, err
//	}
//	for _, ingress := range ingressList.Items {
//		result = append(result, cm.Ingress{
//			Name:      ingress.GetName(),
//			Namespace: ingress.GetNamespace(),
//			Labels:    ingress.Labels,
//			Host:      ingress.Spec.Rules[0].Host,
//			Class:     ingress.Spec.IngressClassName,
//			Address:   ingress.Status.LoadBalancer.Ingress[0].IP,
//			Created:   ingress.GetCreationTimestamp().String(),
//		})
//	}
//
//	return result, nil
//}

//// controllers/ingress/:namespace/:name
//func (kh K8sHandler) GetIngressSpecific(name string, namespace string) (cm.Ingress, error) {
//	ret := cm.Ingress{}
//
//	ingress, err := kh.K8sClient.NetworkingV1().Ingresses(namespace).Get(context.TODO(), name, metav1.GetOptions{})
//	if err != nil {
//		return ret, err
//	}
//
//	ret.Name = ingress.GetName()
//	ret.Namespace = ingress.GetNamespace()
//
//	ret.Details = generateDescribeString(name, ret.Namespace, "ingress")
//
//	return ret, nil
//}

//// controllers/jobs/overview
//func (kh K8sHandler) GetJobsOverview() ([]cm.Job, error) {
//	var result []cm.Job
//	jobList, err := kh.GetJobsList()
//	if err != nil {
//		return []cm.Job{}, err
//	}
//
//	for _, job := range jobList.Items {
//		result = append(result, cm.Job{
//			Name:      job.Name,
//			Namespace: job.Namespace,
//			Failed:    job.Status.Failed,
//			Succeeded: job.Status.Succeeded,
//			Created:   job.CreationTimestamp.Time.String(),
//		})
//	}
//	return result, nil
//}

//// controllers/job/:namespace/:name
//func (kh K8sHandler) GetJobSpecific(namespace string, name string) (cm.Job, error) {
//	ret := cm.Job{}
//
//	job, err := kh.K8sClient.BatchV1().Jobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
//	if err != nil {
//		return ret, err
//	}
//
//	ret.Name = job.GetName()
//	ret.Namespace = job.GetNamespace()
//
//	ret.Details = generateDescribeString(name, ret.Namespace, "job")
//
//	return ret, nil
//}
//
//// controllers/daemonsets/overview
//func (kh K8sHandler) GetDaemonSetsOverview() ([]cm.DaemonSet, error) {
//	var result []cm.DaemonSet
//	daemonSetList, err := kh.GetDaemonSetList()
//	if err != nil {
//		return []cm.DaemonSet{}, err
//	}
//
//	for _, daemonSet := range daemonSetList.Items {
//		result = append(result, cm.DaemonSet{
//			Name:           daemonSet.GetName(),
//			Namespace:      daemonSet.GetNamespace(),
//			Labels:         daemonSet.Labels,
//			UpdateStrategy: string(daemonSet.Spec.UpdateStrategy.Type),
//			Created:        daemonSet.CreationTimestamp.Time.String(),
//		})
//	}
//	return result, nil
//}
//
//// controllers/daemonset/:namespace/:name
//func (kh K8sHandler) GetDaemonSetSpecific(namespace string, name string) (cm.DaemonSet, error) {
//	ret := cm.DaemonSet{}
//
//	daemonset, err := kh.K8sClient.AppsV1().DaemonSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
//	if err != nil {
//		return ret, err
//	}
//
//	ret.Name = daemonset.GetName()
//	ret.Namespace = daemonset.GetNamespace()
//
//	ret.Details = generateDescribeString(name, ret.Namespace, "daemonset")
//
//	return ret, nil
//}
//
//// resources/namespaces/overview
//func (kh K8sHandler) GetNamespaceOverview() ([]cm.Namespace, error) {
//	var result []cm.Namespace
//
//	namespaceList, err := kh.GetNamespaceList()
//	if err != nil {
//		return []cm.Namespace{}, err
//	}
//
//	for _, namespace := range namespaceList.Items {
//		result = append(result, cm.Namespace{
//			Name:   namespace.GetName(),
//			Labels: namespace.Labels,
//			Status: string(namespace.Status.Phase),
//		})
//	}
//	return result, nil
//}
//
//// resources/namespaces/overview
//func (kh K8sHandler) GetNamespaceDetail(name string) (cm.Namespace, error) {
//	var result cm.Namespace
//
//	namespace, err := kh.GetNamespaceByName(name)
//	if err != nil {
//		return result, err
//	}
//
//	result = cm.Namespace{
//		Name:        namespace.GetName(),
//		Labels:      namespace.Labels,
//		Status:      string(namespace.Status.Phase),
//		Annotations: namespace.Annotations,
//		Finalizers:  namespace.Finalizers,
//		Created:     namespace.CreationTimestamp.Time.String(),
//	}
//	return result, nil
//}

// resources/pods/overview
//func (kh K8sHandler) GetPodOverview() ([]cm.Pod, error) {
//	var result []cm.Pod
//
//	podList, err := kh.GetPodList()
//	if err != nil {
//		return []cm.Pod{}, err
//	}
//
//	for _, pod := range podList.Items {
//		// Find Container's name
//		var containerNames []string
//		containerStats := pod.Status.ContainerStatuses
//		for _, containerStat := range containerStats {
//			containerNames = append(containerNames, containerStat.ContainerID)
//		}
//
//		result = append(result, cm.Pod{
//			Name:             pod.GetName(),
//			Namespace:        pod.GetNamespace(),
//			PodIP:            pod.Status.PodIP,
//			Status:           string(pod.Status.Phase),
//			ServiceConnected: pod.Spec.EnableServiceLinks,
//			Restarts:         GetRestartCount(pod),
//			Image:            CheckContainerOfPod(pod).Image,
//			Age:              pod.CreationTimestamp.String(),
//			ContainerNames:   containerNames,
//			Timestamp:        time.Now(), // not pod's creation time , just for db query
//		})
//	}
//	return result, nil
//}

//// resources/pod/:name
//// for node/:name
//func (kh K8sHandler) GetPodDetail(podName string) (cm.Pod, error) {
//
//	result, err := kh.GetRecordOfPod(podName)
//	if err != nil {
//		return result, err
//	}
//
//	//pod, err := kh.GetPodByName(namespace, podName)
//	//if err != nil {
//	//	return cm.Pod{}, errs
//	//}
//	//result = cm.Pod{
//	//	Name:             pod.GetName(),
//	//	Namespace:        pod.GetNamespace(),
//	//	PodIP:            pod.Status.PodIP,
//	//	Status:           string(pod.Status.Phase),
//	//	ServiceConnected: pod.Spec.EnableServiceLinks,
//	//	Restarts:         GetRestartCount(*pod),
//	//	Image:            pod.Status.ContainerStatuses[0].Image,
//	//	Age:              pod.CreationTimestamp.String(),
//	//}
//	//return result, nil
//
//	return result, nil
//}

//func (kh K8sHandler) GetServiceOverview() ([]cm.Service, error) {
//
//	var result []cm.Service
//	services, err := kh.GetServiceList()
//	if err != nil {
//		return result, err
//	}
//	for _, service := range services.Items {
//
//		result = append(result, cm.Service{
//			Name:       service.GetName(),
//			Namespace:  service.GetNamespace(),
//			Type:       string(service.Spec.Type),
//			ClusterIP:  service.Spec.ClusterIP,
//			ExternalIP: service.Spec.ExternalName,
//			Port:       service.Spec.Ports[0].Port,
//			NodePort:   service.Spec.Ports[0].NodePort,
//		})
//	}
//	return result, err
//}
//
//func (kh K8sHandler) GetServiceDetail(namespace string, name string) (cm.Service, error) {
//	var result cm.Service
//
//	service, err := kh.GetServiceByName(namespace, name)
//	if err != nil {
//		return result, err
//	}
//
//	result = cm.Service{
//		Name:       service.GetName(),
//		Namespace:  service.GetNamespace(),
//		Type:       string(service.Spec.Type),
//		ClusterIP:  service.Spec.ClusterIP,
//		ExternalIP: service.Spec.ExternalName,
//		Port:       service.Spec.Ports[0].Port,
//		NodePort:   service.Spec.Ports[0].NodePort,
//		Selector:   service.Spec.Selector,
//		Conditions: service.Status.Conditions,
//		Labels:     service.Labels,
//		Created:    service.CreationTimestamp.Time.String(),
//	}
//
//	return result, err
//}

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
		result.Created = event.LastTimestamp.Time.String()
		result.Type = event.Type
		result.Name = event.InvolvedObject.Name
		result.Status = event.Reason
		result.Message = event.Message

		kh.StoreEventInDB(result)
	}
}

// overview
func (kh K8sHandler) GetOverviewStatus() (cm.Overview, error) {
	var result cm.Overview

	ready, notReady := kh.nodeStatus()
	running, pending, error := kh.podStatus()

	result = cm.Overview{
		Version: kh.GetKubeVersion(),
		NodeStatus: cm.NodeStatus{
			NotReady: notReady,
			Ready:    ready,
		},
		PodStatus: cm.PodStatus{
			Error:   error,
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
	var cpuUsage, ramUsage float64

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

		// Find pod details
		podName = pod.GetName()
		namespace = pod.GetNamespace()
		cpuUsage, ramUsage, err = kh.calculatePodUsage(podName, namespace)
		if err != nil {
			log.Println(err)
		}

		// Find volume details
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
		})
	}

	return result, nil
}

func (kh K8sHandler) GetController() ([]cm.Controller, error) {

	var result []cm.Controller

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
		volumes := deploy.Spec.Template.Spec.Volumes
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

		result = append(result, cm.Controller{
			Name:      statefulSet.GetName(),
			Type:      "StatefulSet",
			Namespace: statefulSet.GetNamespace(),
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

		result = append(result, cm.Controller{
			Name:      daemonSet.GetName(),
			Type:      "DaemonSet",
			Namespace: daemonSet.GetNamespace(),
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
