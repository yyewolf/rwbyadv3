package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
	"github.com/yyewolf/rwbyadv3/ent/schema/enums"
)

// LootBox holds the schema definition for the LootBox entity.
type LootBox struct {
	ent.Schema
}

// Fields of the LootBox.
func (LootBox) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()).Default(uuid.New),
		field.String("player_id"),
		field.Enum("type").
			GoType(enums.LootBoxType("")),
	}
}

func (LootBox) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the LootBox.
func (LootBox) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owned_by", Player.Type).
			Ref("lootboxes").
			Unique().
			Required().
			Field("player_id"),
	}
}
