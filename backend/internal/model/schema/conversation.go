package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Conversation holds the schema definition for the Conversation entity.
type Conversation struct {
	ent.Schema
}

// Fields of the Conversation.
func (Conversation) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable(),
		field.String("agent_id").
			NotEmpty(),
		field.String("user_id").
			NotEmpty(),
		field.String("title").
			Optional(),
		field.JSON("messages", []interface{}{}).
			Optional(),
		field.JSON("context", map[string]interface{}{}).
			Optional(),
		field.JSON("metadata", map[string]interface{}{}).
			Optional(),
		field.String("status").
			Default("active"), // active, closed
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Time("last_message_at").
			Optional(),
	}
}

// Indexes of the Conversation.
func (Conversation) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("agent_id"),
		index.Fields("user_id"),
		index.Fields("status"),
		index.Fields("created_at"),
	}
}
