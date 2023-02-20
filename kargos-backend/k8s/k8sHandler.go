package k8s

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	cm "github.com/boanlab/kargos/common"
	"gopkg.in/mgo.v2"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubectl/pkg/scheme"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"log"
	"strings"
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
		//K8sClient:       cm.ClientSetOutofCluster(),
		//MetricK8sClient: cm.MetricClientSetOutofCluster(),
		session: GetDBSession(),
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
			Name:           podName,
			CpuUsage:       cpuUsage,
			RamUsage:       ramUsage,
			ContainerNames: containerNames,
			// TODO Network , Disk Usage
			Timestamp: time.Now().Format("2006-01-02 15:04"),
		})
	}
	return result, nil
}

func (kh K8sHandler) GetController() []cm.Controller {

	var result []cm.Controller
	var volumes []v1.Volume

	deployList, _ := kh.GetDeploymentList()

	if deployList != nil {

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
				Type:      "deployment",
				Namespace: deploy.GetNamespace(),
				Pods:      pods,
				Volumes:   volumeList,
			})
		}
	}

	statefulSetList, _ := kh.GetStatefulSetList()

	if statefulSetList != nil {

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
				Type:      "statefulset",
				Namespace: statefulSet.GetNamespace(),
				Pods:      pods,
				Volumes:   volumeOfState,
			})
		}
	}

	daemonSetList, _ := kh.GetDaemonSetList()

	if daemonSetList != nil {

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
				Type:      "daemonset",
				Namespace: daemonSet.GetNamespace(),
				Pods:      pods,
				Volumes:   volumeOfDaemon,
			})
		}
	}

	JobList, _ := kh.GetJobList()

	if JobList != nil {

		for _, job := range JobList.Items {

			podList, err := kh.K8sClient.CoreV1().Pods(job.Namespace).List(context.TODO(), metav1.ListOptions{
				LabelSelector: metav1.FormatLabelSelector(job.Spec.Selector),
			})
			if err != nil {
				log.Println(err)
			}
			var pods []string
			for _, pod := range podList.Items {
				pods = append(pods, pod.Name)
			}
			var volumeOfJob []string
			volumes = job.Spec.Template.Spec.Volumes
			for _, volume := range volumes {
				volumeOfJob = append(volumeOfJob, volume.Name)
			}

			result = append(result, cm.Controller{
				Name:      job.GetName(),
				Type:      "job",
				Namespace: job.GetNamespace(),
				Pods:      pods,
				Volumes:   volumeOfJob,
			})
		}
	}

	CronJobList, _ := kh.GetCronJobList()

	if CronJobList != nil {

		for _, cronjob := range CronJobList.Items {

			podList, err := kh.K8sClient.CoreV1().Pods(cronjob.Namespace).List(context.TODO(), metav1.ListOptions{
				LabelSelector: metav1.FormatLabelSelector(cronjob.Spec.JobTemplate.Spec.Selector),
			})
			if err != nil {
				log.Println(err)
			}
			var pods []string
			for _, pod := range podList.Items {
				pods = append(pods, pod.Name)
			}
			var volumeOfCronJob []string
			volumes = cronjob.Spec.JobTemplate.Spec.Template.Spec.Volumes
			for _, volume := range volumes {
				volumeOfCronJob = append(volumeOfCronJob, volume.Name)
			}

			result = append(result, cm.Controller{
				Name:      cronjob.GetName(),
				Type:      "cronjob",
				Namespace: cronjob.GetNamespace(),
				Pods:      pods,
				Volumes:   volumeOfCronJob,
			})
		}
	}

	ReplicaSetList, _ := kh.GetReplicaSetList()

	if ReplicaSetList != nil {
		for _, replicaSet := range ReplicaSetList.Items {

			podList, err := kh.K8sClient.CoreV1().Pods(replicaSet.Namespace).List(context.TODO(), metav1.ListOptions{
				LabelSelector: metav1.FormatLabelSelector(replicaSet.Spec.Selector),
			})
			if err != nil {
				log.Println(err)
			}
			var pods []string
			for _, pod := range podList.Items {
				pods = append(pods, pod.Name)
			}
			var volumeOfReplicaSet []string
			volumes = replicaSet.Spec.Template.Spec.Volumes
			for _, volume := range volumes {
				volumeOfReplicaSet = append(volumeOfReplicaSet, volume.Name)
			}

			result = append(result, cm.Controller{
				Name:      replicaSet.GetName(),
				Type:      "replicaset",
				Namespace: replicaSet.GetNamespace(),
				Pods:      pods,
				Volumes:   volumeOfReplicaSet,
			})
		}
	}

	return result
}

