package common

import "github.com/containerd/containerd"

// Process stores data for a single process including CPU and RAM usage.
type Process struct {
	Name   string
	Status string
	PID    uint32
	CPU    float64
	RAM    float64
}

// ContainerInfo stores data for a single container.
type ContainerInfo struct {
	ID        string
	Namespace string
	Data      containerd.Container
	Processes []Process
}
