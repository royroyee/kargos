package main

import (
	cm "github.com/boanlab/kargos/common"
	"github.com/boanlab/kargos/http"
	"github.com/boanlab/kargos/k8s"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"log"
)

var handlers Handlers

type Handlers struct {
	k8sHandler  *k8s.K8sHandler
	httpHandler *http.HTTPHandler

	// TODO gRPCHandler
}

var k8sHandler k8s.K8sHandler
var httpHandler http.HTTPHandler

var k8sClient *kubernetes.Clientset
var metricK8sClient *versioned.Clientset

func initHandlers() {

	handlers.k8sHandler = k8s.NewK8sHandler(k8sClient, metricK8sClient)
	handlers.httpHandler = http.NewHTTPHandler(handlers.k8sHandler)

	// TODO gRPC Server
}

func initClients() {
	// In Cluster
	//clientSet = cm.InitK8sClient()
	//metriClientSet = cm.MetricClientSetOutofCluster()

	// Out of Cluster
	k8sClient = cm.ClientSetOutofCluster()
	metricK8sClient = cm.MetricClientSetOutofCluster()
}

func init() {
	log.SetPrefix("Kargos: ")
}

func main() {
	initClients()
	initHandlers()

	log.Println("Welcome Kargos!")
	log.Println("Start HTTP Server .. ")
	// Start HTTP Servers
	handlers.httpHandler.StartHTTPServer()

	// TODO gRPC Server

}
