-- migrate:up
ALTER TABLE players
ADD COLUMN liens_bidded bigint NOT NULL DEFAULT 0,
ADD COLUMN username varchar(50) NOT NULL DEFAULT '';

ALTER TABLE cards
ADD COLUMN metadata jsonb NOT NULL DEFAULT '{}',
ADD COLUMN available BOOLEAN NOT NULL DEFAULT True,
ADD COLUMN owned_at timestamptz NOT NULL DEFAULT now();

CREATE TABLE IF NOT EXISTS listings (
    id varchar(50) PRIMARY KEY,
    player_id varchar(50) NOT NULL REFERENCES players (id),
    card_id varchar(50) NOT NULL REFERENCES cards (id),

    price bigint NOT NULL,
    note varchar(500) NOT NULL,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);

CREATE TABLE IF NOT EXISTS auctions (
    id varchar(50) PRIMARY KEY,
    player_id varchar(50) NOT NULL REFERENCES players (id),
    card_id varchar(50) NOT NULL REFERENCES cards (id),

    time_extensions int NOT NULL DEFAULT 0,
    ends_at timestamptz NOT NULL,

    created_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);

CREATE TABLE IF NOT EXISTS auctions_bids (
    id varchar(50) PRIMARY KEY,
    auction_id varchar(50) NOT NULL REFERENCES auctions (id),
    player_id varchar(50) NOT NULL REFERENCES players (id),

    price bigint NOT NULL,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);

-- migrate:down
DROP TABLE auctions_bids;
DROP TABLE auctions;
DROP TABLE listings;

ALTER TABLE cards
DROP COLUMN metadata,
DROP COLUMN available,
DROP COLUMN owned_at;

ALTER TABLE players
DROP COLUMN liens_bidded,
DROP COLUMN username;