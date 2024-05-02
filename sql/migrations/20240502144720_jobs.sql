-- migrate:up
CREATE TABLE IF NOT EXISTS jobs (
    id varchar(100),
    jobkey varchar(100),

    retries int NOT NULL DEFAULT 0,
    run_at timestamptz NOT NULL,
    params jsonb NOT NULL,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz,

    PRIMARY KEY (id, jobkey)
);

-- migrate:down
DROP TABLE jobs;
