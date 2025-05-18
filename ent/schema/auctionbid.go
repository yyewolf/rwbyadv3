package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// AuctionBid holds the schema definition for the AuctionBid entity.
type AuctionBid struct {
	ent.Schema
}

// Fields of the AuctionBid.
func (AuctionBid) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()).Default(uuid.New),
		field.UUID("auction_id", uuid.New()),
		field.String("player_id"),
		field.Int64("price"),
	}
}

func (AuctionBid) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the AuctionBid.
func (AuctionBid) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("auction", Auction.Type).
			Ref("bids").
			Unique().
			Required().
			Field("auction_id"),
		edge.To("player", Player.Type).Unique().Required().Field("player_id"),
	}
}
