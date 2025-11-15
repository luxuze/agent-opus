package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// KnowledgeBase holds the schema definition for the KnowledgeBase entity.
type KnowledgeBase struct {
	ent.Schema
}

// Fields of the KnowledgeBase.
func (KnowledgeBase) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable(),
		field.String("name").
			NotEmpty(),
		field.Text("description").
			Optional(),
		field.String("type").
			Default("document"), // document, database, api
		field.String("embedding_model").
			Default("text-embedding-ada-002"),
		field.JSON("chunk_config", map[string]interface{}{}).
			Optional(),
		field.JSON("documents", []interface{}{}).
			Optional(),
		field.JSON("metadata", map[string]interface{}{}).
			Optional(),
		field.String("created_by").
			NotEmpty(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Int("document_count").
			Default(0),
		field.Int("vector_count").
			Default(0),
	}
}

// Indexes of the KnowledgeBase.
func (KnowledgeBase) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("type"),
		index.Fields("created_by"),
		index.Fields("created_at"),
	}
}
