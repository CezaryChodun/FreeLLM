CREATE TABLE IF NOT EXISTS remaining_resources (
	model TEXT PRIMARY KEY,
	input_tokens_per_minute INTEGER NOT NULL,
	output_tokens_per_minute INTEGER NOT NULL,
	requests_per_day INTEGER NOT NULL,
	last_used TIMESTAMP NOT NULL
);
