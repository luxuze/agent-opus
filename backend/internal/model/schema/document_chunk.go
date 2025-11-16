package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// DocumentChunk holds the schema definition for the DocumentChunk entity.
type DocumentChunk struct {
	ent.Schema
}

// Fields of the DocumentChunk.
func (DocumentChunk) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable(),
		field.String("knowledge_base_id").
			NotEmpty(),
		field.String("document_id").
			NotEmpty(),
		field.Int("chunk_index").
			Default(0),
		field.Text("content").
			NotEmpty(),
		field.JSON("embedding", []float32{}).
			Optional().
			Comment("Vector embedding stored as JSON, converted to vector type in PostgreSQL"),
		field.JSON("metadata", map[string]interface{}{}).
			Optional(),
		field.Time("created_at"),
	}
}

// Indexes of the DocumentChunk.
func (DocumentChunk) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("knowledge_base_id"),
		index.Fields("document_id"),
		index.Fields("knowledge_base_id", "document_id"),
	}
}
