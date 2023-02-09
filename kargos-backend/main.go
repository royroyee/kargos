package main

import (
	"fmt"
	grpcH "github.com/boanlab/kargos/gRPC"
	"github.com/boanlab/kargos/http"
	"github.com/boanlab/kargos/k8s"
	"log"
	"sync"
)

var handlers Handlers

type Handlers struct {
	k8sHandler  *k8s.K8sHandler
	httpHandler *http.HTTPHandler
	grpcHandler *grpcH.Handler
}

func initHandlers() {
	handlers.k8sHandler = k8s.NewK8sHandler()
	handlers.httpHandler = http.NewHTTPHandler(handlers.k8sHandler)
	handlers.grpcHandler = grpcH.NewGRPCHandler(handlers.k8sHandler)
}

func init() {
	log.SetPrefix("Kargos: ")
}

// printLogo prints out the logo of our Kargos in ASCII art.
func printLogo() {
	fmt.Println("  _  __                         ")
	fmt.Println(" | |/ /__ _ _ __ __ _  ___  ___ ")
	fmt.Println(" | ' // _` | '__/ _` |/ _ \\/ __|")
	fmt.Println(" | . \\ (_| | | | (_| | (_) \\__ \\")
	fmt.Println(" |_|\\_\\__,_|_|  \\__, |\\___/|___/")
	fmt.Println("                |___/           ")
	fmt.Printf("A Kubernetes Management and Monitoring System - https://github.com/boanlab/kargos\n\n")
}

func main() {
	printLogo()
	log.Println("Welcome to Kargos!")

	// Handlers
	initHandlers()

	var wg sync.WaitGroup
	wg.Add(3)

	// Start DB Session
	go handlers.k8sHandler.DBSession()

	// Start HTTP Servers (goroutine, channel?)
	go handlers.httpHandler.StartHTTPServer()

	go handlers.grpcHandler.StartGRPCServer()
	// TODO work on channel for gRPC <-> REST API data

	wg.Wait()
	log.Println("Kargos Backend finished. Bye.")
}
