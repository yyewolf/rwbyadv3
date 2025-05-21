package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// PlayerDeck holds the schema definition for the PlayerDeck entity.
type PlayerDeck struct {
	ent.Schema
}

func (PlayerDeck) Annotations() []schema.Annotation {
	return []schema.Annotation{
		field.ID("player_id", "card_id"),
	}
}

// Fields of the PlayerCardDeck.
func (PlayerDeck) Fields() []ent.Field {
	return []ent.Field{
		field.String("player_id"),
		field.UUID("card_id", uuid.New()),
		field.Float("position"),
	}
}

func (PlayerDeck) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("card_id").
			Unique(),
	}
}

func (PlayerDeck) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the PlayerCardDeck.
func (PlayerDeck) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("player", Player.Type).
			Unique().
			Required().
			Field("player_id"),
		edge.To("card", Card.Type).
			Unique().
			Required().
			Field("card_id"),
	}
}