//
//func (kh K8sHandler) GetController() ([]cm.Controller, error) {
//	var result []cm.Controller
//	var volumes []v1.Volume
//
//	// define a list of controller types
//	controllerTypes := []string{"Deployment", "StatefulSet", "DaemonSet"}
//
//	// create a map of controller type to a function that retrieves the list of controllers for that type
//	controllerFuncs := map[string]func() ([]metav1.Object, error){
//		"Deployment":  kh.GetDeploymentList(),
//		"StatefulSet": kh.GetStatefulSetList(),
//		"DaemonSet":   kh.GetDaemonSetList(),
//	}
//
//	// iterate over all controller types
//	for _, controllerType := range controllerTypes {
//		// retrieve the list of controllers for the current type
//		controllers, err := controllerFuncs[controllerType]()
//		if err != nil {
//			return []cm.Controller{}, err
//		}
//
//		// iterate over all controllers for the current type
//		for _, controller := range controllers {
//			// retrieve the pods for the current controller
//			podList, err := kh.K8sClient.CoreV1().Pods(controller.GetNamespace()).List(context.TODO(), metav1.ListOptions{
//				LabelSelector: metav1.FormatLabelSelector(controller.(metav1.Object).GetLabels()),
//			})
//			if err != nil {
//				log.Println(err)
//			}
//
//			// create a list of pod names
//			var pods []string
//			for _, pod := range podList.Items {
//				pods = append(pods, pod.GetName())
//			}
//
//			// create a list of volume names
//			var volumeList []string
//			volumes = controller.(v1beta1.ControllerRevisionInterface).GetTemplate().Spec.Volumes
//			for _, volume := range volumes {
//				volumeList = append(volumeList, volume.Name)
//			}
//
//			// create a controller object and append it to the result list
//			controllerObj := cm.Controller{
//				Name:      controller.GetName(),
//				Type:      controllerType,
//				Namespace: controller.GetNamespace(),
//				Pods:      pods,
//				Volumes:   volumeList,
//			}
//			result = append(result, controllerObj)
//		}
//	}
//
//	return result, nil
//}

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

	// read the logs and format them with timestamps and pod name
	scanner := bufio.NewScanner(logs)
	for scanner.Scan() {
		line := scanner.Text()

		// extract the "$date" field from the JSON object in the log line
		dateStart := strings.Index(line, "{\"$date\":\"")
		dateEnd := strings.Index(line[dateStart+1:], "\"}")
		dateStr := line[dateStart+9 : dateStart+1+dateEnd]

		// parse the timestamp in the log line
		ts, err := time.Parse(time.Now().Format("2006-01-02 15:04"), dateStr)
		if err != nil {
			return result, err
		}

		// format the log line with the timestamp and pod name
		formatted := fmt.Sprintf("%s [%s] %s", ts.Format("2006-01-02 15:04"), podName, line[dateStart+1+dateEnd+3:])
		result = append(result, formatted)
	}

	if err := scanner.Err(); err != nil {
		return result, err
	}

	// return the result slice
	return result, nil
}

func (kh K8sHandler) GetLogsOfNode(nodeName string) ([]string, error) {
	// create a time range for the logs
	now := time.Now()
	before := now.Add(-24 * time.Hour) // get logs from the last 24 hours

	// create the REST client for the nodes API
	restConfig := kh.K8sClient.RESTClient()

	// create the URL for the node logs
	nodeLogURL := restConfig.Post().
		Resource("nodes").
		Name(nodeName).
		SubResource("log").
		VersionedParams(&v1.PodLogOptions{
			Timestamps: true,
			SinceTime:  &metav1.Time{Time: before},
		}, scheme.ParameterCodec).URL()

	// create the request for the node logs
	req := restConfig.Get().
		AbsPath(nodeLogURL.String())

	// start the request
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	readCloser, err := req.Stream(ctx)
	if err != nil {
		return nil, err
	}
	defer readCloser.Close()

	// read the logs and store them in a buffer
	var buf bytes.Buffer
	_, err = io.Copy(&buf, readCloser)
	if err != nil {
		return nil, err
	}

	// split the logs into lines
	scanner := bufio.NewScanner(&buf)
	var result []string
	for scanner.Scan() {
		line := scanner.Text()

		// parse the timestamp in the log line
		tsEnd := strings.Index(line, " ")
		tsStr := line[:tsEnd]

		// format the log line with the timestamp and node name
		ts, err := time.Parse("2006-01-02 15:04", tsStr)
		if err != nil {
			return result, err
		}
		formatted := fmt.Sprintf("%s [%s] %s", ts.Format(time.RFC3339), nodeName, line[tsEnd+1:])
		result = append(result, formatted)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// return the result slice
	return result, nil
}

func (kh K8sHandler) GetNodeInfo(nodeName string) (cm.NodeInfo, error) {

	var result cm.NodeInfo

	// get the node with the specified name
	node, err := kh.K8sClient.CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return result, err
	}

	result.OS = node.Status.NodeInfo.OSImage
	// get the hostname and IP address of the node
	result.HostName = node.ObjectMeta.Name
	result.IP = node.Status.Addresses[0].Address
	result.Status = isNodeReady(node)

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

	result.RamCapacity = node.Status.Capacity.Memory().Value() / 1024 / 1024 / 1024

	return result, nil

}

