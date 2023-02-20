package common

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Overview main
type Overview struct {
	NodeStatus NodeStatus `json:"node_status"`
	PodStatus  PodStatus  `json:"pod_status"`
}

//// Node
//type Node struct {
//	Name          string                  `json:"name"`
//	CpuUsage      float64                 `json:"cpu_usage"`
//	RamUsage      float64                 `json:"ram_usage"`
//	DiskAllocated float64                 `json:"disk_allocated"`
//	IP            string                  `json:"ip"`
//	Ready         string                  `json:"ready"`
//	OsImage       string                  `json:"os_image"`
//	Pods          []Pod                   `json:"pods"`
//	Record        map[string]RecordOfNode `json:"record"`
//}
//
//// NodeMetric (DB ( last 24 hours etc ..)
//type RecordOfNode struct {
//	Name          string    `json:"name"`
//	CpuUsage      float64   `json:"cpu_usage"`
//	RamUsage      float64   `json:"ram_usage"`
//	DiskAllocated float64   `json:"disk_allocated"`
//	Timestamp     time.Time `json:"timestamp"`
//}

//// Pod
//type Pod struct {
//	Name             string    `json:"name"`
//	Namespace        string    `json:"namespace"`
//	PodIP            string    `json:"pod_ip"`
//	Status           string    `json:"status"` // Running  or Pending
//	ServiceConnected *bool     `json:"service_connected"`
//	Restarts         int32     `json:"restarts"`
//	Image            string    `json:"image"`
//	Age              string    `json:"age"`
//	Timestamp        time.Time `json:"timestamp"` // not pod's created , just for db query
//
//	// Container struct
//	Containers     []Container `json:"containers"`
//	ContainerNames []string
//}

// Deployment
type Deployment struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Image     string            `json:"image"`
	Status    string            `json:"status"`
	Labels    map[string]string `json:"label"`
	Created   string            `json:"created"`

	// detail
	Details string `json:"details"`
}

// Ingress
type Ingress struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
	Host      string            `json:"host"`
	Class     *string           `json:"class"`
	Address   string            `json:"address"`
	Created   string            `json:"created"`

	Details string `json:"details"`
}

// Namespace
type Namespace struct {
	Name        string            `json:"name"`
	Labels      map[string]string `json:"labels"`
	Status      string            `json:"status"`
	Annotations map[string]string `json:"annotations"`
	Finalizers  []string          `json:"finalizers"`
	Created     string            `json:"created"`

	// Infra agent
	process []Process `json:"process"` // inner struct
}

// Service
type Service struct {
	Name       string             `json:"name"`
	Namespace  string             `json:"namespace"`
	Type       string             `json:"Type"`
	ClusterIP  string             `json:"cluster_ip"`
	ExternalIP string             `json:"external_ip"`
	Port       int32              `json:"port"`
	NodePort   int32              `json:"node_port"`
	Selector   map[string]string  `json:"selector"`
	Conditions []metav1.Condition `json:"conditions"`
	Labels     map[string]string  `json:"labels"`
	Created    string             `json:"created"`
}

// Job
type Job struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Failed    int32  `json:"failed"`
	Succeeded int32  `json:"succeeded"`
	Created   string `json:"created"`

	Details string `json:"details"`
}

// DaemonSet
type DaemonSet struct {
	Name           string            `json:"name"`
	Namespace      string            `json:"namespace "`
	Labels         map[string]string `json:"labels"`
	UpdateStrategy string            `json:"update_strategy"`
	Created        string            `json:"created"`

	Details string `json:"details"`
}

// Process (Infra agent)
type Process struct {
	Name     string  `json:"name"`
	Status   string  `json:"status"`
	PID      int32   `json:"pid"`
	CpuUsage float32 `json:"cpu_usage"`
	RamUsage float32 `json:"ram_usage"`
}

// Container stores data for a single container.
type Container struct {
	ID        string    `json:"name"`
	Namespace string    `json:"image"`
	Processes []Process `json:"processes"`
}

// 02.11 ~ //

// Event
type Event struct {
	Created    string `json:"created"`
	EventLevel string `json:"event_level"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Type       string `json:"type"`
}

type NodeStatus struct {
	NotReady []string `json:"not_ready""`
	Ready    []string `json:"ready"`
}

type PodStatus struct {
	Error   []string `json:"error"`
	Pending []string `json:"pending"`
	Running int      `json:"running"`
}

// Controller
type Controller struct {
	Namespace string   `json:"namespace"`
	Type      string   `json:"type"`
	Name      string   `json:"name"`
	Pods      []string `json:"pods"`
	Volumes   []string `json:"volumes"`
}

type ControllerOverview struct {
	Namespace string   `json:"namespace"`
	Type      string   `json:"type"`
	Name      string   `json:"name"`
	Pods      []string `json:"pods"`
}

