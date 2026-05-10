ALTER TABLE used_resources RENAME TO usage_tracking;
ALTER TABLE usage_tracking ADD COLUMN requests_per_minute INTEGER NOT NULL DEFAULT 0;
