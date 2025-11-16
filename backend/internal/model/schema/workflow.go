package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Workflow holds the schema definition for the Workflow entity.
type Workflow struct {
	ent.Schema
}

// Fields of the Workflow.
func (Workflow) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable(),
		field.String("name").
			NotEmpty(),
		field.String("description").
			Optional(),
		field.JSON("steps", []map[string]interface{}{}).
			Comment("Workflow steps configuration"),
		field.JSON("config", map[string]interface{}{}).
			Optional().
			Comment("Workflow configuration"),
		field.String("status").
			Default("draft").
			Comment("draft, active, archived"),
		field.String("created_by").
			NotEmpty(),
		field.Time("created_at"),
		field.Time("updated_at"),
	}
}

// Indexes of the Workflow.
func (Workflow) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("created_by"),
		index.Fields("status"),
	}
}
