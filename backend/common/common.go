package common

// Overview main
type Home struct {
	Version    string `json:"kubernetes_version"` // kubernetes version
	TotalNodes int    `json:"total_nodes"`        // total nodes
	Created    string `json:"created"`            // created

	Tabs map[string]int // total_resources ~ daemon_sets

	TopNamespaces []string `json:"top_namespaces"`
	AlertCount    int      `json:"alert_count"` // warning 등의 이벤트만
}

// Alert
type Alert struct {
	tag     string `json:"tag"`
	message string `json:"message"`
	uuid    string `json:"uuid"`
}

// Node
type Node struct {
	Name      string  `json:"name"`
	CpuUsage  float64 `json:"cpu_usage"`
	RamUsage  float64 `json:"ram_usage"`
	DiskUsage float64 `json:"disk_usage"`
	IP        string  `json:"ip"`
	Ready     string  `json:"ready"`
	OsImage   string  `json:"os_image"`
	Pods      []Pod   `json:"pods"`

	// detail 항목들
}

// Pod
type Pod struct {
	Name             string `json:"name"`
	Status           string `json:"status"` // Running  or Pending
	Image            string `json:"image"`
	Age              string `json:"age"` // created
	Namespace        string `json:"namespace"`
	PodIP            string `json:"pod_ip"`
	ServiceConnected *bool  `json:"service_connected"`
	Restarts         int    `json:"restarts"`
}

// Deployment
type Deployment struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Image     string `json:"image"`
	Status    string `json:"status"`
	Label     string `json:"label"`
	created   string `json:"pod_count"`

	// detail 항목들
}

// Namespace
type Namespace struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	CpuUsage string `json:"cpu_usage"`
	RamUsage string `json:"ram_usage"`

	// Infra agent
	process []Process `json:"process"` // inner struct
}

// Process (Infra agent)
type Process struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	PID      int32  `json:"pid"`
	CpuUsage int32  `json:"cpu_usage"`
	RamUsage int32  `json:"ram_usage"`
}
