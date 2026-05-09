package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/cezarychodun/freellms/internal/modules/usage"
)

func Intercept(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello!")

	targetURL, _ := url.Parse("http://0.0.0.0:4000")
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	proxy.ModifyResponse = func(resp *http.Response) error {
		fmt.Println("Response received, processing usage information...")
		if resp.Body == nil || resp.Body == http.NoBody {
			return nil
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()

		usage := getUsageFromResponse(bodyBytes)
		registerUsage(usage)

		resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		resp.ContentLength = int64(len(bodyBytes))
		resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(bodyBytes)))

		return nil
	}

	replaceModelInBody(r)

	// Update headers for SSL/forwarding
	r.URL.Host = targetURL.Host
	r.URL.Scheme = targetURL.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = targetURL.Host

	proxy.ServeHTTP(w, r)
}

func replaceModelInBody(r *http.Request) {
	// Rewrite request body before forwarding.
	if r.Body != nil && r.Body != http.NoBody {
		bodyBytes, err := io.ReadAll(r.Body)
		if err == nil && len(bodyBytes) > 0 && strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			var payload any
			if json.Unmarshal(bodyBytes, &payload) == nil {
				payload = replaceModelValue(payload)
				if rewrittenBody, err := json.Marshal(payload); err == nil {
					r.Body = io.NopCloser(bytes.NewReader(rewrittenBody))
					r.ContentLength = int64(len(rewrittenBody))
					r.Header.Set("Content-Length", fmt.Sprintf("%d", len(rewrittenBody)))
				}
			} else {
				r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			}
		} else {
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}
	}
}

func getUsageFromResponse(bodyBytes []byte) usage.Usage {
	fmt.Println("Calculating usage...")
	fmt.Println("Body:", string(bodyBytes))

	lines := strings.Split(string(bodyBytes), "\n")
	stats := usage.Usage{}
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" || (!strings.HasPrefix(line, "data:") && !strings.HasPrefix(line, "{")) {
			continue
		}

		body := line
		if strings.HasPrefix(line, "data:") {
			body = strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		}

		if body == "" || body == "[DONE]" {
			continue
		}

		var payload usage.OpenAIResponse
		if err := json.Unmarshal([]byte(body), &payload); err != nil {
			fmt.Printf("Error: %+v\n", err)
			continue
		}

		fmt.Printf("Registred usage: in: %d, out: %d \n", payload.Usage.InputTokens, payload.Usage.OutputTokens)
		if payload.Model != "" {
			stats.ModelName = payload.Model
		}
		stats.InputTokens += payload.Usage.InputTokens
		stats.OutputTokens += payload.Usage.OutputTokens
		stats.Timestamp = max(payload.Timestamp, stats.Timestamp)
	}

	return stats
}

func registerUsage(usage usage.Usage) {
	fmt.Println("Registering usage...")
	fmt.Println("Usage:", usage)
}

func replaceModelValue(v any) any {
	switch value := v.(type) {
	case map[string]any:
		for k, child := range value {
			if k == "model" {
				value[k] = "gemma-1b" //"gemini-flash"
				continue
			}
			value[k] = replaceModelValue(child)
		}
		return value
	case []any:
		for i, item := range value {
			value[i] = replaceModelValue(item)
		}
		return value
	default:
		return v
	}
}
