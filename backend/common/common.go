package common

// Overview main
type Home struct {
	Version    string // kuberentes version
	TotalNodes int    // total nodes
	Created    string // created

	Tabs map[string]int // total_resources ~ daemon_sets

	TopNamespaces []string
	AlertCount    int // warning 등의 이벤트만
}

// Alert
type Alert struct {
	tag     string
	message string
	uuid    string
}

// Node
type Node struct {
	Name      string
	CpuUsage  float64
	RamUsage  float64
	DiskUsage float64
	IP        string
	OsImage   string
	Pods      []Pod

	// detail 항목들
}

// Pod
type Pod struct {
	Name             string
	Status           string // Running  or Pending
	Image            string
	Age              string // created
	Namespace        string `json:"namespace"`
	PodIP            string `json:"podIP"`
	ServiceConnected *bool  `json:"serviceConnected"`
	Restarts         int    `json:"restarts"`
}

// Deployment
type Deployment struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Image     string
	Status    string
	Label     string
	created   string `json:"pod_count"`

	// detail 항목들
}

// Namespace
type Namespace struct {
	Name      string
	Status    string
	Cpu_usage string
	Ram_usage string

	// Infra agent
	process []Process // inner struct
}

// Process (Infra agent)
type Process struct {
	Name     string
	Status   string
	PID      int32
	CpuUsage int32
	RamUsage int32
}
