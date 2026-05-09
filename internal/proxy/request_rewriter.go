package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func RewriteModelInRequest(r *http.Request, modelName string) {
	if r.Body == nil || r.Body == http.NoBody {
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		return
	}

	if len(bodyBytes) == 0 || !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		return
	}

	var payload any
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		return
	}

	payload = replaceModelValue(payload, modelName)

	rewrittenBody, err := json.Marshal(payload)
	if err != nil {
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		return
	}

	r.Body = io.NopCloser(bytes.NewReader(rewrittenBody))
	r.ContentLength = int64(len(rewrittenBody))
	r.Header.Set("Content-Length", fmt.Sprintf("%d", len(rewrittenBody)))
}

func replaceModelValue(v any, modelName string) any {
	switch value := v.(type) {
	case map[string]any:
		for k, child := range value {
			if k == "model" {
				value[k] = modelName
				continue
			}

			value[k] = replaceModelValue(child, modelName)
		}

		return value
	case []any:
		for i, item := range value {
			value[i] = replaceModelValue(item, modelName)
		}

		return value
	default:
		return v
	}
}
