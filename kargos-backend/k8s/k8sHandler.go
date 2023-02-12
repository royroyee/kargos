package k8s

import (
	"context"
	cm "github.com/boanlab/kargos/common"
	"github.com/boanlab/kargos/k8s/kubectl"
	"gopkg.in/mgo.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"log"
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

// for Overview/main
func (kh K8sHandler) GetHome() (cm.Home, error) {
	var result cm.Home

	totalResources, namespaces, deployments, pods, ingresses, services, persistentVolumes, jobs, daemonSets, err := kh.GetTotalResources()

	if err != nil {
		return cm.Home{}, err
	}
	result = cm.Home{
		Version:    kh.GetKubeVersion(),
		TotalNodes: kh.GetTotalNodes(),
		Created:    kh.GetCreatedOfCluster(),
		Tabs: map[string]int{
			"TotalResources":   totalResources,
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

	return result, nil
}

// for nodes/overview
func (kh K8sHandler) GetNodeOverview() ([]cm.Node, error) {
	result, err := kh.GetNodeList()
	return result, err
}

// for node/:name
func (kh K8sHandler) GetNodeDetail(nodeName string) (cm.Node, error) {
	var result cm.Node

	node, err := kh.GetNode(nodeName)
	if err != nil {
		return cm.Node{}, err
	}

	cpuUsage, ramUsage, diskAllocated := kh.GetMetricUsage(*node)

	hours24, hours12, hours6 := kh.GetRecordOfNode(nodeName)

	podList, err := kh.GetPodsByNode(nodeName)

	if err != nil {
		return cm.Node{}, err
	}

	result = cm.Node{
		Name:          nodeName,
		CpuUsage:      cpuUsage,
		RamUsage:      ramUsage,
		DiskAllocated: diskAllocated,
		IP:            node.Status.Addresses[0].Address,
		Ready:         string(node.Status.Conditions[0].Status),
		OsImage:       node.Status.NodeInfo.OSImage,
		Pods:          kh.TransferPod(podList),
		Record: map[string]cm.RecordOfNode{
			"24hours": hours24,
			"12hours": hours12,
			"6hours":  hours6,
		},
	}
	return result, nil
}

// controllers/deployments/overview
func (kh K8sHandler) GetDeploymentOverview() ([]cm.Deployment, error) {

	var result []cm.Deployment

	deployList, err := kh.GetDeploymentList()
	if err != nil {
		return []cm.Deployment{}, err
	}

	for _, deploy := range deployList.Items {
		status, image := CheckContainerOfDeploy(deploy)
		result = append(result, cm.Deployment{
			Name:      deploy.GetName(),
			Namespace: deploy.GetNamespace(),
			Image:     image,
			Status:    status,
			Labels:    deploy.Labels,
			Created:   deploy.GetCreationTimestamp().String(),
		})
	}
	return result, nil
}

// GetDeploymentSpecific retrieves information of a deployment. This will also get details as well.
func (kh K8sHandler) GetDeploymentSpecific(namespace string, name string) (cm.Deployment, error) {
	ret := cm.Deployment{}

	deployment, err := kh.K8sClient.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return ret, err
	}

	ret.Name = deployment.GetName()
	ret.Namespace = deployment.GetNamespace()
	//ret.Image = deployment.Spec.Template.Spec.Containers[0].Image
	//ret.Status = string(deployment.Status.Conditions[0].Status)
	//ret.Labels = deployment.Labels
	//ret.Created = deployment.GetCreationTimestamp().String()
	ret.Details = generateDescribeString(name, ret.Namespace, "deployment")

	return ret, nil
}

// controlelrs/ingresses/overview
func (kh K8sHandler) GetIngressOverview() ([]cm.Ingress, error) {
	var result []cm.Ingress

	ingressList, err := kh.GetIngressList()
	if err != nil {
		return []cm.Ingress{}, err
	}
	for _, ingress := range ingressList.Items {
		result = append(result, cm.Ingress{
			Name:      ingress.GetName(),
			Namespace: ingress.GetNamespace(),
			Labels:    ingress.Labels,
			Host:      ingress.Spec.Rules[0].Host,
			Class:     ingress.Spec.IngressClassName,
			Address:   ingress.Status.LoadBalancer.Ingress[0].IP,
			Created:   ingress.GetCreationTimestamp().String(),
		})
	}

	return result, nil
}

// controllers/ingress/:namespace/:name
func (kh K8sHandler) GetIngressSpecific(name string, namespace string) (cm.Ingress, error) {
	ret := cm.Ingress{}

	ingress, err := kh.K8sClient.NetworkingV1().Ingresses(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return ret, err
	}

	ret.Name = ingress.GetName()
	ret.Namespace = ingress.GetNamespace()

	ret.Details = generateDescribeString(name, ret.Namespace, "ingress")

	return ret, nil
}

// controllers/jobs/overview
func (kh K8sHandler) GetJobsOverview() ([]cm.Job, error) {
	var result []cm.Job
	jobList, err := kh.GetJobsList()
	if err != nil {
		return []cm.Job{}, err
	}

	for _, job := range jobList.Items {
		result = append(result, cm.Job{
			Name:      job.Name,
			Namespace: job.Namespace,
			Failed:    job.Status.Failed,
			Succeeded: job.Status.Succeeded,
			Created:   job.CreationTimestamp.Time.String(),
		})
	}
	return result, nil
}

// controllers/job/:namespace/:name
func (kh K8sHandler) GetJobSpecific(namespace string, name string) (cm.Job, error) {
	ret := cm.Job{}

	job, err := kh.K8sClient.BatchV1().Jobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return ret, err
	}

	ret.Name = job.GetName()
	ret.Namespace = job.GetNamespace()

	ret.Details = generateDescribeString(name, ret.Namespace, "job")

	return ret, nil
}

// controllers/daemonsets/overview
func (kh K8sHandler) GetDaemonSetsOverview() ([]cm.DaemonSet, error) {
	var result []cm.DaemonSet
	daemonSetList, err := kh.GetDaemonSetList()
	if err != nil {
		return []cm.DaemonSet{}, err
	}

	for _, daemonSet := range daemonSetList.Items {
		result = append(result, cm.DaemonSet{
			Name:           daemonSet.GetName(),
			Namespace:      daemonSet.GetNamespace(),
			Labels:         daemonSet.Labels,
			UpdateStrategy: string(daemonSet.Spec.UpdateStrategy.Type),
			Created:        daemonSet.CreationTimestamp.Time.String(),
		})
	}
	return result, nil
}

// controllers/daemonset/:namespace/:name
func (kh K8sHandler) GetDaemonSetSpecific(namespace string, name string) (cm.DaemonSet, error) {
	ret := cm.DaemonSet{}

	daemonset, err := kh.K8sClient.AppsV1().DaemonSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return ret, err
	}

	ret.Name = daemonset.GetName()
	ret.Namespace = daemonset.GetNamespace()

	ret.Details = generateDescribeString(name, ret.Namespace, "daemonset")

	return ret, nil
}

