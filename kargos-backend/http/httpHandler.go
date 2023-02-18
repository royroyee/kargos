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
	r.GET("/node/info/:name", httpHandler.GetNodeInfo)
	r.GET("/nodes/count", httpHandler.GetNumberOfNodes)
	//r.GET("/node/logs/:name", httpHandler.GetLogsOfNode) (TODO REST client 필요)

	// Monitor
	r.GET("/workload/namespaces", httpHandler.GetNamespace)
	r.GET("/workload/controllers", httpHandler.GetControllersByFilter) // Filtering by Namespace
	// Example "/monitor/controllers" : all Controllers   "/monitor/controllers?namespace=kargos : Controllers of Kargos
	// localhost:9000/monitor/controllers?namespace=kargos&page=1&per_page=10 (&pagination)
	// localhost:9000/monitor/controllers?namespace=kargos&controller=daemonset

	//r.GET("/monitor/namespaces/controller", httpHandler.GetControllers)

	r.GET("/workload/pods/:controller", httpHandler.GetPodList) // pod List of controller
	r.GET("/pod/detail/:name", httpHandler.GetPodDetail)        // Information of Pod (detail page)
	//	r.GET("pod/usage/:name", httpHandler.GetPodUsage)          // Usage(cpu,ram) of Pod (detail page)

	r.GET("/pod/logs/:namespace/:name", httpHandler.GetLogsOfPod) // TODO

	r.GET("/workload/controller/events/:namespace/:name", httpHandler.GetEventsByController) // Only 10

	log.Fatal(http.ListenAndServe(":9000", r))

}
