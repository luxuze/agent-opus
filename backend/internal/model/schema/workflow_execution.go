package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// WorkflowExecution holds the schema definition for the WorkflowExecution entity.
type WorkflowExecution struct {
	ent.Schema
}

// Fields of the WorkflowExecution.
func (WorkflowExecution) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable(),
		field.String("workflow_id").
			NotEmpty(),
		field.JSON("input", map[string]interface{}{}).
			Optional().
			Comment("Initial workflow input"),
		field.JSON("output", map[string]interface{}{}).
			Optional().
			Comment("Final workflow output"),
		field.JSON("context", map[string]interface{}{}).
			Optional().
			Comment("Execution context and intermediate results"),
		field.String("status").
			Default("pending").
			Comment("pending, running, completed, failed, cancelled"),
		field.String("current_step").
			Optional().
			Comment("Current step being executed"),
		field.String("error").
			Optional().
			Comment("Error message if failed"),
		field.String("started_by").
			NotEmpty(),
		field.Time("started_at"),
		field.Time("completed_at").
			Optional(),
		field.Time("updated_at"),
	}
}

// Indexes of the WorkflowExecution.
func (WorkflowExecution) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("workflow_id"),
		index.Fields("started_by"),
		index.Fields("status"),
	}
}
