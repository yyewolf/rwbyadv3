-- migrate:up
CREATE TABLE IF NOT EXISTS github_stars (
    player_id VARCHAR(50) PRIMARY KEY,

    github_user_id VARCHAR(50) UNIQUE,
    has_starred BOOLEAN NOT NULL DEFAULT false,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

INSERT INTO github_stars (player_id)
SELECT id AS player_id
FROM players;

ALTER TABLE github_stars
ADD CONSTRAINT fk_player_id FOREIGN KEY (player_id) REFERENCES players (id)
DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE players
ADD CONSTRAINT fk_github_star FOREIGN KEY (id) REFERENCES github_stars (player_id)
DEFERRABLE INITIALLY DEFERRED;

CREATE TYPE auth_github_states_type AS ENUM ('check_star');

CREATE TABLE IF NOT EXISTS auth_github_states (
    state VARCHAR(50) PRIMARY KEY,

    player_id VARCHAR(50) NOT NULL REFERENCES players (id),
    expires_at TIMESTAMPTZ NOT NULL,
    type AUTH_GITHUB_STATES_TYPE NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

-- migrate:down
DROP TABLE auth_github_states;
ALTER TABLE github_stars DROP CONSTRAINT fk_player_id;
ALTER TABLE players DROP CONSTRAINT fk_github_star;
DROP TYPE AUTH_GITHUB_STATES_TYPE;
DROP TABLE github_stars;
