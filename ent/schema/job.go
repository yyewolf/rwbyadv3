package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// Job holds the schema definition for the Job entity.
type Job struct {
	ent.Schema
}

// Fields of the Job.
func (Job) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()).Default(uuid.New),
		field.String("jobkey"),
		field.Int("retries").Default(0).NonNegative(),
		field.Time("run_at"),
		field.JSON("params", map[string]any{}),
		field.Int64("last_run_id").Default(0).NonNegative(),
		field.Bool("recurring").Default(false),
		field.Int64("delta_time").Default(0).NonNegative(),
		field.Bool("errored").Default(false),
	}
}

func (Job) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the Job.
func (Job) Edges() []ent.Edge {
	return nil
}
