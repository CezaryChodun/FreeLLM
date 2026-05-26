package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/cezarychodun/freellms/internal/modules/usage"
)

func NewReverseProxy(targetURL *url.URL, usageRepo *usage.ModelResourcesRepository, selector *ModelSelector) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read body to extract model name
		var bodyBytes []byte
		if r.Body != nil && r.Body != http.NoBody {
			bodyBytes, _ = io.ReadAll(r.Body)
		}

		groupName := extractModelFromBody(bodyBytes)

		// Restore body for downstream use
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		model, err := selector.Select(groupName)
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
			respBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			resp.Body.Close()

			if model.ID > 0 {
				stats := ExtractUsageFromResponse(respBytes)
				RegisterUsage(usageRepo, model.ID, stats)
			}

			resp.Body = io.NopCloser(bytes.NewReader(respBytes))
			resp.ContentLength = int64(len(respBytes))
			resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(respBytes)))
			return nil
		}

		rp.ServeHTTP(w, r)
	})
}

func extractModelFromBody(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	var payload struct {
		Model string `json:"model"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return ""
	}
	return payload.Model
}
