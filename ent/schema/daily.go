package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// Daily holds the schema definition for the Daily entity.
type Daily struct {
	ent.Schema
}

// Fields of the Daily.
func (Daily) Fields() []ent.Field {
	return []ent.Field{
		field.String("player_id"),
		field.Bool("has_voted").Default(false),
		field.Time("last_vote_at").Default(time.Now),
		field.Int64("streak").Default(0),
	}
}

func (Daily) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the Daily.
func (Daily) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("of", Player.Type).
			Ref("daily").
			Unique().
			Required().
			Field("player_id"),
	}
}
