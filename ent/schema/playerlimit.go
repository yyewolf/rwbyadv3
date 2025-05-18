package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// PlayerLimit holds the schema definition for the PlayerLimit entity.
type PlayerLimit struct {
	ent.Schema
}

// Fields of the PlayerLimit.
func (PlayerLimit) Fields() []ent.Field {
	return []ent.Field{
		field.String("player_id"),
		field.Int("dungeons_left").Default(3).NonNegative(),
		field.Time("dungeons_reset_at").Optional(),
	}
}

func (PlayerLimit) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the PlayerLimit.
func (PlayerLimit) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("of", Player.Type).
			Ref("limits").
			Unique().
			Required().
			Field("player_id"),
	}
}
