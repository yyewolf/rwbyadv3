package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Auction holds the schema definition for the Auction entity.
type Auction struct {
	ent.Schema
}

// Fields of the Auction.
func (Auction) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()).Default(uuid.New),
		field.String("player_id"),
		field.UUID("card_id", uuid.New()),
		field.Int("time_extensions").Default(0),
		field.Time("ends_at"),
	}
}

// Edges of the Auction.
func (Auction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owned_by", Player.Type).Ref("auctions").Unique().Required().Field("player_id"),
		edge.To("card", Card.Type).Unique().Required().Field("card_id"),
		edge.To("bids", AuctionBid.Type),
	}
}
