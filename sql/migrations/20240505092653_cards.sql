-- migrate:up
ALTER TABLE cards
ADD COLUMN xp INT NOT NULL DEFAULT 0,
ADD COLUMN next_level_xp INT NOT NULL,
ADD COLUMN card_type varchar(50) NOT NULL,
ADD COLUMN individual_value FLOAT NOT NULL,
ADD COLUMN rarity INT NOT NULL,
ADD COLUMN level INT NOT NULL,
ADD COLUMN buffs INT NOT NULL;

CREATE TABLE IF NOT EXISTS cards_stats (
    card_id varchar(50) PRIMARY KEY REFERENCES cards (id),

    health INT NOT NULL,
    armor INT NOT NULL,
    damage INT NOT NULL,
    healing INT NOT NULL,
    speed INT NOT NULL,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);

ALTER TABLE cards
ADD CONSTRAINT fk_cards_stats FOREIGN KEY (id) REFERENCES cards_stats (card_id)
DEFERRABLE INITIALLY DEFERRED;

-- migrate:down
ALTER TABLE cards DROP CONSTRAINT fk_cards_stats;

DROP TABLE cards_stats;

ALTER TABLE cards
DROP COLUMN xp,
DROP COLUMN next_level_xp,
DROP COLUMN card_type,
DROP COLUMN individual_value,
DROP COLUMN rarity,
DROP COLUMN level,
DROP COLUMN buffs;