package containerd

import (
	cm "Infra/common"
	"context"
	"errors"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"os"
)

// Handler is a handler for Containerd actions.
type Handler struct {
	client      *containerd.Client
	initialized bool
}

// InitHandler will initialize Containerd client.
func (ch *Handler) InitHandler() error {
	var err error
	// Containerd socket will be set as /run/containerd/containerd.sock by default.
	// For custom usage, use CONTAINERD_SOCK env variable as your containerd.sock.
	sock := os.Getenv("CONTAINERD_SOCK")
	if len(sock) == 0 {
		sock = "/run/containerd/containerd.sock"
	}

	ch.client, err = containerd.New(sock)
	if err != nil {
		return err
	} else {
		ch.initialized = true
		return nil
	}
}

// StopHandler will destruct Containerd client.
func (ch *Handler) StopHandler() {
	if ch.initialized {
		defer func(client *containerd.Client) {
			_ = client.Close()
		}(ch.client)
	}
}

// GetContainers will get all container that were created by Containerd in all namespaces. Thanks chatGPT :)
func (ch *Handler) GetContainers() ([]cm.ContainerInfo, error) {
	var ret []cm.ContainerInfo
	if ch.initialized {
		nss, err := ch.client.NamespaceService().List(context.Background())
		if err != nil {
			return ret, err
		}

		// Iterate over the namespaces
		for _, ns := range nss {
			// Use the default namespace
			ctx := namespaces.WithNamespace(context.Background(), ns)

			// Get a list of container
			containers, err := ch.client.Containers(ctx)
			if err != nil {
				return ret, err
			}

			// Iterate over container and generate ContainerInfo
			for _, container := range containers {
				tmp := cm.ContainerInfo{ID: container.ID(), Data: container, Namespace: ns}
				ret = append(ret, tmp)
			}
		}
		return ret, nil
	} else {
		return ret, errors.New("handler was not initialized")
	}
}
