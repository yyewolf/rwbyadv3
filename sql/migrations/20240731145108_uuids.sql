-- migrate:up

---- Delete all constraints
ALTER TABLE "auctions" DROP CONSTRAINT "auctions_card_id_fkey";
ALTER TABLE "auctions_bids" DROP CONSTRAINT "auctions_bids_auction_id_fkey";
ALTER TABLE "cards" DROP CONSTRAINT "fk_cards_stats";
ALTER TABLE "cards_stats" DROP CONSTRAINT "cards_stats_card_id_fkey";
ALTER TABLE "listings" DROP CONSTRAINT "listings_card_id_fkey";
ALTER TABLE "player_card_favorites" DROP CONSTRAINT "player_card_favorites_card_id_fkey";
ALTER TABLE "player_cards" DROP CONSTRAINT "player_cards_card_id_fkey";
ALTER TABLE "player_cards_deck" DROP CONSTRAINT "player_cards_deck_card_id_fkey";
ALTER TABLE "players" DROP CONSTRAINT "players_selected_card_id_fkey";
---- Delete all constraints

ALTER TABLE auctions
ALTER COLUMN id TYPE uuid USING id::uuid,
ALTER COLUMN card_id TYPE uuid USING card_id::uuid;

ALTER TABLE auctions_bids
ALTER COLUMN id TYPE uuid USING id::uuid,
ALTER COLUMN auction_id TYPE uuid USING auction_id::uuid;

ALTER TABLE cards
ALTER COLUMN id TYPE uuid USING id::uuid;

ALTER TABLE cards_stats
ALTER COLUMN card_id TYPE uuid USING card_id::uuid;

ALTER TABLE dungeons
ALTER COLUMN id TYPE uuid USING id::uuid;

ALTER TABLE jobs
ALTER COLUMN id TYPE uuid USING id::uuid;

ALTER TABLE listings
ALTER COLUMN id TYPE uuid USING id::uuid,
ALTER COLUMN card_id TYPE uuid USING card_id::uuid;

ALTER TABLE loot_boxes
ALTER COLUMN id TYPE uuid USING id::uuid;

ALTER TABLE player_card_favorites
ALTER COLUMN card_id TYPE uuid USING card_id::uuid;

ALTER TABLE player_cards
ALTER COLUMN card_id TYPE uuid USING card_id::uuid;

ALTER TABLE player_cards_deck
ALTER COLUMN card_id TYPE uuid USING card_id::uuid;

ALTER TABLE players
ALTER COLUMN selected_card_id TYPE uuid USING selected_card_id::uuid;


---- Create all constraints
ALTER TABLE "auctions" ADD CONSTRAINT "auctions_card_id_fkey" FOREIGN KEY (
    card_id
) REFERENCES cards (id) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "auctions_bids" ADD CONSTRAINT "auctions_bids_auction_id_fkey" FOREIGN KEY (
    auction_id
) REFERENCES auctions (id) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "cards" ADD CONSTRAINT "fk_cards_stats" FOREIGN KEY (
    id
) REFERENCES cards_stats (card_id) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "cards_stats" ADD CONSTRAINT "cards_stats_card_id_fkey" FOREIGN KEY (
    card_id
) REFERENCES cards (id) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "listings" ADD CONSTRAINT "listings_card_id_fkey" FOREIGN KEY (
    card_id
) REFERENCES cards (id) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "player_card_favorites" ADD CONSTRAINT "player_card_favorites_card_id_fkey" FOREIGN KEY (
    card_id
) REFERENCES cards (id) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "player_cards" ADD CONSTRAINT "player_cards_card_id_fkey" FOREIGN KEY (
    card_id
) REFERENCES cards (id) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "player_cards_deck" ADD CONSTRAINT "player_cards_deck_card_id_fkey" FOREIGN KEY (
    card_id
) REFERENCES cards (id) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "players" ADD CONSTRAINT "players_selected_card_id_fkey" FOREIGN KEY (
    selected_card_id
) REFERENCES cards (id) DEFERRABLE INITIALLY DEFERRED;
---- Create all constraints

-- migrate:down
