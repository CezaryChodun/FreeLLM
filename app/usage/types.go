package usage

type GeminiResponse struct {
	Model string        `json:"model"`
	Usage UsageResponse `json:"usage"`
}

type UsageResponse struct {
	InputTokens  int `json:"prompt_tokens"`
	OutputTokens int `json:"completion_tokens"`
}
