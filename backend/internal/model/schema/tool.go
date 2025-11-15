package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Tool holds the schema definition for the Tool entity.
type Tool struct {
	ent.Schema
}

// Fields of the Tool.
func (Tool) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable(),
		field.String("name").
			NotEmpty(),
		field.Text("description").
			Optional(),
		field.String("type").
			Default("function"), // function, api, plugin
		field.JSON("schema", map[string]interface{}{}).
			Optional(),
		field.Text("implementation").
			Optional(),
		field.String("version").
			Default("1.0.0"),
		field.Bool("is_public").
			Default(false),
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
		field.String("category").
			Optional(),
	}
}

// Indexes of the Tool.
func (Tool) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("type"),
		index.Fields("created_by"),
		index.Fields("is_public"),
	}
}
