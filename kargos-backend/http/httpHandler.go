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

	// overview/main
	r.GET("/overview/main", httpHandler.GetOverview)

	// node
	r.GET("/nodes/overview", httpHandler.GetNodeOverview)
	r.GET("/node/:name", httpHandler.GetNodeDetail)

	// Controllers
	r.GET("/controllers/deployments/overview", httpHandler.GetDeploymentOverview)
	r.GET("/controllers/deployment/:namespace/:name", httpHandler.GetDeploymentSpecific)

	r.GET("/controllers/ingresses/overview", httpHandler.GetIngressOverview)
	r.GET("/controllers/ingress/:namespace/:name", httpHandler.GetIngressSpecific)

	r.GET("/controllers/jobs/overview", httpHandler.GetJobsOverview)
	r.GET("/controllers/job/:namespace/:name", httpHandler.GetJobSpecific)

	r.GET("/controllers/daemonsets/overview", httpHandler.GetDaemonSetsOverview)
	r.GET("/controllers/daemonset/:namespace/:name", httpHandler.GetDaemonSetSpecific)

	// Resources
	r.GET("/resources/namespaces/overview", httpHandler.GetNamespaceOverview)
	r.GET("/resources/namespace/:name", httpHandler.GetNamespaceDetail)

	r.GET("/resources/pods/overview", httpHandler.GetPodOverview)
	r.GET("/resources/pod/:name", httpHandler.GetPodDetail)

	r.GET("/resources/services/overview", httpHandler.GetServiceOverview)
	r.GET("/resources/service/:namespace/:name", httpHandler.GetServiceDetail)

	r.GET("/resources/persistentvolumes/overview", httpHandler.GetPersistentVolumeOverview)
	r.GET("/resources/persistentvolume/:name", httpHandler.GetPersistentVolumeDetail)

	log.Fatal(http.ListenAndServe(":9000", r))
}
