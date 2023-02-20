package http

import (
	"github.com/boanlab/kargos/k8s"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

// HTTP
type HTTPHandler struct {
	k8sHandler k8s.K8sHandler
}

func NewHTTPHandler(k8sHandler *k8s.K8sHandler) *HTTPHandler {
	httpHandler := &HTTPHandler{
		k8sHandler: *k8sHandler,
	}
	return httpHandler
}

func (httpHandler HTTPHandler) StartHTTPServer() {
	log.Println("Start HTTP Server .. ")

	r := httprouter.New()

	log.Println("Success to Start HTTP Server")

	// Overview
	r.GET("/overview/status", httpHandler.GetOverviewStatus)
	r.GET("/overview/nodes/usage", httpHandler.GetNodeUsage)
	r.GET("/overview/nodes/top", httpHandler.GetTopNode)
	r.GET("/overview/pods/top", httpHandler.GetTopPod)

	// Events
	r.GET("/events", httpHandler.GetEvents) // Example : /events/?event=warning&page=1&per_page=10
	r.GET("/events/count", httpHandler.GetNumberOfEvents)

	// Nodes
	r.GET("/nodes", httpHandler.GetNodeOverview)
	r.GET("/node/usage/:name", httpHandler.GetNodeUsage)
	r.GET("/node/info/:name", httpHandler.GetNodeInfo)
	r.GET("/nodes/count", httpHandler.GetNumberOfNodes)
	r.GET("/node/logs/:name", httpHandler.GetLogsOfNode) // TODO (ERROR)

	// Workload
	r.GET("/workload/namespaces", httpHandler.GetNamespace)
	r.GET("/workload", httpHandler.GetControllersByFilter) // Filtering by Namespace, Type

	r.GET("/workload/count", httpHandler.GetNumberOfControllers)
	r.GET("/workload/info/:namespace/:name", httpHandler.GetControllerInfo)

	r.GET("/pod/info/:name", httpHandler.GetPodInfo) // Information of Pod (detail page)
	r.GET("/pod/usage/:name", httpHandler.GetPodUsage)

	r.GET("/pod/logs/:namespace/:name", httpHandler.GetLogsOfPod) // TODO (ERROR)

	//r.GET("/workload/controller/events/:namespace/:name", httpHandler.GetEventsByController) // Only 10

	log.Fatal(http.ListenAndServe(":9000", r))

}
