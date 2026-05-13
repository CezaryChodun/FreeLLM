CREATE TABLE IF NOT EXISTS model_groups (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS model_group_members (
    model_group_id INTEGER NOT NULL REFERENCES model_groups(id) ON DELETE CASCADE,
    model_id INTEGER NOT NULL REFERENCES models(id) ON DELETE CASCADE,
    PRIMARY KEY (model_group_id, model_id)
);
