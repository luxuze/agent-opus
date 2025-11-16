package repository

import (
	"agent-platform/internal/model/ent"
	"agent-platform/internal/model/ent/conversation"
	"context"
	"fmt"
	"time"
)

// ConversationRepository handles conversation data access
type ConversationRepository struct {
	client *ent.Client
}

// NewConversationRepository creates a new conversation repository
func NewConversationRepository(client *ent.Client) *ConversationRepository {
	return &ConversationRepository{client: client}
}

// Create creates a new conversation
func (r *ConversationRepository) Create(ctx context.Context, c *ent.Conversation) (*ent.Conversation, error) {
	builder := r.client.Conversation.
		Create().
		SetID(c.ID).
		SetAgentID(c.AgentID).
		SetUserID(c.UserID).
		SetStatus(c.Status)

	// Set optional fields
	if c.Title != "" {
		builder = builder.SetTitle(c.Title)
	}
	if c.Messages != nil {
		builder = builder.SetMessages(c.Messages)
	}
	if c.Context != nil {
		builder = builder.SetContext(c.Context)
	}
	if c.Metadata != nil {
		builder = builder.SetMetadata(c.Metadata)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating conversation: %w", err)
	}

	return created, nil
}

// Get retrieves a conversation by ID
func (r *ConversationRepository) Get(ctx context.Context, id string) (*ent.Conversation, error) {
	c, err := r.client.Conversation.
		Query().
		Where(conversation.ID(id)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("conversation not found: %s", id)
		}
		return nil, fmt.Errorf("failed querying conversation: %w", err)
	}

	return c, nil
}

// List retrieves conversations with pagination and filters
func (r *ConversationRepository) List(ctx context.Context, page, pageSize int32, agentID, userID, status string) ([]*ent.Conversation, int, error) {
	query := r.client.Conversation.Query()

	// Apply filters
	if agentID != "" {
		query = query.Where(conversation.AgentID(agentID))
	}
	if userID != "" {
		query = query.Where(conversation.UserID(userID))
	}
	if status != "" {
		query = query.Where(conversation.Status(status))
	}

	// Get total count
	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed counting conversations: %w", err)
	}

	// Apply pagination
	offset := int((page - 1) * pageSize)
	conversations, err := query.
		Order(ent.Desc(conversation.FieldUpdatedAt)).
		Offset(offset).
		Limit(int(pageSize)).
		All(ctx)

	if err != nil {
		return nil, 0, fmt.Errorf("failed listing conversations: %w", err)
	}

	return conversations, total, nil
}

// AddMessage adds a message to a conversation
func (r *ConversationRepository) AddMessage(ctx context.Context, id string, message interface{}) (*ent.Conversation, error) {
	c, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	messages := c.Messages
	if messages == nil {
		messages = []interface{}{}
	}
	messages = append(messages, message)

	updated, err := r.client.Conversation.
		UpdateOneID(id).
		SetMessages(messages).
		SetLastMessageAt(time.Now()).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed adding message: %w", err)
	}

	return updated, nil
}

// Update updates an existing conversation
func (r *ConversationRepository) Update(ctx context.Context, id string, updates map[string]interface{}) (*ent.Conversation, error) {
	updateQuery := r.client.Conversation.UpdateOneID(id)

	for key, value := range updates {
		switch key {
		case "title":
			if v, ok := value.(string); ok {
				updateQuery = updateQuery.SetTitle(v)
			}
		case "status":
			if v, ok := value.(string); ok {
				updateQuery = updateQuery.SetStatus(v)
			}
		case "context":
			if v, ok := value.(map[string]interface{}); ok {
				updateQuery = updateQuery.SetContext(v)
			}
		case "metadata":
			if v, ok := value.(map[string]interface{}); ok {
				updateQuery = updateQuery.SetMetadata(v)
			}
		}
	}

	updated, err := updateQuery.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("conversation not found: %s", id)
		}
		return nil, fmt.Errorf("failed updating conversation: %w", err)
	}

	return updated, nil
}

// Delete deletes a conversation by ID
func (r *ConversationRepository) Delete(ctx context.Context, id string) error {
	err := r.client.Conversation.
		DeleteOneID(id).
		Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("conversation not found: %s", id)
		}
		return fmt.Errorf("failed deleting conversation: %w", err)
	}

	return nil
}
