package proxy

type OpenAIResponse struct {
	Model     string       `json:"model"`
	Usage     UsagePayload `json:"usage"`
	Timestamp int          `json:"timestamp"`
}

type UsagePayload struct {
	InputTokens  int `json:"prompt_tokens"`
	OutputTokens int `json:"completion_tokens"`
}
