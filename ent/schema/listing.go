package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/yyewolf/entvis"
	"github.com/yyewolf/rwbyadv3/ent/schema/mixin"
)

// Listing holds the schema definition for the Listing entity.
type Listing struct {
	ent.Schema
}

// Fields of the Listing.
func (Listing) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()).Default(uuid.New).Annotations(entvis.Visibility(RoleSelf, RolePublic)),
		field.String("player_id"),
		field.UUID("card_id", uuid.New()).Annotations(entvis.Visibility(RoleSelf, RolePublic)),
		field.Int64("price").NonNegative().Annotations(entvis.Visibility(RoleSelf, RolePublic)),
		field.String("note").Optional().Annotations(entvis.Visibility(RoleSelf, RolePublic)),
	}
}

func (Listing) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.PublicTime{},
	}
}

// Edges of the Listing.
func (Listing) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owned_by", Player.Type).Ref("listings").Unique().Required().Field("player_id"),
		edge.To("card", Card.Type).Unique().Required().Field("card_id"),
	}
}
