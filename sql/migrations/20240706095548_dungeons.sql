-- migrate:up
CREATE TABLE IF NOT EXISTS dungeons (
    id varchar(50) PRIMARY KEY,
    player_id varchar(50) NOT NULL REFERENCES players (id),

    seed bigint NOT NULL,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);

CREATE TABLE IF NOT EXISTS player_limits (
    player_id varchar(50) PRIMARY KEY,

    -- dungeons
    dungeons_left int NOT NULL DEFAULT 3,
    dungeons_reset_at timestamptz,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);

ALTER TABLE player_limits
ADD CONSTRAINT fk_player_id FOREIGN KEY (player_id) REFERENCES players (id)
DEFERRABLE INITIALLY DEFERRED;

INSERT INTO player_limits (player_id)
SELECT id AS player_id
FROM players;

-- migrate:down
DROP TABLE player_limits;
DROP TABLE dungeons;
