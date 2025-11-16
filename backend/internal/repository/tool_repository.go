package repository

import (
	"agent-platform/internal/model/ent"
	"agent-platform/internal/model/ent/tool"
	"context"
	"fmt"
)

// ToolRepository handles tool data access
type ToolRepository struct {
	client *ent.Client
}

// NewToolRepository creates a new tool repository
func NewToolRepository(client *ent.Client) *ToolRepository {
	return &ToolRepository{client: client}
}

// Create creates a new tool
func (r *ToolRepository) Create(ctx context.Context, t *ent.Tool) (*ent.Tool, error) {
	builder := r.client.Tool.
		Create().
		SetID(t.ID).
		SetName(t.Name).
		SetDescription(t.Description).
		SetType(t.Type).
		SetVersion(t.Version).
		SetIsPublic(t.IsPublic).
		SetCreatedBy(t.CreatedBy)

	// Set optional fields
	if t.Schema != nil {
		builder = builder.SetSchema(t.Schema)
	}
	if t.Implementation != "" {
		builder = builder.SetImplementation(t.Implementation)
	}
	if t.Tags != nil {
		builder = builder.SetTags(t.Tags)
	}
	if t.Category != "" {
		builder = builder.SetCategory(t.Category)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating tool: %w", err)
	}

	return created, nil
}

// Get retrieves a tool by ID
func (r *ToolRepository) Get(ctx context.Context, id string) (*ent.Tool, error) {
	t, err := r.client.Tool.
		Query().
		Where(tool.ID(id)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("tool not found: %s", id)
		}
		return nil, fmt.Errorf("failed querying tool: %w", err)
	}

	return t, nil
}

// List retrieves tools with pagination and filters
func (r *ToolRepository) List(ctx context.Context, page, pageSize int32, toolType, category, createdBy string, isPublic *bool) ([]*ent.Tool, int, error) {
	query := r.client.Tool.Query()

	// Apply filters
	if toolType != "" {
		query = query.Where(tool.Type(toolType))
	}
	if category != "" {
		query = query.Where(tool.Category(category))
	}
	if createdBy != "" {
		query = query.Where(tool.CreatedBy(createdBy))
	}
	if isPublic != nil {
		query = query.Where(tool.IsPublic(*isPublic))
	}

	// Get total count
	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed counting tools: %w", err)
	}

	// Apply pagination
	offset := int((page - 1) * pageSize)
	tools, err := query.
		Order(ent.Desc(tool.FieldCreatedAt)).
		Offset(offset).
		Limit(int(pageSize)).
		All(ctx)

	if err != nil {
		return nil, 0, fmt.Errorf("failed listing tools: %w", err)
	}

	return tools, total, nil
}

// Update updates an existing tool
func (r *ToolRepository) Update(ctx context.Context, id string, updates map[string]interface{}) (*ent.Tool, error) {
	updateQuery := r.client.Tool.UpdateOneID(id)

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
		case "schema":
			if v, ok := value.(map[string]interface{}); ok {
				updateQuery = updateQuery.SetSchema(v)
			}
		case "implementation":
			if v, ok := value.(string); ok {
				updateQuery = updateQuery.SetImplementation(v)
			}
		case "version":
			if v, ok := value.(string); ok {
				updateQuery = updateQuery.SetVersion(v)
			}
		case "is_public":
			if v, ok := value.(bool); ok {
				updateQuery = updateQuery.SetIsPublic(v)
			}
		case "tags":
			if v, ok := value.([]string); ok {
				updateQuery = updateQuery.SetTags(v)
			}
		case "category":
			if v, ok := value.(string); ok {
				updateQuery = updateQuery.SetCategory(v)
			}
		}
	}

	updated, err := updateQuery.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("tool not found: %s", id)
		}
		return nil, fmt.Errorf("failed updating tool: %w", err)
	}

	return updated, nil
}

// Delete deletes a tool by ID
func (r *ToolRepository) Delete(ctx context.Context, id string) error {
	err := r.client.Tool.
		DeleteOneID(id).
		Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("tool not found: %s", id)
		}
		return fmt.Errorf("failed deleting tool: %w", err)
	}

	return nil
}