//func (kh K8sHandler) GetEventsByController(namespace string, controllerName string) ([]string, error) {
//
//	var result []string
//
//	// Set up a field selector to only retrieve events for the deployment
//	fieldSelector := fmt.Sprintf("involvedObject.name=%s,involvedObject.namespace=%s", controllerName, namespace)
//
//	// Retrieve the events for the deployment
//	eventList, err := kh.K8sClient.CoreV1().Events(namespace).List(context.Background(), metav1.ListOptions{
//		FieldSelector: fieldSelector,
//		Limit:         10,
//	})
//	if err != nil {
//		log.Println(err)
//		return result, err
//	}
//
//	// Create a string array to hold the event messages
//	eventMessages := make([]string, len(eventList.Items))
//
//	// Add each event message to the string array
//	for i, event := range eventList.Items {
//		eventMessages[i] = fmt.Sprintf("%s: %s", event.LastTimestamp.Format("2023-01-02 15:04:05"), event.Message)
//	}
//
//	// Print out the string array of event messagesç
//	result = eventMessages
//	return result, nil
//
//}

func (kh K8sHandler) NumberOfNodes() (cm.Count, error) {
	var result cm.Count

	nodes, err := kh.K8sClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return result, err
	}

	result.Count = len(nodes.Items)
	return result, err
}

