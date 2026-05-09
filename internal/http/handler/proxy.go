package handler

import (
	"log"
	"net/http"
	"net/url"

	"github.com/cezarychodun/freellms/internal/modules/usage"
	"github.com/cezarychodun/freellms/internal/proxy"
)

type ProxyHandler struct {
	proxy http.Handler
}

func NewProxyHandler(modelResourcesRepository *usage.ModelResourcesRepository) *ProxyHandler {
	targetURL, err := url.Parse("http://0.0.0.0:4000")
	if err != nil {
		log.Fatalf("failed to parse proxy target URL: %v", err)
	}

	return &ProxyHandler{
		proxy: proxy.NewReverseProxy(targetURL, modelResourcesRepository),
	}
}

func (h *ProxyHandler) Intercept(w http.ResponseWriter, r *http.Request) {
	h.proxy.ServeHTTP(w, r)
}
