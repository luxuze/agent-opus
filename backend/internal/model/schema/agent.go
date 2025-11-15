package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Agent holds the schema definition for the Agent entity.
type Agent struct {
	ent.Schema
}

// Fields of the Agent.
func (Agent) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable(),
		field.String("name").
			NotEmpty(),
		field.Text("description").
			Optional(),
		field.String("type").
			Default("single"), // single, multi
		field.JSON("model_config", map[string]interface{}{}).
			Optional(),
		field.JSON("tools", []string{}).
			Optional(),
		field.JSON("knowledge_bases", []string{}).
			Optional(),
		field.Text("prompt_template").
			Optional(),
		field.JSON("parameters", map[string]interface{}{}).
			Optional(),
		field.String("status").
			Default("draft"), // draft, published, archived
		field.String("version").
			Default("1.0.0"),
		field.String("created_by").
			NotEmpty(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.JSON("tags", []string{}).
			Optional(),
		field.String("folder").
			Optional(),
		field.Bool("is_public").
			Default(false),
	}
}

// Indexes of the Agent.
func (Agent) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("created_by"),
		index.Fields("status"),
		index.Fields("created_at"),
	}
}
