package http

import (
	"fmt"
	"github.com/boanlab/kargos/k8s"
	"github.com/julienschmidt/httprouter"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
)

func ServerValidator(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var kh k8s.K8sHandler

func (httpHandler HTTPHandler) GetOverview(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	fmt.Println("check getOverview 1")
	home := httpHandler.k8sHandler.GetHome()

	result, err := json.Marshal(&home)
	ServerValidator(w, err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)

}
