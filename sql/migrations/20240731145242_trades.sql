-- migrate:up
CREATE TABLE IF NOT EXISTS trades (
    id uuid PRIMARY KEY,
    sender_id varchar(50) NOT NULL REFERENCES players (id),
    receiver_id varchar(50) NOT NULL REFERENCES players (id),

    metadata json NOT NULL,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);

-- migrate:down
DROP TABLE IF EXISTS trades;
