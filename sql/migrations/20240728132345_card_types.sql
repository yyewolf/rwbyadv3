-- migrate:up
CREATE TABLE IF NOT EXISTS card_types (
    card_type varchar(50) PRIMARY KEY,

    name varchar(50) NOT NULL,
    categories varchar(50) NOT NULL
);

-- migrate:down
DROP TABLE IF EXISTS card_types;
