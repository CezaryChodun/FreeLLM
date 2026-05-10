CREATE TABLE IF NOT EXISTS rate_limits (
    model TEXT PRIMARY KEY,
    input_tokens_per_minute INTEGER NOT NULL,
    output_tokens_per_minute INTEGER NOT NULL,
    requests_per_day INTEGER NOT NULL
);
