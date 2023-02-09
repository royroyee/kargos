package container

import (
	cm "Infra/common"
)

// GenerateContainersInfo will generate ContainersInfo from cm.ContainerInfo for gRPC communication.
func GenerateContainersInfo(info []cm.ContainerInfo) *ContainersInfo {
	var containers []*SingleContainerInfo
	for _, cont := range info {
		tmp := convertSingleContainerInfo(cont)
		containers = append(containers, tmp)
	}

	ret := ContainersInfo{Containers: containers}
	return &ret
}

// convertProcessInfo converts cm.Process into ProcessInfo for gRPC communication.
func convertProcessInfo(info cm.Process) *ProcessInfo {
	ret := ProcessInfo{
		Name:   info.Name,
		Status: info.Status,
		PID:    int32(info.PID),
		CPU:    float32(info.CPU),
		RAM:    float32(info.RAM),
	}

	return &ret
}

// convertSingleContainerInfo converts cm.ContainerInfo into a SingleContainerInfo for gRPC communication.
func convertSingleContainerInfo(info cm.ContainerInfo) *SingleContainerInfo {
	// Convert all cm.Process structs into slice of ProcessInfo pointers.
	var processes []*ProcessInfo
	for _, proc := range info.Processes {
		tmp := convertProcessInfo(proc)
		processes = append(processes, tmp)
	}

	ret := SingleContainerInfo{
		ID:        info.ID,
		Namespace: info.Namespace,
		Processes: processes,
	}

	return &ret
}
