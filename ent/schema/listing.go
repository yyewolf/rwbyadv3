package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// Listing holds the schema definition for the Listing entity.
type Listing struct {
	ent.Schema
}

// Fields of the Listing.
func (Listing) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()).Default(uuid.New),
		field.String("player_id"),
		field.UUID("card_id", uuid.New()),
		field.Int64("price").NonNegative(),
		field.String("note"),
	}
}

func (Listing) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the Listing.
func (Listing) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owned_by", Player.Type).Ref("listings").Unique().Required().Field("player_id"),
		edge.To("card", Card.Type).Unique().Required().Field("card_id"),
	}
}
