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

// PlayerFavoriteCards holds the schema definition for the PlayerFavoriteCards entity.
type PlayerFavoriteCards struct {
	ent.Schema
}

func (PlayerFavoriteCards) Annotations() []schema.Annotation {
	return []schema.Annotation{
		field.ID("player_id", "card_id"),
	}
}

// Fields of the PlayerCardFavorite.
func (PlayerFavoriteCards) Fields() []ent.Field {
	return []ent.Field{
		field.String("player_id"),
		field.UUID("card_id", uuid.New()),
		field.Float("position"),
	}
}

func (PlayerFavoriteCards) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("card_id").
			Unique(),
	}
}

func (PlayerFavoriteCards) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the PlayerCardFavorite.
func (PlayerFavoriteCards) Edges() []ent.Edge {
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
