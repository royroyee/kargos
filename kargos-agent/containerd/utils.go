package containerd

import (
	cm "Infra/common"
	"context"
	"fmt"
	"github.com/containerd/containerd/namespaces"
	"github.com/shirou/gopsutil/v3/process"
	"os"
)

// GetContainerPIDs will get all PIDs running in a specific container.
func GetContainerPIDs(info cm.ContainerInfo) []int32 {
	var ret []int32
	ctx := namespaces.WithNamespace(context.Background(), info.Namespace)
	task, err := info.Data.Task(ctx, nil)
	if err != nil {
		return ret
	}

	pids, err := task.Pids(ctx)
	if err != nil {
		return ret
	}

	for _, v := range pids {
		ret = append(ret, int32(v.Pid))
	}

	return ret
}

// PrintPIDs will print out all PID and commandline in human-readable format.
func PrintPIDs(pids []int32) {
	for i, pid := range pids {
		tmp, err := process.NewProcess(pid)
		if err != nil {
			continue
		}

		cmd, err := tmp.Cmdline()
		if err != nil {
			continue
		}

		fmt.Printf("    #%d %s\n", i, cmd)
	}
}

// PrintProcess will print out PID and its data according to the data by GenerateProcessInfo.
func PrintProcess(pids []int32) {
	for _, pid := range pids {
		tmp := GenerateProcessInfo(pid)
		fmt.Printf("    %d, %s, %f, %f (%s)\n", tmp.PID, tmp.Name, tmp.CPU, tmp.RAM, tmp.Status)
	}
}

// GenerateProcessInfo generates a Process struct using original PID.
func GenerateProcessInfo(pid int32) cm.Process {
	var ret cm.Process
	org, err := process.NewProcess(pid)
	if err != nil {
		return ret
	}

	// Store PID of process.
	ret.PID = uint32(pid)

	// Store status of process.
	status, err := org.Status()
	if err != nil {
		ret.Status = "Running"
	} else {
		ret.Status = status[0]
	}

	// Store name of process.
	cmd, err := org.Cmdline()
	if err != nil {
		ret.Name = "UNKNOWN"
	} else {
		ret.Name = cmd
	}

	// Stores how many percent of CPU time this process uses.
	cpu, err := org.CPUPercent()
	if err != nil {
		ret.CPU = -1.0
	} else {
		ret.CPU = cpu
	}

	// Stores how many percent of the total RAM this process uses.
	ram, err := org.MemoryPercent()
	if err != nil {
		ret.RAM = -1.0
	} else {
		ret.RAM = float64(ram)
	}

	return ret
}

// GetAllContainers will return data on all the container in this node.
func GetAllContainers() []cm.ContainerInfo {
	var ret []cm.ContainerInfo

	cdh := Handler{}
	err := cdh.InitHandler()
	if err != nil {
		fmt.Println(err)
		return ret
	}

	containers, err := cdh.GetContainers()
	if err != nil {
		fmt.Println(err)
		return ret
	}

	// Get all container and its PIDs.
	for _, cnt := range containers {
		pids := GetContainerPIDs(cnt)
		var procs []cm.Process

		for _, pid := range pids {
			tmp := GenerateProcessInfo(pid)
			procs = append(procs, tmp)
		}
		cnt.Processes = procs
		ret = append(ret, cnt)
	}

	cdh.StopHandler()
	return ret
}

// PrintContainers will print information on all container.
func PrintContainers(containers []cm.ContainerInfo) {
	for i, cnt := range containers {
		fmt.Printf("[%d] %s", i, cnt.ID)
		for _, proc := range cnt.Processes {
			fmt.Printf("    %d, %s, %f, %f (%s)\n", proc.PID, proc.Name, proc.CPU, proc.RAM, proc.Status)
		}
	}
}

// FindHostname finds out the name of the real host using /etc/hostname
func FindHostname() string {
	content, err := os.ReadFile("/etc/hostname")
	if err != nil {
		content = []byte("UNKNOWN")
	}
	return string(content)
}
