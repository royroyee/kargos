package container

import (
	"context"
	"fmt"
	cm "github.com/boanlab/kargos/common"
	"github.com/boanlab/kargos/k8s"
	"log"
	"strings"
)

type Containers struct {
	K8sHandler *k8s.K8sHandler
}

// translateContainersInfo translates ContainersInfo into slice of cm.ContainerInfo
func translateContainersInfo(info *ContainersInfo) []cm.Container {
	var ret []cm.Container
	for _, container := range info.Containers {
		tmp := translateSingleContainerInfo(container)
		ret = append(ret, tmp)
	}

	return ret
}

// translateSingleContainerInfo translates SingleContainerInfo into cm.ContainerInfo
func translateSingleContainerInfo(info *SingleContainerInfo) cm.Container {
	// Translate all processes
	var processes []cm.Process
	for _, process := range info.Processes {
		tmp := translateProcessesInfo(process)
		processes = append(processes, tmp)
	}

	// Store other data
	ret := cm.Container{
		ID:        info.ID,
		Namespace: info.Namespace,
		Processes: processes,
	}

	return ret
}

// translateProcessInfo translates ProcessInfo into cm.Process struct.
func translateProcessesInfo(info *ProcessInfo) cm.Process {
	ret := cm.Process{
		Name:     info.Name,
		Status:   info.Status,
		PID:      info.PID,
		CpuUsage: info.CPU,
		RamUsage: info.RAM,
	}

	return ret
}

// printContainers will print information on all container.
func printContainers(containers []cm.Container) {
	for i, cnt := range containers {
		fmt.Printf("[%d] %s", i, cnt.ID)
		for _, proc := range cnt.Processes {
			fmt.Printf("    %d, %s, %f, %f (%s)\n", proc.PID, proc.Name, proc.CpuUsage, proc.RamUsage, proc.Status)
		}
	}
}

// matchPodContainers finds out the pod from container.
func matchPodContainers(containers []cm.Container, pod *cm.PodUsage) {
	// Look for containers that matches the pod's container.
	// This is O(n^2) operation, therefore, this must be optimized.
	for _, podContainerName := range pod.ContainerNames {
		for _, container := range containers {
			if strings.Contains(podContainerName, container.ID) {
				pod.Containers = append(pod.Containers, container)
			}
		}
	}
}

// SendContainerData receives container is a callback function for gRPC with service of Containers.
func (c Containers) SendContainerData(ctx context.Context, info *ContainersInfo) (*Response, error) {
	containers := translateContainersInfo(info)
	//printContainers(containers) (debug)

	//	pods, _ := c.K8sHandler.PodOverview()
	pods, _ := c.K8sHandler.GetPodUsage()
	for i, pod := range pods {
		matchPodContainers(containers, &pod)
		pods[i] = pod
		fmt.Println("SendContainerData")
		printContainers(pods[i].Containers)
		fmt.Println()
	}

	// Store Pod Data into DB
	c.K8sHandler.StorePodUsageInDB(pods)

	log.Printf("received data from agent %s\n", info.NodeInfo)
	return &Response{
		Status: 100,
	}, nil
}

func (c Containers) mustEmbedUnimplementedContainersServer() {
	//TODO implement me
	panic("implement me")
}
