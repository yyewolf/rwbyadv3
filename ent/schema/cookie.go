package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// Cookie holds the schema definition for the Cookie entity.
type Cookie struct {
	ent.Schema
}

// Fields of the Cookie.
func (Cookie) Fields() []ent.Field {
	return []ent.Field{
		field.String("id"),
		field.String("player_id"),
		field.Time("expires_at"),
	}
}

func (Cookie) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the Cookie.
func (Cookie) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("player", Player.Type).Unique().Required().Field("player_id"),
	}
}
