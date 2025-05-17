package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
	"github.com/yyewolf/rwbyadv3/ent/schema/enums"
)

// AuthState holds the schema definition for the AuthState entity.
type AuthState struct {
	ent.Schema
}

// Fields of the AuthState.
func (AuthState) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()).Default(uuid.New),
		field.String("player_id"),
		field.String("redirect_uri").Optional(),
		field.Time("expires_at"),
		field.Enum("type").GoType(enums.AuthStateTypes("")),
	}
}

func (AuthState) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the AuthState.
func (AuthState) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("player", Player.Type).Unique().Required().Field("player_id"),
	}
}