// resources/namespaces/overview
func (kh K8sHandler) GetNamespaceOverview() ([]cm.Namespace, error) {
	var result []cm.Namespace

	namespaceList, err := kh.GetNamespaceList()
	if err != nil {
		return []cm.Namespace{}, err
	}

	for _, namespace := range namespaceList.Items {
		result = append(result, cm.Namespace{
			Name:   namespace.GetName(),
			Labels: namespace.Labels,
			Status: string(namespace.Status.Phase),
		})
	}
	return result, nil
}

// resources/namespaces/overview
func (kh K8sHandler) GetNamespaceDetail(name string) (cm.Namespace, error) {
	var result cm.Namespace

	namespace, err := kh.GetNamespaceByName(name)
	if err != nil {
		return result, err
	}

	result = cm.Namespace{
		Name:        namespace.GetName(),
		Labels:      namespace.Labels,
		Status:      string(namespace.Status.Phase),
		Annotations: namespace.Annotations,
		Finalizers:  namespace.Finalizers,
		Created:     namespace.CreationTimestamp.Time.String(),
	}
	return result, nil
}

// resources/pods/overview
func (kh K8sHandler) GetPodOverview() ([]cm.Pod, error) {
	var result []cm.Pod

	podList, err := kh.GetPodList()
	if err != nil {
		return []cm.Pod{}, err
	}

	for _, pod := range podList.Items {
		// Find Container's name
		var containerNames []string
		containerStats := pod.Status.ContainerStatuses
		for _, containerStat := range containerStats {
			containerNames = append(containerNames, containerStat.ContainerID)
		}

		result = append(result, cm.Pod{
			Name:             pod.GetName(),
			Namespace:        pod.GetNamespace(),
			PodIP:            pod.Status.PodIP,
			Status:           string(pod.Status.Phase),
			ServiceConnected: pod.Spec.EnableServiceLinks,
			Restarts:         GetRestartCount(pod),
			Image:            CheckContainerOfPod(pod),
			Age:              pod.CreationTimestamp.String(),
			ContainerNames:   containerNames,
			Timestamp:        time.Now(), // not pod's creation time , just for db query
		})
	}
	return result, nil
}

