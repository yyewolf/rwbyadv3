package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// CardStats holds the schema definition for the CardStats entity.
type CardStats struct {
	ent.Schema
}

// Fields of the CardStats.
func (CardStats) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("card_id", uuid.New()).Default(uuid.New),
		field.Int64("health"),
		field.Int64("armor"),
		field.Int64("damage"),
		field.Int64("healing"),
		field.Int64("speed"),
	}
}

func (CardStats) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the CardStats.
func (CardStats) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("of", Card.Type).
			Ref("stats").
			Unique().
			Required().
			Field("card_id"),
	}
}
