package common

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

//// Overview main
//type Overview struct {
//	Version    string `json:"kubernetes_version"` // kubernetes version
//	TotalNodes int    `json:"total_nodes"`        // total nodes
//	Created    string `json:"created"`            // created
//
//	Tabs map[string]int // total_resources ~ daemon_sets
//
//	TopNamespaces []string `json:"top_namespaces"`
//	AlertCount    int      `json:"alert_count"` // warning 등의 이벤트만
//}

// Overview main
type Overview struct {
	Version    string     `json:"kubernetes_version"` // kubernetes version
	NodeStatus NodeStatus `json:"node_status"`
	PodStatus  PodStatus  `json:"pod_status"`
}

// Node
type Node struct {
	Name          string                  `json:"name"`
	CpuUsage      float64                 `json:"cpu_usage"`
	RamUsage      float64                 `json:"ram_usage"`
	DiskAllocated float64                 `json:"disk_allocated"`
	IP            string                  `json:"ip"`
	Ready         string                  `json:"ready"`
	OsImage       string                  `json:"os_image"`
	Pods          []Pod                   `json:"pods"`
	Record        map[string]RecordOfNode `json:"record"`
}

// NodeMetric (DB ( last 24 hours etc ..)
type RecordOfNode struct {
	Name          string    `json:"name"`
	CpuUsage      float64   `json:"cpu_usage"`
	RamUsage      float64   `json:"ram_usage"`
	DiskAllocated float64   `json:"disk_allocated"`
	Timestamp     time.Time `json:"timestamp"`
}

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

// Persistent Volume
type PersistentVolume struct {
	Name          string                           `json:"name"`
	Capacity      v1.ResourceList                  `json:"capacity"`
	AccessModes   []v1.PersistentVolumeAccessMode  `json:"access_modes"`
	ReclaimPolicy v1.PersistentVolumeReclaimPolicy `json:"reclaim_policy"`
	Status        string                           `json:"status"`
	Claim         string                           `json:"claim"`
	StorageClass  string                           `json:"storage_class"`
	Reason        string                           `json:"reason"`
	MountOption   []string                         `json:"mount_option"`
	Labels        map[string]string                `json:"labels"`
	Created       string                           `json:"created"`
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
	Created string `json:"created"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type NodeStatus struct {
	NotReady int `json:"not_ready""`
	Ready    int `json:"ready"`
}

type PodStatus struct {
	Error   int `json:"error"`
	Pending int `json:"pending"`
	Running int `json:"running"`
}

// Controllers
type Controllers struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Namespace   string   `json:"namespace"`
	Replicasets int      `json:"replicasets"`
	Pods        []string `json:"pods"`
}

// Pod
type Pod struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	CpuUsage  int64  `json:"cpuUsage"`
	RamUsage  int64  `json:"ramUsage"`
	Restarts  int32  `json:"restarts"`
	PodIP     string `json:"pod_ip"`
	Status    string `json:"status"`
	Image     string `json:"image""`

	// Container struct
	Containers     []Container `json:"containers"`
	ContainerNames []string
}

type PodOverview struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	CpuUsage  int64  `json:"cpuUsage"`
	RamUsage  int64  `json:"ramUsage"`
	Restarts  int32  `json:"restarts"`
	PodIP     string `json:"pod_ip"`
	Status    string `json:"status"`
}
