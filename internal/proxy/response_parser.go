package proxy

import (
	"encoding/json"
	"errors"
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

func RegisterUsage(repository *usage.ModelResourcesRepository, stats usage.Usage) {
	if stats.ModelName == "" {
		fmt.Println("Skipping usage registration: response does not contain model name")
		return
	}

	if stats.InputTokens == 0 && stats.OutputTokens == 0 {
		fmt.Println("Skipping usage registration: response does not contain token usage")
		return
	}

	if err := repository.AddTokenUsage(stats.ModelName, stats.InputTokens, stats.OutputTokens); err != nil {
		if errors.Is(err, usage.ErrModelResourcesNotFound) {
			fmt.Printf("Skipping usage registration: unknown model %q\n", stats.ModelName)
			return
		}

		fmt.Printf("Failed to register usage for model %q: %+v\n", stats.ModelName, err)
		return
	}

	fmt.Printf(
		"Registered usage for model %q: input_tokens=%d output_tokens=%d\n",
		stats.ModelName,
		stats.InputTokens,
		stats.OutputTokens,
	)
}
