package ratelimits

type RateLimit struct {
	Name                  string `db:"name" json:"name"`
	InputTokensPerMinute  int    `db:"input_tokens_per_minute" json:"input_tokens_per_minute"`
	OutputTokensPerMinute int    `db:"output_tokens_per_minute" json:"output_tokens_per_minute"`
	RequestsPerDay        int    `db:"requests_per_day" json:"requests_per_day"`
}
