package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// Player holds the schema definition for the Player entity.
type Player struct {
	ent.Schema
}

// Fields of the Player.
func (Player) Fields() []ent.Field {
	return []ent.Field{
		field.String("id"),
		field.String("username").Default(""),

		field.Int64("liens").Default(500).NonNegative(),
		field.Int64("liens_in_auction").Default(0).NonNegative(),

		field.Int64("level").Default(1).Positive(),
		field.Int64("experience_points").Default(0).NonNegative(),
		field.Int64("experience_points_threshold").Default(0).NonNegative(),

		field.Int64("backpack_level").Default(1).Positive(),
		field.Int64("backpack_reserved_slots").Default(0).NonNegative(),

		// Foreign Keys
		field.UUID("selected_card_id", uuid.New()).Optional(),
	}
}

func (Player) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the Player.
func (Player) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("limits", PlayerLimit.Type).Unique(),
		edge.To("cards", Card.Type),
		edge.To("favorite_cards", Card.Type).Through("player_favorite_cards", PlayerFavoriteCards.Type),
		edge.To("deck", Card.Type).Through("player_decks", PlayerDeck.Type),
		edge.To("lootboxes", LootBox.Type),
		edge.To("selected_card", Card.Type).Unique().Field("selected_card_id"),
		edge.To("github_star", GithubStar.Type).Unique(),
		edge.To("daily", Daily.Type).Unique(),
		edge.To("listings", Listing.Type),
		edge.To("dungeons", Dungeon.Type),
	}
}
