package main

import (
	"github.com/boanlab/kargos/http"
	"github.com/boanlab/kargos/k8s"
	"log"
)

var handlers Handlers

type Handlers struct {
	k8sHandler  *k8s.K8sHandler
	httpHandler *http.HTTPHandler

	// TODO gRPCHandler & dbHandler
}

func initHandlers() {

	handlers.k8sHandler = k8s.NewK8sHandler()
	handlers.httpHandler = http.NewHTTPHandler(handlers.k8sHandler)

	// TODO gRPC Server & DB
}

func init() {
	log.SetPrefix("Kargos: ")
}

func main() {
	// Handlers
	initHandlers()

	log.Println("Welcome Kargos!")
	log.Println("Start HTTP Server .. ")

	// Start HTTP Servers (goroutine, channel?)
	handlers.httpHandler.StartHTTPServer()

	// TODO gRPC Server (goroutine, channel?)

}
