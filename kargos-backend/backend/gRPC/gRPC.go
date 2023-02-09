package gRPC

import (
	"github.com/boanlab/kargos/gRPC/protos/v3/container"
	"github.com/boanlab/kargos/k8s"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

// Handler struct stores data for gRPC server.
type Handler struct {
	gRPCServer  *grpc.Server
	K8sHandler  *k8s.K8sHandler
	initialized bool
}

// NewGRPCHandler generates a new handler struct and returns the handler.
func NewGRPCHandler(K8sHandler *k8s.K8sHandler) *Handler {
	ret := Handler{
		initialized: false,
		K8sHandler:  K8sHandler,
	}

	return &ret
}

// registerServices registers services of all gRPC services.
func (gh Handler) registerServices() {
	container.RegisterContainersServer(gh.gRPCServer, &container.Containers{gh.K8sHandler})
}

// StartGRPCServer initializes gRPC server and starts the server.
func (gh Handler) StartGRPCServer() {
	// If environment variable GRPC_LISTEN_PORT was not set, just use 50001 as default port.
	port := os.Getenv("GRPC_LISTEN_PORT")
	if len(port) == 0 {
		port = "50001"
	}

	// Prepare listening gRPC server.
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Printf("could not listen gRPC server: %s\n", err)
		panic(err)
	}

	// Generate new gRPC server and register services.
	gh.gRPCServer = grpc.NewServer()
	gh.registerServices()
	log.Printf("starting gRPC server in :%s\n", port)

	// Start serving gRPC.
	err = gh.gRPCServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
