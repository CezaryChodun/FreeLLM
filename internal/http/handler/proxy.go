package handler

import (
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/cezarychodun/freellms/internal/modules/usage"
	"github.com/cezarychodun/freellms/internal/proxy"
)

type ProxyHandler struct {
	proxy http.Handler
}

func NewProxyHandler(modelResourcesRepository *usage.ModelResourcesRepository, selector *proxy.ModelSelector) *ProxyHandler {
	target := os.Getenv("LITELLM_URL")
	if target == "" {
		target = "http://localhost:4000"
	}
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("failed to parse proxy target URL: %v", err)
	}

	return &ProxyHandler{
		proxy: proxy.NewReverseProxy(targetURL, modelResourcesRepository, selector),
	}
}

func (h *ProxyHandler) Intercept(w http.ResponseWriter, r *http.Request) {
	h.proxy.ServeHTTP(w, r)
}
