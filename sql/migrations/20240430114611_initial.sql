-- migrate:up
CREATE TABLE IF NOT EXISTS players (
    id varchar(50) PRIMARY KEY,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);

CREATE TABLE IF NOT EXISTS cards (
    id varchar(50) PRIMARY KEY,
    player_id varchar(50) NOT NULL REFERENCES players (id),

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);


-- migrate:down
DROP TABLE cards;
DROP TABLE players;
