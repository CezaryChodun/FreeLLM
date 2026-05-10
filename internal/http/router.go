package http

import (
	nethttp "net/http"

	"github.com/cezarychodun/freellms/internal/http/handler"
	"github.com/cezarychodun/freellms/internal/modules/usage"
	"github.com/cezarychodun/freellms/internal/proxy"
	"github.com/gorilla/mux"
)

func NewRouter(modelResourcesRepository *usage.ModelResourcesRepository, selector *proxy.ModelSelector) *mux.Router {
	router := mux.NewRouter()

	RegisterRoutes(router, modelResourcesRepository, selector)

	return router
}

func RegisterRoutes(router *mux.Router, modelResourcesRepository *usage.ModelResourcesRepository, selector *proxy.ModelSelector) {
	proxyHandler := handler.NewProxyHandler(modelResourcesRepository, selector)

	router.PathPrefix("/").Handler(nethttp.HandlerFunc(proxyHandler.Intercept))
}
