package http

import (
	"github.com/julienschmidt/httprouter"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"strconv"
)

func (httpHandler HTTPHandler) GetOverviewStatus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	overview, err := httpHandler.k8sHandler.GetOverviewStatus()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&overview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetNodeUsage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	nodeUsage, err := httpHandler.k8sHandler.GetNodeUsage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&nodeUsage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetTopNode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	topNodes, err := httpHandler.k8sHandler.GetTopNode()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&topNodes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetTopPod(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	topPods, err := httpHandler.k8sHandler.GetTopPod()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&topPods)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

//func (httpHandler HTTPHandler) GetNodeDetail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	nodeDetail, err := httpHandler.k8sHandler.GetNodeDetail(ps.ByName("name"))
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//	result, err := json.Marshal(&nodeDetail)
//
//	w.Header().Set("Content-Type", "application/json")
//	w.Write(result)
//	w.WriteHeader(http.StatusOK)
//}

//func (httpHandler HTTPHandler) GetDeploymentOverview(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//	deployOverview, err := httpHandler.k8sHandler.GetDeploymentOverview()
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	result, err := json.Marshal(&deployOverview)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.Write(result)
//	w.WriteHeader(http.StatusOK)
//}

//// GetDeploymentSpecific is a callback function for endpoint /controllers/deployment/:name
//func (httpHandler HTTPHandler) GetDeploymentSpecific(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	deploySpecific, err := httpHandler.k8sHandler.GetDeploymentSpecific(ps.ByName("namespace"), ps.ByName("name"))
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	result, err := json.Marshal(&deploySpecific)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.Write(result)
//	w.WriteHeader(http.StatusOK)
//}

//func (httpHandler HTTPHandler) GetIngressOverview(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//	ingressOverview, err := httpHandler.k8sHandler.GetIngressOverview()
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//	result, err := json.Marshal(&ingressOverview)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.Write(result)
//	w.WriteHeader(http.StatusOK)
//}
//
//// controllers/ingress/:namespace/:name
//func (httpHandler HTTPHandler) GetIngressSpecific(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	ingressSpecific, err := httpHandler.k8sHandler.GetIngressSpecific(ps.ByName("namespace"), ps.ByName("name"))
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	result, err := json.Marshal(&ingressSpecific)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.Write(result)
//	w.WriteHeader(http.StatusOK)
//}

func (httpHandler HTTPHandler) GetPodDetail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	podDetail, err := httpHandler.k8sHandler.GetRecordOfPod(ps.ByName("name"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := json.Marshal(&podDetail)

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

//func (httpHandler HTTPHandler) GetVolumeDetail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//
//	podDetail, err := httpHandler.k8sHandler.GetVolumeDetail(ps.ByName("name"))
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//	result, err := json.Marshal(&podDetail)
//
//	w.Header().Set("Content-Type", "application/json")
//	w.Write(result)
//	w.WriteHeader(http.StatusOK)
//}

func (httpHandler HTTPHandler) GetLogsOfPod(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	logsOfPod, err := httpHandler.k8sHandler.GetLogsOfPod(ps.ByName("namespace"), ps.ByName("name"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := json.Marshal(&logsOfPod)

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

//func (httpHandler HTTPHandler) GetLogsOfNode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//
//	logsOfNode, err := httpHandler.k8sHandler.GetLogsOfNode(ps.ByName("name"))
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//	result, err := json.Marshal(&logsOfNode)
//
//	w.Header().Set("Content-Type", "application/json")
//	w.Write(result)
//	w.WriteHeader(http.StatusOK)
//}

func (httpHandler HTTPHandler) GetEvents(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	eventType := r.URL.Query().Get("event")

	// Parse the query parameters
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	perPage, err := strconv.Atoi(r.URL.Query().Get("per_page"))
	if err != nil {
		perPage = 10
	}

	// Get the data from db
	events, err := httpHandler.k8sHandler.GetEvents(eventType, page, perPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetPodList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	podList, err := httpHandler.k8sHandler.GetPodsOfController(params.ByName("controller"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&podList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetNodeOverview(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	// Parse the query parameters
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	perPage, err := strconv.Atoi(r.URL.Query().Get("per_page"))
	if err != nil {
		perPage = 10
	}

	nodeOverview, err := httpHandler.k8sHandler.GetNodeOverview(page, perPage)
	result, err := json.Marshal(&nodeOverview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetNodeInfo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	nodeInfo, err := httpHandler.k8sHandler.GetNodeInfo(params.ByName("name"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&nodeInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetControllers(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	// Parse the query parameters
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	perPage, err := strconv.Atoi(r.URL.Query().Get("per_page"))
	if err != nil {
		perPage = 10
	}

	controller, err := httpHandler.k8sHandler.GetControllers(page, perPage)
	result, err := json.Marshal(&controller)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetNamespace(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	namespaceList, err := httpHandler.k8sHandler.GetNamespaceName()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&namespaceList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}
func (httpHandler HTTPHandler) GetControllersByFilter(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	// Parse the query parameters

	namespace := r.URL.Query().Get("namespace")
	controller := r.URL.Query().Get("controller")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	perPage, err := strconv.Atoi(r.URL.Query().Get("per_page"))
	if err != nil {
		perPage = 10
	}

	controllers, err := httpHandler.k8sHandler.GetControllersByFilter(namespace, controller, page, perPage)
	result, err := json.Marshal(&controllers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetControllersByType(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	// Parse the query parameters
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	perPage, err := strconv.Atoi(r.URL.Query().Get("per_page"))
	if err != nil {
		perPage = 10
	}

	controller, err := httpHandler.k8sHandler.GetControllersByType(params.ByName("controller"), page, perPage)
	result, err := json.Marshal(&controller)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetEventsByController(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	nodeInfo, err := httpHandler.k8sHandler.GetEventsByController(params.ByName("namespace"), params.ByName("name"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&nodeInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetNumberOfNodes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	count, err := httpHandler.k8sHandler.NumberOfNodes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}
func (httpHandler HTTPHandler) GetNumberOfEvents(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	count, err := httpHandler.k8sHandler.NumberOfEvents()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

//func (httpHandler HTTPHandler) GetPersistentVolume(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
//
//	// Parse the query parameters
//	page, err := strconv.Atoi(r.URL.Query().Get("page"))
//	if err != nil {
//		page = 1
//	}
//	perPage, err := strconv.Atoi(r.URL.Query().Get("per_page"))
//	if err != nil {
//		perPage = 10
//	}
//
//	pv, err := httpHandler.k8sHandler.GetPersistentVolume(page, perPage)
//	result, err := json.Marshal(&pv)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.Write(result)
//	w.WriteHeader(http.StatusOK)
//}
