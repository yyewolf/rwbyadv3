-- migrate:up
CREATE TYPE loot_boxes_type AS ENUM ('classic', 'rare', 'limited', 'special');

CREATE TABLE IF NOT EXISTS loot_boxes (
    id varchar(50) PRIMARY KEY,
    player_id varchar(50) NOT NULL REFERENCES players (id),

    type loot_boxes_type NOT NULL,
    metadata json,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);

-- migrate:down
DROP TABLE loot_boxes;
DROP TYPE loot_boxes_type;
