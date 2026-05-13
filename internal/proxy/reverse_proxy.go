package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/cezarychodun/freellms/internal/modules/usage"
)

func NewReverseProxy(targetURL *url.URL, usageRepo *usage.ModelResourcesRepository, selector *ModelSelector) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		model, err := selector.Next()
		if err != nil {
			fmt.Printf("model selection failed: %v\n", err)
		} else {
			RewriteModelInRequest(r, model.Name)
		}

		rp := httputil.NewSingleHostReverseProxy(targetURL)
		rp.Director = func(req *http.Request) {
			req.URL.Host = targetURL.Host
			req.URL.Scheme = targetURL.Scheme
			req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
			req.Host = targetURL.Host
		}
		rp.ModifyResponse = func(resp *http.Response) error {
			if resp.Body == nil || resp.Body == http.NoBody {
				return nil
			}
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			resp.Body.Close()

			if model.ID > 0 {
				stats := ExtractUsageFromResponse(bodyBytes)
				RegisterUsage(usageRepo, model.ID, stats)
			}

			resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			resp.ContentLength = int64(len(bodyBytes))
			resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(bodyBytes)))
			return nil
		}

		rp.ServeHTTP(w, r)
	})
}
