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

func (httpHandler HTTPHandler) GetNodeUsageOverview(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	nodeUsage, err := httpHandler.k8sHandler.GetNodeUsageAvg()
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
func (httpHandler HTTPHandler) GetNodeUsage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	nodeUsage, err := httpHandler.k8sHandler.GetNodeUsage(ps.ByName("name"))
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

func (httpHandler HTTPHandler) GetPodUsage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	podUsage, err := httpHandler.k8sHandler.GetPodUsageDetail(ps.ByName("name"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&podUsage)
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

func (httpHandler HTTPHandler) GetPodInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	podInfo, err := httpHandler.k8sHandler.GetInfoOfPod(ps.ByName("name"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := json.Marshal(&podInfo)

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetContainers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	containers, err := httpHandler.k8sHandler.GetContainersOfPod(ps.ByName("name"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := json.Marshal(&containers)

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

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

func (httpHandler HTTPHandler) GetLogsOfNode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	logsOfNode, err := httpHandler.k8sHandler.GetLogsOfNode(ps.ByName("name"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := json.Marshal(&logsOfNode)

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetEvents(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	eventType := r.URL.Query().Get("event")

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

func (httpHandler HTTPHandler) GetControllerInfo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	controllerType := r.URL.Query().Get("type")

	controllerInfo, err := httpHandler.k8sHandler.GetControllerInfo(controllerType, params.ByName("namespace"), params.ByName("name"))
	result, err := json.Marshal(&controllerInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetControllerDetail(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	controllerInfo, err := httpHandler.k8sHandler.GetControllerDetail(params.ByName("namespace"), params.ByName("name"))
	result, err := json.Marshal(&controllerInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetConditions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	controllerType := r.URL.Query().Get("type")
	conditions, err := httpHandler.k8sHandler.GetConditions(controllerType, ps.ByName("namespace"), ps.ByName("name"))
	result, err := json.Marshal(&conditions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

//func (httpHandler HTTPHandler) GetTemplateContainers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//
//	controllerType := r.URL.Query().Get("type")
//	containers, err := httpHandler.k8sHandler.GetTemplateContainers(controllerType, ps.ByName("namespace"), ps.ByName("name"))
//	result, err := json.Marshal(&containers)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.Write(result)
//	w.WriteHeader(http.StatusOK)
//}

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

	nodeInfo, err := httpHandler.k8sHandler.GetEventsByController(params.ByName("name"))
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

func (httpHandler HTTPHandler) GetNumberOfControllers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	namespace := r.URL.Query().Get("namespace")
	controllerType := r.URL.Query().Get("type")
	count, err := httpHandler.k8sHandler.NumberOfControllers(namespace, controllerType)
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
func (httpHandler HTTPHandler) GetNumberOfEvents(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	eventLevel := r.URL.Query().Get("level")
	count, err := httpHandler.k8sHandler.NumberOfEvents(eventLevel)
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
