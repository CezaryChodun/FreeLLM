package usage

type Usage struct {
	ModelName    string `json:"model"`
	InputTokens  int    `json:"prompt_tokens"`
	OutputTokens int    `json:"completion_tokens"`
	Timestamp    int    `json:"timestamp"`
}

type OpenAIResponse struct {
	Model     string        `json:"model"`
	Usage     UsageResponse `json:"usage"`
	Timestamp int           `json:"timestamp"`
}

type UsageResponse struct {
	InputTokens  int `json:"prompt_tokens"`
	OutputTokens int `json:"completion_tokens"`
}
