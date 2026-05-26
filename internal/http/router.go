package http

import (
	nethttp "net/http"

	"github.com/cezarychodun/freellms/internal/http/handler"
	"github.com/cezarychodun/freellms/internal/modules/modelgroups"
	"github.com/cezarychodun/freellms/internal/modules/usage"
	"github.com/cezarychodun/freellms/internal/proxy"
	"github.com/gorilla/mux"
)

func NewRouter(modelResourcesRepository *usage.ModelResourcesRepository, selector *proxy.ModelSelector, modelGroupRepo *modelgroups.ModelGroupRepository) *mux.Router {
	router := mux.NewRouter()

	RegisterRoutes(router, modelResourcesRepository, selector, modelGroupRepo)

	return router
}

func RegisterRoutes(router *mux.Router, modelResourcesRepository *usage.ModelResourcesRepository, selector *proxy.ModelSelector, modelGroupRepo *modelgroups.ModelGroupRepository) {
	modelsHandler := handler.NewModelsHandler(modelGroupRepo)
	proxyHandler := handler.NewProxyHandler(modelResourcesRepository, selector)

	router.HandleFunc("/models", modelsHandler.ListModels).Methods("GET")
	router.HandleFunc("/v1/models", modelsHandler.ListModels).Methods("GET")

	router.PathPrefix("/forward/").Handler(nethttp.StripPrefix("/forward", nethttp.HandlerFunc(proxyHandler.Intercept)))

	router.PathPrefix("/").Handler(nethttp.HandlerFunc(proxyHandler.Intercept))
}
