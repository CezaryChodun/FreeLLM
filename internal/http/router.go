package http

import (
	nethttp "net/http"

	"github.com/cezarychodun/freellms/internal/http/handler"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	RegisterRoutes(router)

	return router
}

func RegisterRoutes(router *mux.Router) {
	proxyHandler := handler.NewProxyHandler()

	router.PathPrefix("/").Handler(nethttp.HandlerFunc(proxyHandler.Intercept))
}
