-- Create models table
CREATE TABLE IF NOT EXISTS models (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    provider TEXT NOT NULL,
    instance INTEGER NOT NULL DEFAULT 1
);

-- Drop old tables and recreate with model_id FK
DROP TABLE IF EXISTS rate_limits;
CREATE TABLE rate_limits (
    model_name TEXT NOT NULL,
    model_provider TEXT NOT NULL,
    input_tokens_per_minute INTEGER NOT NULL,
    output_tokens_per_minute INTEGER NOT NULL,
    requests_per_minute INTEGER NOT NULL,
    requests_per_day INTEGER NOT NULL,
    PRIMARY KEY (model_name, model_provider)
);

DROP TABLE IF EXISTS usage_tracking;
CREATE TABLE usage_tracking (
    model_id INTEGER PRIMARY KEY REFERENCES models(id) ON DELETE CASCADE,
    input_tokens_per_minute INTEGER NOT NULL,
    output_tokens_per_minute INTEGER NOT NULL,
    requests_per_minute INTEGER NOT NULL,
    requests_per_day INTEGER NOT NULL,
    last_used TIMESTAMP NOT NULL
);
