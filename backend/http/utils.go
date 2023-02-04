package http

import (
	"github.com/boanlab/kargos/backend/k8s"
	"github.com/julienschmidt/httprouter"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
)

var kh k8s.K8sHandler

func ServerValidator(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler HTTPHandler) GetOverview(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	home := kh.GetHome()

	result, err := json.Marshal(&home)
	ServerValidator(w, err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)

}
