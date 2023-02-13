package http

import (
	"github.com/julienschmidt/httprouter"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"strconv"
)

func (httpHandler HTTPHandler) GetOverview(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	overview, err := httpHandler.k8sHandler.GetOverview()
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

func (httpHandler HTTPHandler) GetDeploymentOverview(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	deployOverview, err := httpHandler.k8sHandler.GetDeploymentOverview()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&deployOverview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

// GetDeploymentSpecific is a callback function for endpoint /controllers/deployment/:name
func (httpHandler HTTPHandler) GetDeploymentSpecific(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	deploySpecific, err := httpHandler.k8sHandler.GetDeploymentSpecific(ps.ByName("namespace"), ps.ByName("name"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&deploySpecific)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetIngressOverview(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ingressOverview, err := httpHandler.k8sHandler.GetIngressOverview()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := json.Marshal(&ingressOverview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

// controllers/ingress/:namespace/:name
func (httpHandler HTTPHandler) GetIngressSpecific(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ingressSpecific, err := httpHandler.k8sHandler.GetIngressSpecific(ps.ByName("namespace"), ps.ByName("name"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&ingressSpecific)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

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

func (httpHandler HTTPHandler) GetNamespaceOverview(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	namespaceOverview, err := httpHandler.k8sHandler.GetNamespaceOverview()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&namespaceOverview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetNamespaceDetail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	namespaceDetail, err := httpHandler.k8sHandler.GetNamespaceDetail(ps.ByName("name"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&namespaceDetail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

//func (httpHandler HTTPHandler) GetJobsOverview(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//
//	jobOverview, err := httpHandler.k8sHandler.GetJobsOverview()
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	result, err := json.Marshal(&jobOverview)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.Write(result)
//	w.WriteHeader(http.StatusOK)
//}

// controllers/job/:namespace/:name
func (httpHandler HTTPHandler) GetJobSpecific(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	jobSpecific, err := httpHandler.k8sHandler.GetJobSpecific(ps.ByName("namespace"), ps.ByName("name"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&jobSpecific)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetDaemonSetsOverview(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	daemonSetOverview, err := httpHandler.k8sHandler.GetDaemonSetsOverview()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&daemonSetOverview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

// controllers/job/:namespace/:name
func (httpHandler HTTPHandler) GetDaemonSetSpecific(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	daemonSetSpecific, err := httpHandler.k8sHandler.GetDaemonSetSpecific(ps.ByName("namespace"), ps.ByName("name"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&daemonSetSpecific)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetServiceOverview(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	serviceOverview, err := httpHandler.k8sHandler.GetServiceOverview()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&serviceOverview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}
func (httpHandler HTTPHandler) GetServiceDetail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	serviceDetail, err := httpHandler.k8sHandler.GetServiceDetail(ps.ByName("namespace"), ps.ByName("name"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&serviceDetail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

// "events/alerts?page=<page>&per_page=<per_page>"
// Get Alerts (type = "Warning", "Critical")
func (httpHandler HTTPHandler) GetAlerts(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

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
	alerts, err := httpHandler.k8sHandler.GetAlerts(page, perPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&alerts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

// "events/info?page=<page>&per_page=<per_page>"
// Get Info (events other than the warning, critical type
func (httpHandler HTTPHandler) GetInfo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

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
	info, err := httpHandler.k8sHandler.GetInfo(page, perPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetPodOverview(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	// Parse the query parameters
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	perPage, err := strconv.Atoi(r.URL.Query().Get("per_page"))
	if err != nil {
		perPage = 10
	}

	podOverview, err := httpHandler.k8sHandler.GetPodOverview(page, perPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&podOverview)
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

func (httpHandler HTTPHandler) GetPersistentVolume(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	// Parse the query parameters
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	perPage, err := strconv.Atoi(r.URL.Query().Get("per_page"))
	if err != nil {
		perPage = 10
	}

	pv, err := httpHandler.k8sHandler.GetPersistentVolume(page, perPage)
	result, err := json.Marshal(&pv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}
