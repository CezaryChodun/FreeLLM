package usage

type Usage struct {
	ModelName    string `json:"model"`
	InputTokens  int    `json:"prompt_tokens"`
	OutputTokens int    `json:"completion_tokens"`
	Timestamp    int    `json:"timestamp"`
}