// resources/pod/:name
// for node/:name
func (kh K8sHandler) GetPodDetail(podName string) (cm.Pod, error) {

	result, err := kh.GetRecordOfPod(podName)
	if err != nil {
		return result, err
	}

	//pod, err := kh.GetPodByName(namespace, podName)
	//if err != nil {
	//	return cm.Pod{}, errs
	//}
	//result = cm.Pod{
	//	Name:             pod.GetName(),
	//	Namespace:        pod.GetNamespace(),
	//	PodIP:            pod.Status.PodIP,
	//	Status:           string(pod.Status.Phase),
	//	ServiceConnected: pod.Spec.EnableServiceLinks,
	//	Restarts:         GetRestartCount(*pod),
	//	Image:            pod.Status.ContainerStatuses[0].Image,
	//	Age:              pod.CreationTimestamp.String(),
	//}
	//return result, nil

	return result, nil
}

func (kh K8sHandler) GetServiceOverview() ([]cm.Service, error) {

	var result []cm.Service
	services, err := kh.GetServiceList()
	if err != nil {
		return result, err
	}
	for _, service := range services.Items {

		result = append(result, cm.Service{
			Name:       service.GetName(),
			Namespace:  service.GetNamespace(),
			Type:       string(service.Spec.Type),
			ClusterIP:  service.Spec.ClusterIP,
			ExternalIP: service.Spec.ExternalName,
			Port:       service.Spec.Ports[0].Port,
			NodePort:   service.Spec.Ports[0].NodePort,
		})
	}
	return result, err
}

func (kh K8sHandler) GetServiceDetail(namespace string, name string) (cm.Service, error) {
	var result cm.Service

	service, err := kh.GetServiceByName(namespace, name)
	if err != nil {
		return result, err
	}

	result = cm.Service{
		Name:       service.GetName(),
		Namespace:  service.GetNamespace(),
		Type:       string(service.Spec.Type),
		ClusterIP:  service.Spec.ClusterIP,
		ExternalIP: service.Spec.ExternalName,
		Port:       service.Spec.Ports[0].Port,
		NodePort:   service.Spec.Ports[0].NodePort,
		Selector:   service.Spec.Selector,
		Conditions: service.Status.Conditions,
		Labels:     service.Labels,
		Created:    service.CreationTimestamp.Time.String(),
	}

	return result, err
}

func (kh K8sHandler) GetPersistentVolumeOverview() ([]cm.PersistentVolume, error) {
	var result []cm.PersistentVolume

	pvs, err := kh.GetPersistentVolumeList()
	if err != nil {
		return result, err
	}

	for _, pv := range pvs.Items {
		result = append(result, cm.PersistentVolume{
			Name:          pv.GetName(),
			Capacity:      pv.Spec.Capacity,
			AccessModes:   pv.Spec.AccessModes,
			ReclaimPolicy: pv.Spec.PersistentVolumeReclaimPolicy,
			Status:        string(pv.Status.Phase),
			Claim:         GetPersistentVolumeClaim(&pv),
			StorageClass:  pv.Spec.StorageClassName,
		})
	}
	return result, err
}

func (kh K8sHandler) GetPersistentVolumeDetail(name string) (cm.PersistentVolume, error) {
	var result cm.PersistentVolume

	pv, err := kh.GetPersistentVolumeByName(name)
	if err != nil {
		return result, err
	}

	result = cm.PersistentVolume{
		Name:          pv.GetName(),
		Capacity:      pv.Spec.Capacity,
		AccessModes:   pv.Spec.AccessModes,
		ReclaimPolicy: pv.Spec.PersistentVolumeReclaimPolicy,
		Status:        string(pv.Status.Phase),
		Claim:         GetPersistentVolumeClaim(pv),
		StorageClass:  pv.Spec.StorageClassName,
		Reason:        pv.Status.Reason,
		MountOption:   pv.Spec.MountOptions,
		Labels:        pv.Labels,
		Created:       pv.CreationTimestamp.Time.String(),
	}

	return result, err
}

// generateDescribeString generates string that represent kubernetes resource like "kubectl describe"
// The code originated from kubectl source code's kubectl/pkg/cmd/cmd.go
func generateDescribeString(name string, namespace string, resourceType string) string {
	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag().WithDiscoveryBurst(300).WithDiscoveryQPS(50.0)
	cmdutil.NewMatchVersionFlags(kubeConfigFlags)
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)
	f := cmdutil.NewFactory(matchVersionKubeConfigFlags)
	flags := kubectl.NewDescribeFlags(f, genericclioptions.IOStreams{})
	o, _ := flags.ToOptions("kubectl", []string{resourceType, name, "namespace", namespace})
	ret := o.Run()
	return ret
}

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
