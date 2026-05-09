package http

import (
	nethttp "net/http"

	"github.com/cezarychodun/freellms/internal/http/handler"
	"github.com/cezarychodun/freellms/internal/modules/usage"
	"github.com/gorilla/mux"
)

func NewRouter(modelResourcesRepository *usage.ModelResourcesRepository) *mux.Router {
	router := mux.NewRouter()

	RegisterRoutes(router, modelResourcesRepository)

	return router
}

func RegisterRoutes(router *mux.Router, modelResourcesRepository *usage.ModelResourcesRepository) {
	proxyHandler := handler.NewProxyHandler(modelResourcesRepository)

	router.PathPrefix("/").Handler(nethttp.HandlerFunc(proxyHandler.Intercept))
}
