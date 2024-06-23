-- migrate:up
CREATE TYPE auth_discord_states_type AS ENUM ('login');

CREATE TABLE IF NOT EXISTS auth_discord_states (
    state VARCHAR(50) PRIMARY KEY,

    player_id VARCHAR(50) REFERENCES players (id),
    redirect_uri TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    type auth_discord_states_type NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS auth_cookies (
    id VARCHAR(100) PRIMARY KEY,
    player_id VARCHAR(50) NOT NULL REFERENCES players (id),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

-- migrate:down
DROP TABLE auth_cookies;
DROP TABLE auth_discord_states;
DROP TYPE auth_discord_states_type;
