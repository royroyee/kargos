package main

import (
	"github.com/boanlab/kargos/backend/http"
	"github.com/boanlab/kargos/backend/k8s"
)

type Handlers struct {
	k8sHandler  k8s.K8sHandler
	httpHandler http.HTTPHandler

	// TODO gRPCHandler
}

var k8sHandler k8s.K8sHandler
var httpHandler http.HTTPHandler

func InitHandlers() {
	k8sHandler = *k8s.NewK8sHandler()
	httpHandler = *http.NewHTTPHandler()

	// TODO gRPC Server
}

func main() {

	// Start HTTP Servers
	httpHandler.StartHTTPServer()

	// TODO gRPC Server
}
