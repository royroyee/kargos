package grpc

import (
	cd "Infra/containerd"
	pbc "Infra/grpc/protos/v3/container"
	"context"
	"errors"
	"google.golang.org/grpc"
	"log"
	"os"
)

// Handler is for storing information on connection status and initialization status.
type Handler struct {
	conn        *grpc.ClientConn
	initialized bool
}

// NewHandler generate
func NewHandler() *Handler {
	ret := Handler{
		initialized: false,
	}

	return &ret
}

// SendContainerInfo sends container information to the gRPC server that is in backend server.
func (gh *Handler) SendContainerInfo() error {
	if gh.initialized {
		client := pbc.NewContainersClient(gh.conn)
		containers := cd.GetAllContainers()
		data := pbc.GenerateContainersInfo(containers)

		data.NodeInfo = cd.FindHostname()
		response, err := client.SendContainerData(context.Background(), data)
		if err != nil {
			log.Printf("Error sending data: %s", err)
			return err
		}

		// With response 100, everything went correct.
		switch response.Status {
		case 100:
			{
				return nil
			}
		default:
			return nil
		}

	} else {
		return errors.New("handler was not initialized before")
	}
}

// InitializeHandler will initialize handler for gRPC communication.
func (gh *Handler) InitializeHandler() {
	serverIP := os.Getenv("SERVER_IP")
	serverPort := os.Getenv("SERVER_PORT")

	log.Printf("server: %s:%s\n", serverIP, serverPort)
	// Environment variable SERVER_IP or SERVER_PORT was not set.
	if len(serverIP) == 0 || len(serverPort) == 0 {
		log.Printf("ENV SERVER_IP or SERVER_PORT is not set.")
		return
	}

	// Try dialing gRPC server which is implemented in the backend.
	var err error
	gh.conn, err = grpc.Dial(serverIP+":"+serverPort, grpc.WithInsecure())
	if err != nil {
		log.Printf("invalid server IP or port : %s:%s", serverIP, serverPort)
		gh.initialized = false
		return
	} else {
		gh.initialized = true
	}
}