// Pod
type Pod struct {
	Name       string   `json:"name"`
	Namespace  string   `json:"namespace"`
	CpuUsage   int64    `json:"cpu_usage"`
	RamUsage   int64    `json:"ram_usage"`
	Restarts   int32    `json:"restarts"`
	PodIP      string   `json:"pod_ip"`
	Node       string   `json:"node"`
	Volumes    []string `json:"volumes"`
	Controller string   `json:"controller"`
	Status     string   `json:"status"`
	Image      string   `json:"image"`
	Timestamp  string   `json:"timestamp"`

	//ControllerKind string `json:"controller_kind"`
	//ControllerName string `json:"controller_name"`

	// Container struct
	Containers     []Container `json:"containers"`
	ContainerNames []string    `json:"container_names"`
}

type PodUsage struct {
	Name     string `json:"name"`
	CpuUsage int64  `json:"cpu_usage"`
	RamUsage int64  `json:"ram_usage"`
	// TODO Network, Disk Usage
	Timestamp string `json:"timestamp"`

	// Container struct
	Containers     []Container `json:"containers"`
	ContainerNames []string
}
type GetPodUsage struct {
	CpuUsage     []int `json:"cpu_usage"`
	RamUsage     []int `json:"ram_usage"`
	NetworkUsage []int `json:"network_usage"`
	DiskUsage    []int `json:"disk_usage"`
}

type PodInfo struct {
	Name       string   `json:"name"`
	Namespace  string   `json:"namespace"`
	Image      string   `json:"image"`
	Node       string   `json:"node"`
	PodIP      string   `json:"pod_ip"`
	Restarts   int32    `json:"restarts"`
	Volumes    []string `json:"volumes"`
	Controller string   `json:"controller"`
	Status     string   `json:"status"`
}

type PodsOfController struct {
	Pods []string `json:"pods"`
}

type Node struct {
	Name          string  `json:"name"`
	CpuUsage      float64 `json:"cpu_usage"`
	RamUsage      float64 `json:"ram_usage"`
	DiskAllocated float64 `json:"disk_allocated"`
	NetworkUsage  float64 `json:"network_usage"`
	IP            string  `json:"ip"`
	Status        string  `json:"status"`
	Timestamp     string  `json:"timestamp"`
}

type NodeOverview struct {
	Name          string  `json:"name"`
	CpuUsage      float64 `json:"cpu_usage"`
	RamUsage      float64 `json:"ram_usage"`
	DiskAllocated float64 `json:"disk_allocated"`
	NetworkUsage  float64 `json:"network_usage"`
	IP            string  `json:"ip"`
	Status        string  `json:"status"`
}
type NodeUsage struct {
	CpuUsage     []int `json:"cpu_usage"`
	RamUsage     []int `json:"ram_usage"`
	NetworkUsage []int `json:"network_usage"`
	DiskUsage    []int `json:"disk_usage"`
}

type NodeCpuUsage struct {
	Name     string `json:"name"`
	CpuUsage int    `json:"cpu_usage"`
}
type NodeRamUsage struct {
	Name     string `json:"name"`
	RamUsage int    `json:"ram_usage"`
}
type TopNode struct {
	Cpu []NodeCpuUsage `json:"cpu"`
	Ram []NodeRamUsage `json:"ram"`
}

type PodCpuUsage struct {
	Name     string `json:"name"`
	CpuUsage int    `json:"cpu_usage"`
}
type PodRamUsage struct {
	Name     string `json:"name"`
	RamUsage int    `json:"ram_usage"`
}
type TopPod struct {
	Cpu []PodCpuUsage `json:"cpu"`
	Ram []PodRamUsage `json:"ram"`
}

type NodeInfo struct {
	OS                      string `json:"os"`
	HostName                string `json:"host_name"`
	IP                      string `json:"ip"`
	KubeletVersion          string `json:"kubelet_version"`
	ContainerRuntimeVersion string `json:"container_runtime_version"`
	NumContainers           int    `json:"num_containers"`
	CpuCores                int64  `json:"cpu_cores"`
	RamCapacity             int64  `json:"ram_capacity"`
	Status                  bool   `json:"status"`
}

type VolumeInfo struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	ClaimName string `json:"claim_name"`
	ReadOnly  string `json:"read_only"`
}

type PersistentVolume struct {
	Name         string                          `json:"name"`
	Capacity     int64                           `json:"capacity"`
	AccessModes  []v1.PersistentVolumeAccessMode `json:"access_modes"`
	Claim        string                          `json:"claim"`
	StorageClass string                          `json:"storage_class"`
	Status       string                          `json:"status"`
}

type Count struct {
	Count int `json:"count"`
}

type ControllerInfo struct {
	Labels       []string `json:"labels"`
	Limits       []string `json:"limits"`
	Environment  []string `json:"environment"`
	Mounts       []string `json:"mounts"`
	Volumes      []string `json:"volumes"`
	ControlledBy string   `json:"controlled_by"`
}

type Containers struct {
	ContainerNames []string `json:"container_names"`
}
