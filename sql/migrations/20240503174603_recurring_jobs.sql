-- migrate:up
ALTER TABLE jobs
ADD COLUMN last_run_id BIGINT NOT NULL DEFAULT 0,
ADD COLUMN recurring BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN delta_time BIGINT NOT NULL DEFAULT 0,
ADD COLUMN errored BOOLEAN NOT NULL DEFAULT false;

-- migrate:down
ALTER TABLE jobs
DROP COLUMN last_run_id,
DROP COLUMN recurring,
DROP COLUMN delta_time,
DROP COLUMN errored;
