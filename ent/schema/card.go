package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// Card holds the schema definition for the Card entity.
type Card struct {
	ent.Schema
}

type CardMetadata struct {
	Location string
}

// Fields of the Card.
func (Card) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()).Default(uuid.New),
		field.String("player_id"),
		field.String("card_type"),
		field.Float("position"),
		field.Bool("available").Default(true),

		field.Int64("level").Default(1).Positive(),
		field.Int64("experience_points").Default(0).NonNegative(),
		field.Int64("experience_points_threshold").Default(0).NonNegative(),

		field.Int("rarity").NonNegative(),
		field.Int("buffs").NonNegative(),

		field.Float("individual_value").Positive(),

		field.JSON("metadata", CardMetadata{}).Default(CardMetadata{}),
		field.Time("owned_at").Default(time.Now),
	}
}

func (Card) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the Card.
func (Card) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owned_by", Player.Type).
			Ref("cards").
			Unique().
			Required().
			Field("player_id"),
		edge.From("in_favorites_of", Player.Type).
			Ref("favorite_cards").
			Through("player_favorite_cards", PlayerFavoriteCards.Type),
		edge.From("in_deck_of", Player.Type).
			Ref("deck").
			Through("player_decks", PlayerDeck.Type),
		edge.To("type", CardType.Type).Unique().Required().Field("card_type"),
		edge.To("stats", CardStats.Type).Unique(),
	}
}
