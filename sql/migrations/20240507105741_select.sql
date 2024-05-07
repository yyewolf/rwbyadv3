-- migrate:up
ALTER TABLE players
ADD COLUMN selected_card_id VARCHAR(50) REFERENCES cards (id);

CREATE TABLE IF NOT EXISTS player_cards_deck (
    player_id VARCHAR(50) REFERENCES players (id),
    card_id VARCHAR(50) REFERENCES cards (id),
    position INT NOT NULL,
    PRIMARY KEY (player_id, card_id)
);

CREATE TABLE IF NOT EXISTS player_card_favorites (
    player_id VARCHAR(50) REFERENCES players (id),
    card_id VARCHAR(50) REFERENCES cards (id),
    position INT NOT NULL,
    PRIMARY KEY (player_id, card_id)
);

-- migrate:down
DROP TABLE player_cards_deck;
DROP TABLE player_card_favorites;

ALTER TABLE players
DROP COLUMN selected_card_id;
