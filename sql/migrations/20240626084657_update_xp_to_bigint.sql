-- migrate:up
ALTER TABLE cards
ALTER COLUMN xp TYPE bigint,
ALTER COLUMN next_level_xp TYPE bigint;

-- migrate:down
ALTER TABLE cards
ALTER COLUMN xp TYPE bigint,
ALTER COLUMN next_level_xp TYPE bigint;
