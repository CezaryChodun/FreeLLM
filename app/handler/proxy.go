package handler

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func Intercept(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello!")

	targetURL, _ := url.Parse("http://0.0.0.0:4000")
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	// Update headers for SSL/forwarding
	r.URL.Host = targetURL.Host
	r.URL.Scheme = targetURL.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = targetURL.Host

	proxy.ServeHTTP(w, r)
}
