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

	"github.com/cezarychodun/freellms/app/usage"
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

		printUsageFromResponse(bodyBytes)

		resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		resp.ContentLength = int64(len(bodyBytes))
		resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(bodyBytes)))

		return nil
	}

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

	// Update headers for SSL/forwarding
	r.URL.Host = targetURL.Host
	r.URL.Scheme = targetURL.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = targetURL.Host

	proxy.ServeHTTP(w, r)
}

func printUsageFromResponse(bodyBytes []byte) {
	var payload usage.GeminiResponse
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		fmt.Println("Error:", err)
		return
	}

	usageBytes, err := json.Marshal(payload.Usage)
	if err != nil {
		fmt.Println("usage:", payload)
		return
	}

	fmt.Println("usage:", string(usageBytes))
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
