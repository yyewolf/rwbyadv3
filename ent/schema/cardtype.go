package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// CardType holds the schema definition for the CardType entity.
type CardType struct {
	ent.Schema
}

// Fields of the CardType.
func (CardType) Fields() []ent.Field {
	return []ent.Field{
		field.String("id"),
		field.String("name"),
		field.Strings("categories"),
		field.String("s_categories"),
	}
}

// Edges of the CardType.
func (CardType) Edges() []ent.Edge {
	return nil
}