func (kh K8sHandler) GetControllerInfo(controllerType string, namespace string, controllerName string) (cm.ControllerInfo, error) {
	var result cm.ControllerInfo
	var limits, volumes, mounts, envs, labels []string
	var controlleredByName string

	if controllerType == "deployment" {
		controller, err := kh.K8sClient.AppsV1().Deployments(namespace).Get(context.TODO(), controllerName, metav1.GetOptions{})
		if err != nil {
			return result, err
		}

		container := controller.Spec.Template.Spec.Containers

		if len(container) > 0 {
			for _, container := range controller.Spec.Template.Spec.Containers {
				for resourceName, resourceLimit := range container.Resources.Limits {
					limits = append(limits, fmt.Sprintf("%s=%s", resourceName, resourceLimit.String()))
				}
			}
			for _, volume := range controller.Spec.Template.Spec.Volumes {
				volumes = append(volumes, volume.Name)
			}
			for _, volumeMount := range controller.Spec.Template.Spec.Containers[0].VolumeMounts {
				mounts = append(mounts, volumeMount.Name)
			}
			for _, env := range controller.Spec.Template.Spec.Containers[0].Env {
				envs = append(envs, env.Name)
			}

			for key, value := range controller.Labels {
				labels = append(labels, fmt.Sprintf("%s=%s", key, value))
			}
			if len(controller.OwnerReferences) > 0 {
				controlleredBy := controller.OwnerReferences[0]
				controlleredByName = controlleredBy.Name
			}

		}

	} else if controllerType == "daemonset" {
		controller, err := kh.K8sClient.AppsV1().DaemonSets(namespace).Get(context.TODO(), controllerName, metav1.GetOptions{})
		if err != nil {
			return result, err
		}

		container := controller.Spec.Template.Spec.Containers

		if len(container) > 0 {
			for _, container := range controller.Spec.Template.Spec.Containers {
				for resourceName, resourceLimit := range container.Resources.Limits {
					limits = append(limits, fmt.Sprintf("%s=%s", resourceName, resourceLimit.String()))
				}
			}
			for _, volume := range controller.Spec.Template.Spec.Volumes {
				volumes = append(volumes, volume.Name)
			}
			for _, volumeMount := range controller.Spec.Template.Spec.Containers[0].VolumeMounts {
				mounts = append(mounts, volumeMount.Name)
			}
			for _, env := range controller.Spec.Template.Spec.Containers[0].Env {
				envs = append(envs, env.Name)
			}
			if len(controller.OwnerReferences) > 0 {
				controlleredBy := controller.OwnerReferences[0]
				controlleredByName = controlleredBy.Name
			}

		}

	} else if controllerType == "staefulset" {
		controller, err := kh.K8sClient.AppsV1().StatefulSets(namespace).Get(context.TODO(), controllerName, metav1.GetOptions{})
		if err != nil {
			return result, err
		}
		container := controller.Spec.Template.Spec.Containers

		if len(container) > 0 {
			for _, container := range controller.Spec.Template.Spec.Containers {
				for resourceName, resourceLimit := range container.Resources.Limits {
					limits = append(limits, fmt.Sprintf("%s=%s", resourceName, resourceLimit.String()))
				}
			}
			for _, volume := range controller.Spec.Template.Spec.Volumes {
				volumes = append(volumes, volume.Name)
			}
			for _, volumeMount := range controller.Spec.Template.Spec.Containers[0].VolumeMounts {
				mounts = append(mounts, volumeMount.Name)
			}
			for _, env := range controller.Spec.Template.Spec.Containers[0].Env {
				envs = append(envs, env.Name)
			}
			if len(controller.OwnerReferences) > 0 {
				controlleredBy := controller.OwnerReferences[0]
				controlleredByName = controlleredBy.Name
			}

		}

	} else if controllerType == "job" {
		controller, err := kh.K8sClient.BatchV1().Jobs(namespace).Get(context.TODO(), controllerName, metav1.GetOptions{})
		if err != nil {
			return result, err
		}
		container := controller.Spec.Template.Spec.Containers

		if len(container) > 0 {
			for _, container := range controller.Spec.Template.Spec.Containers {
				for resourceName, resourceLimit := range container.Resources.Limits {
					limits = append(limits, fmt.Sprintf("%s=%s", resourceName, resourceLimit.String()))
				}
			}
			for _, volume := range controller.Spec.Template.Spec.Volumes {
				volumes = append(volumes, volume.Name)
			}
			for _, volumeMount := range controller.Spec.Template.Spec.Containers[0].VolumeMounts {
				mounts = append(mounts, volumeMount.Name)
			}
			for _, env := range controller.Spec.Template.Spec.Containers[0].Env {
				envs = append(envs, env.Name)
			}
			if len(controller.OwnerReferences) > 0 {
				controlleredBy := controller.OwnerReferences[0]
				controlleredByName = controlleredBy.Name
			}

		}

	} else if controllerType == "cronjob" {
		controller, err := kh.K8sClient.BatchV1().CronJobs(namespace).Get(context.TODO(), controllerName, metav1.GetOptions{})
		if err != nil {
			return result, err
		}
		container := controller.Spec.JobTemplate.Spec.Template.Spec.Containers

		if len(container) > 0 {
			for _, container := range controller.Spec.JobTemplate.Spec.Template.Spec.Containers {
				for resourceName, resourceLimit := range container.Resources.Limits {
					limits = append(limits, fmt.Sprintf("%s=%s", resourceName, resourceLimit.String()))
				}
			}
			for _, volume := range controller.Spec.JobTemplate.Spec.Template.Spec.Volumes {
				volumes = append(volumes, volume.Name)
			}
			for _, volumeMount := range controller.Spec.JobTemplate.Spec.Template.Spec.Containers[0].VolumeMounts {
				mounts = append(mounts, volumeMount.Name)
			}
			for _, env := range controller.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Env {
				envs = append(envs, env.Name)
			}
			if len(controller.OwnerReferences) > 0 {
				controlleredBy := controller.OwnerReferences[0]
				controlleredByName = controlleredBy.Name
			}

		}

	} else if controllerType == "replicaset" {
		controller, err := kh.K8sClient.AppsV1().ReplicaSets(namespace).Get(context.TODO(), controllerName, metav1.GetOptions{})
		if err != nil {
			return result, err
		}
		container := controller.Spec.Template.Spec.Containers

		if len(container) > 0 {
			for _, container := range controller.Spec.Template.Spec.Containers {
				for resourceName, resourceLimit := range container.Resources.Limits {
					limits = append(limits, fmt.Sprintf("%s=%s", resourceName, resourceLimit.String()))
				}
			}
			for _, volume := range controller.Spec.Template.Spec.Volumes {
				volumes = append(volumes, volume.Name)
			}
			for _, volumeMount := range controller.Spec.Template.Spec.Containers[0].VolumeMounts {
				mounts = append(mounts, volumeMount.Name)
			}
			for _, env := range controller.Spec.Template.Spec.Containers[0].Env {
				envs = append(envs, env.Name)
			}
			if len(controller.OwnerReferences) > 0 {
				controlleredBy := controller.OwnerReferences[0]
				controlleredByName = controlleredBy.Name
			}

		}

	} else {
		err := fmt.Errorf("Invalid Controller Type %v", controllerType)
		return result, err
	}

	result.Limits = limits
	result.Environment = envs
	result.Mounts = mounts
	result.Volumes = volumes
	result.Labels = labels
	result.ControlledBy = controlleredByName

	return result, nil
}
