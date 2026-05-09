package proxy

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cezarychodun/freellms/internal/modules/usage"
)

func ExtractUsageFromResponse(bodyBytes []byte) usage.Usage {
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

		var payload OpenAIResponse
		if err := json.Unmarshal([]byte(body), &payload); err != nil {
			fmt.Printf("failed to parse usage response chunk: %+v\n", err)
			continue
		}

		if payload.Model != "" {
			stats.ModelName = payload.Model
		}

		stats.InputTokens += payload.Usage.InputTokens
		stats.OutputTokens += payload.Usage.OutputTokens
		stats.Timestamp = max(payload.Timestamp, stats.Timestamp)
	}

	return stats
}

func RegisterUsage(usage usage.Usage) {
	fmt.Println("Registering usage...")
	fmt.Println("Usage:", usage)
}
