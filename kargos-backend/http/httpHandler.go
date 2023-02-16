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
	r.GET("/overview/nodes/top", httpHandler.GetTopNode) // 현재 정렬이 되지 않는 문제 있음
	r.GET("/overview/pods/top", httpHandler.GetTopPod)   // 위처럼 정렬문제와 애초에 cpu , ram 사용량 값 자체에 대한 에러 있음..

	// Events
	r.GET("/events", httpHandler.GetEvents) // Example : localhost:9000/events/?event=warning&page=1&per_page=10
	r.GET("/events/count", httpHandler.GetNumberOfEvents)

	// Nodes
	r.GET("/nodes", httpHandler.GetNodeOverview)
	r.GET("/node/info/:name", httpHandler.GetNodeInfo)
	r.GET("/nodes/count", httpHandler.GetNumberOfNodes)
	//r.GET("/node/logs/:name", httpHandler.GetLogsOfNode) (TODO REST client 필요)

	// Monitor
	r.GET("/monitor/namespaces", httpHandler.GetNamespace)            // monitor 페이지에서 필터링 할 namespace 목록 리턴
	r.GET("/monitor/controllers", httpHandler.GetControllersByFilter) // Filtering by Namespace
	// Example "/monitor/controllers" : all Controllers   "/monitor/controllers?namespace=kargos : Controllers of Kargos
	// localhost:9000/monitor/controllers?namespace=kargos&page=1&per_page=10 (&pagination)
	// localhost:9000/monitor/controllers?namespace=kargos&controller=daemonset namespace / type / page 까지!

	//r.GET("/monitor/namespaces/controller", httpHandler.GetControllers) // namespace 로 필터링하고, 거기서 deployment 등의 필터링 (namespace 입력안한 것도 고려해야 함 위의 필터처럼)

	r.GET("/monitor/pods/:controller", httpHandler.GetPodList) // pod List of controller
	r.GET("/pod/detail/:name", httpHandler.GetPodDetail)       // Information of Pod (detail page)
	//	r.GET("pod/usage/:name", httpHandler.GetPodUsage)          // Usage(cpu,ram) of Pod (detail page)

	r.GET("/pod/logs/:namespace/:name", httpHandler.GetLogsOfPod) // (현재 24시간 이전 log만 반환하는데, 반환값이 이쁘지 않을 때도 있고, 너무 많을 때도 있는 듯)

	r.GET("/monitor/controller/events/:namespace/:name", httpHandler.GetEventsByController) // 컨트롤러의 events 반환 (10개만)
	//r.GET("/volume/detail/:name", httpHandler.GetVolumeDetail)
	// Resources/Persistent Volumesx
	//r.GET("/resources/persistentvolumes", httpHandler.GetPersistentVolume)

	//// node
	//r.GET("/nodes/overview", httpHandler.GetNodeOverview)
	//r.GET("/node/:name", httpHandler.GetNodeDetail)

	// Controllers
	//r.GET("/controllers/deployments/overview", httpHandler.GetDeploymentOverview)
	//r.GET("/controllers/deployment/:namespace/:name", httpHandler.GetDeploymentSpecific)
	//
	//r.GET("/controllers/ingresses/overview", httpHandler.GetIngressOverview)
	//r.GET("/controllers/ingress/:namespace/:name", httpHandler.GetIngressSpecific)
	//
	//r.GET("/controllers/jobs/overview", httpHandler.GetJobsOverview)
	//r.GET("/controllers/job/:namespace/:name", httpHandler.GetJobSpecific)
	//
	//r.GET("/controllers/daemonsets/overview", httpHandler.GetDaemonSetsOverview)
	//r.GET("/controllers/daemonset/:namespace/:name", httpHandler.GetDaemonSetSpecific)
	//
	//// Resources
	//r.GET("/resources/namespaces/overview", httpHandler.GetNamespaceOverview)
	//r.GET("/resources/namespace/:name", httpHandler.GetNamespaceDetail)
	//
	//r.GET("/resources/pods/overview", httpHandler.GetPodOverview)
	//r.GET("/resources/pod/:name", httpHandler.GetPodDetail)
	//
	//r.GET("/resources/services/overview", httpHandler.GetServiceOverview)
	//r.GET("/resources/service/:namespace/:name", httpHandler.GetServiceDetail)
	//
	//r.GET("/resources/persistentvolumes/overview", httpHandler.GetPersistentVolumeOverview)
	//r.GET("/resources/persistentvolume/:name", httpHandler.GetPersistentVolumeDetail)

	log.Fatal(http.ListenAndServe(":9000", r))

}
