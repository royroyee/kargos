package http

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// HTTP
type HTTPHandler struct{}

func NewHTTPHandler() *HTTPHandler {
	httpHandler := &HTTPHandler{}

	return httpHandler
}

var httpHandler HTTPHandler

func (httpHandler HTTPHandler) StartHTTPServer() {
	r := httprouter.New()

	r.GET("/overview/main", httpHandler.GetOverview)

	http.ListenAndServe(":9000", r)
}
