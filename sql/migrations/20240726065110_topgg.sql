-- migrate:up
CREATE TABLE IF NOT EXISTS dailies (
    player_id varchar(50) PRIMARY KEY,

    has_voted boolean NOT NULL DEFAULT false,
    last_vote_at timestamptz NOT NULL DEFAULT now(),
    streak integer NOT NULL DEFAULT 0,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);

ALTER TABLE dailies
ADD CONSTRAINT fk_player_id FOREIGN KEY (player_id) REFERENCES players (id)
DEFERRABLE INITIALLY DEFERRED;

INSERT INTO dailies (player_id)
SELECT id AS player_id
FROM players;

-- migrate:down
DROP TABLE dailies;
