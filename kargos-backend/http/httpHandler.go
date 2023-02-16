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
	//	r.GET("/overview/pods/top", httpHandler.GetTopPod)

	// Events
	r.GET("/events", httpHandler.GetEvents) // Example : localhost:9000/events/?event=warning&page=1&per_page=10

	// Nodes
	r.GET("/nodes", httpHandler.GetNodeOverview)
	r.GET("/node/info/:name", httpHandler.GetNodeInfo)

	// Monitor

	r.GET("/monitor/namespaces", httpHandler.GetNamespace)               // monitor 페이지에서 필터링 할 namespace 목록 리턴
	r.GET("/monitor/controllers", httpHandler.GetControllersByNamespace) // Filtering by Namespace
	// Example "/monitor/controllers" : all Controllers   "/monitor/controllers?namespace=kargos : Controllers of Kargos        localhost:9000/monitor/controllers?namespace=kargos&page=1&per_page=10 (&pagination)

	//r.GET("/monitor/namespaces/controller", httpHandler.GetControllers) // namespace 로 필터링하고, 거기서 deployment 등의 필터링 (namespace 입력안한 것도 고려해야 함 위의 필터처럼)

	r.GET("/monitor/pods/:controller", httpHandler.GetPodList) // pod List of controller
	r.GET("/pod/detail/:name", httpHandler.GetPodDetail)       // detail of pod
	// r.GET("/monitor/:controller/events, ..(paging 없이 이벤트 10개만 리미트 걸기)
	r.GET("/pod/logs/:namespace/:name", httpHandler.GetLogsOfPod) // (log인데 일단 보류. k8s쪽의 GetLogsOfPod 함수도 아직 반환값없고 print 만 해놈)

	// Resources/Persistent Volumes
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
