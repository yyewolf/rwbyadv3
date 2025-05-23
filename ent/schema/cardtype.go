package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/yyewolf/entvis"
)

// CardType holds the schema definition for the CardType entity.
type CardType struct {
	ent.Schema
}

// Fields of the CardType.
func (CardType) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Annotations(entvis.Visibility(RolePublic)),
		field.String("name").Annotations(entvis.Visibility(RolePublic)),
		field.Strings("categories").Annotations(entvis.Visibility(RolePublic)),
		field.String("s_categories"),
	}
}

// Edges of the CardType.
func (CardType) Edges() []ent.Edge {
	return nil
}
