package repository

import (
	"agent-platform/internal/model/ent"
	"agent-platform/internal/model/ent/agent"
	"context"
	"fmt"
)

// AgentRepository handles agent data access
type AgentRepository struct {
	client *ent.Client
}

// NewAgentRepository creates a new agent repository
func NewAgentRepository(client *ent.Client) *AgentRepository {
	return &AgentRepository{client: client}
}

// Create creates a new agent
func (r *AgentRepository) Create(ctx context.Context, a *ent.Agent) (*ent.Agent, error) {
	builder := r.client.Agent.
		Create().
		SetID(a.ID).
		SetName(a.Name).
		SetDescription(a.Description).
		SetType(a.Type).
		SetStatus(a.Status).
		SetVersion(a.Version).
		SetCreatedBy(a.CreatedBy).
		SetIsPublic(a.IsPublic)

	// Set optional fields
	if a.ModelConfig != nil {
		builder = builder.SetModelConfig(a.ModelConfig)
	}
	if a.Tools != nil {
		builder = builder.SetTools(a.Tools)
	}
	if a.KnowledgeBases != nil {
		builder = builder.SetKnowledgeBases(a.KnowledgeBases)
	}
	if a.PromptTemplate != "" {
		builder = builder.SetPromptTemplate(a.PromptTemplate)
	}
	if a.Parameters != nil {
		builder = builder.SetParameters(a.Parameters)
	}
	if a.Tags != nil {
		builder = builder.SetTags(a.Tags)
	}
	if a.Folder != "" {
		builder = builder.SetFolder(a.Folder)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating agent: %w", err)
	}

	return created, nil
}

// Get retrieves an agent by ID
func (r *AgentRepository) Get(ctx context.Context, id string) (*ent.Agent, error) {
	a, err := r.client.Agent.
		Query().
		Where(agent.ID(id)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("agent not found: %s", id)
		}
		return nil, fmt.Errorf("failed querying agent: %w", err)
	}

	return a, nil
}

// List retrieves agents with pagination and filters
func (r *AgentRepository) List(ctx context.Context, page, pageSize int32, status, agentType, createdBy string) ([]*ent.Agent, int, error) {
	query := r.client.Agent.Query()

	// Apply filters
	if status != "" {
		query = query.Where(agent.Status(status))
	}
	if agentType != "" {
		query = query.Where(agent.Type(agentType))
	}
	if createdBy != "" {
		query = query.Where(agent.CreatedBy(createdBy))
	}

	// Get total count
	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed counting agents: %w", err)
	}

	// Apply pagination
	offset := int((page - 1) * pageSize)
	agents, err := query.
		Order(ent.Desc(agent.FieldCreatedAt)).
		Offset(offset).
		Limit(int(pageSize)).
		All(ctx)

	if err != nil {
		return nil, 0, fmt.Errorf("failed listing agents: %w", err)
	}

	return agents, total, nil
}

// Update updates an existing agent
func (r *AgentRepository) Update(ctx context.Context, id string, updates map[string]interface{}) (*ent.Agent, error) {
	updateQuery := r.client.Agent.UpdateOneID(id)

	// Apply updates dynamically
	for key, value := range updates {
		switch key {
		case "name":
			if v, ok := value.(string); ok {
				updateQuery = updateQuery.SetName(v)
			}
		case "description":
			if v, ok := value.(string); ok {
				updateQuery = updateQuery.SetDescription(v)
			}
		case "type":
			if v, ok := value.(string); ok {
				updateQuery = updateQuery.SetType(v)
			}
		case "model_config":
			if v, ok := value.(map[string]interface{}); ok {
				updateQuery = updateQuery.SetModelConfig(v)
			}
		case "tools":
			if v, ok := value.([]string); ok {
				updateQuery = updateQuery.SetTools(v)
			}
		case "knowledge_bases":
			if v, ok := value.([]string); ok {
				updateQuery = updateQuery.SetKnowledgeBases(v)
			}
		case "prompt_template":
			if v, ok := value.(string); ok {
				updateQuery = updateQuery.SetPromptTemplate(v)
			}
		case "parameters":
			if v, ok := value.(map[string]interface{}); ok {
				updateQuery = updateQuery.SetParameters(v)
			}
		case "status":
			if v, ok := value.(string); ok {
				updateQuery = updateQuery.SetStatus(v)
			}
		case "tags":
			if v, ok := value.([]string); ok {
				updateQuery = updateQuery.SetTags(v)
			}
		case "folder":
			if v, ok := value.(string); ok {
				updateQuery = updateQuery.SetFolder(v)
			}
		case "is_public":
			if v, ok := value.(bool); ok {
				updateQuery = updateQuery.SetIsPublic(v)
			}
		}
	}

	updated, err := updateQuery.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("agent not found: %s", id)
		}
		return nil, fmt.Errorf("failed updating agent: %w", err)
	}

	return updated, nil
}

// Delete deletes an agent by ID
func (r *AgentRepository) Delete(ctx context.Context, id string) error {
	err := r.client.Agent.
		DeleteOneID(id).
		Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("agent not found: %s", id)
		}
		return fmt.Errorf("failed deleting agent: %w", err)
	}

	return nil
}
