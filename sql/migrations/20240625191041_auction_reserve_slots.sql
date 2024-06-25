-- migrate:up
ALTER TABLE players
ADD COLUMN slots_reserved int NOT NULL DEFAULT 0;

-- migrate:down
ALTER TABLE players
DROP COLUMN slots_reserved;