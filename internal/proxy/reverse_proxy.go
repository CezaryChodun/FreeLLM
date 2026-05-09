package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewReverseProxy(targetURL *url.URL) http.Handler {
	reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)

	originalDirector := reverseProxy.Director
	reverseProxy.Director = func(r *http.Request) {
		originalDirector(r)

		RewriteModelInRequest(r, "gemma-1b")

		r.URL.Host = targetURL.Host
		r.URL.Scheme = targetURL.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = targetURL.Host
	}

	reverseProxy.ModifyResponse = func(resp *http.Response) error {
		if resp.Body == nil || resp.Body == http.NoBody {
			return nil
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()

		usage := ExtractUsageFromResponse(bodyBytes)
		RegisterUsage(usage)

		resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		resp.ContentLength = int64(len(bodyBytes))
		resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(bodyBytes)))

		return nil
	}

	return reverseProxy
}
