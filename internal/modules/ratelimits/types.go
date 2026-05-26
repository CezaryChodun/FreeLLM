package ratelimits

type RateLimit struct {
	ModelName             string `db:"model_name" json:"model_name"`
	ModelProvider         string `db:"model_provider" json:"model_provider"`
	InputTokensPerMinute  int    `db:"input_tokens_per_minute" json:"input_tokens_per_minute"`
	OutputTokensPerMinute int    `db:"output_tokens_per_minute" json:"output_tokens_per_minute"`
	RequestsPerMinute     int    `db:"requests_per_minute" json:"requests_per_minute"`
	RequestsPerDay        int    `db:"requests_per_day" json:"requests_per_day"`
}
