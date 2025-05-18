package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// Dungeon holds the schema definition for the Dungeon entity.
type Dungeon struct {
	ent.Schema
}

// Fields of the Dungeon.
func (Dungeon) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()).Default(uuid.New),
		field.String("player_id"),
		field.Int64("seed"),
	}
}

func (Dungeon) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the Dungeon.
func (Dungeon) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owned_by", Player.Type).
			Ref("dungeons").
			Unique().
			Required().
			Field("player_id"),
	}
}
