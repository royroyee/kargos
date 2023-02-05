package http

import (
	"github.com/julienschmidt/httprouter"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
)

func ServerValidator(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // 500
		return
	}
}

func (httpHandler HTTPHandler) GetOverview(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	home := httpHandler.k8sHandler.GetHome()

	result, err := json.Marshal(&home)
	ServerValidator(w, err) // 500

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetNodeOverview(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	nodeOverview := httpHandler.k8sHandler.GetNodeOverview()

	result, err := json.Marshal(&nodeOverview)
	ServerValidator(w, err) // 500

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetNodeDetail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	nodeDetail := httpHandler.k8sHandler.GetNodeDetail(ps.ByName("name"))

	result, err := json.Marshal(&nodeDetail)
	ServerValidator(w, err) // validation

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}
