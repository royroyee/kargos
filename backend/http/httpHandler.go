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
	r := httprouter.New()

	log.Println("Success to start http server")

	r.GET("/overview/main", httpHandler.GetOverview)

	http.ListenAndServe(":9000", r)
}
