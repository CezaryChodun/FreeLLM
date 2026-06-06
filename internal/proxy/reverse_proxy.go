package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cezarychodun/freellms/internal/modules/usage"
)

var httpClient = &http.Client{}

func NewReverseProxy(targetURL *url.URL, usageRepo *usage.ModelResourcesRepository, selector *ModelSelector) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var bodyBytes []byte
		if r.Body != nil && r.Body != http.NoBody {
			bodyBytes, _ = io.ReadAll(r.Body)
		}

		groupName := extractModelFromBody(bodyBytes)

		var blacklist []int
		for {
			model, err := selector.SelectExcluding(groupName, blacklist)
			if err != nil {
				http.Error(w, "no available model", http.StatusServiceUnavailable)
				return
			}

			rewritten := rewriteBody(bodyBytes, model.Name)

			if model.ID > 0 {
				if err := usageRepo.IncrementRequestCount(model.ID); err != nil {
					fmt.Printf("failed to increment request count for model %s: %v\n", model.Name, err)
				}
			}

			upstreamReq, err := http.NewRequestWithContext(r.Context(), r.Method, targetURL.String()+r.RequestURI, bytes.NewReader(rewritten))
			if err != nil {
				http.Error(w, "failed to create upstream request", http.StatusInternalServerError)
				return
			}
			upstreamReq.Header = r.Header.Clone()
			upstreamReq.ContentLength = int64(len(rewritten))
			upstreamReq.Header.Set("Content-Length", fmt.Sprintf("%d", len(rewritten)))

			resp, err := httpClient.Do(upstreamReq)
			if err != nil {
				fmt.Printf("upstream request failed for model %s: %v\n", model.Name, err)
				blacklist = append(blacklist, model.ID)
				continue
			}

			respBytes, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				blacklist = append(blacklist, model.ID)
				continue
			}

			if resp.StatusCode >= 500 {
				fmt.Printf("model %s returned status %d, retrying\n", model.Name, resp.StatusCode)
				blacklist = append(blacklist, model.ID)
				continue
			}

			if model.ID > 0 {
				stats := ExtractUsageFromResponse(respBytes)
				RegisterUsage(usageRepo, model.ID, stats)
			}

			for k, vals := range resp.Header {
				for _, v := range vals {
					w.Header().Add(k, v)
				}
			}
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(respBytes)))
			w.WriteHeader(resp.StatusCode)
			w.Write(respBytes)
			return
		}
	})
}

func rewriteBody(body []byte, modelName string) []byte {
	if len(body) == 0 {
		return body
	}
	var payload any
	if err := json.Unmarshal(body, &payload); err != nil {
		return body
	}
	payload = replaceModelValue(payload, modelName)
	out, err := json.Marshal(payload)
	if err != nil {
		return body
	}
	return out
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
