package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// GithubStar holds the schema definition for the GithubStar entity.
type GithubStar struct {
	ent.Schema
}

// Fields of the GithubStar.
func (GithubStar) Fields() []ent.Field {
	return []ent.Field{
		field.String("player_id"),
		field.String("github_user_id").Optional(),
		field.Bool("has_starred").Default(false),
	}
}

func (GithubStar) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the GithubStar.
func (GithubStar) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("of", Player.Type).
			Ref("github_star").
			Unique().
			Required().
			Field("player_id"),
	}
}
