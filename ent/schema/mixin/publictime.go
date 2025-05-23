package mixin

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/yyewolf/entvis"
)

type PublicTime struct{ mixin.Time }

// Fields of the time mixin.
func (PublicTime) Fields() []ent.Field {
	return []ent.Field{
		field.Time("create_time").
			Default(time.Now).
			Immutable().
			Annotations(entvis.Visibility("public")),
		field.Time("update_time").
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(entvis.Visibility("public")),
	}
}
