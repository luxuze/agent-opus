package repository

import (
	"agent-platform/internal/model/ent"
	"agent-platform/internal/model/ent/knowledgebase"
	"context"
	"fmt"
)

// KnowledgeBaseRepository handles knowledge base data access
type KnowledgeBaseRepository struct {
	client *ent.Client
}

// NewKnowledgeBaseRepository creates a new knowledge base repository
func NewKnowledgeBaseRepository(client *ent.Client) *KnowledgeBaseRepository {
	return &KnowledgeBaseRepository{client: client}
}

// Create creates a new knowledge base
func (r *KnowledgeBaseRepository) Create(ctx context.Context, kb *ent.KnowledgeBase) (*ent.KnowledgeBase, error) {
	builder := r.client.KnowledgeBase.
		Create().
		SetID(kb.ID).
		SetName(kb.Name).
		SetDescription(kb.Description).
		SetType(kb.Type).
		SetEmbeddingModel(kb.EmbeddingModel).
		SetCreatedBy(kb.CreatedBy).
		SetDocumentCount(kb.DocumentCount).
		SetVectorCount(kb.VectorCount)

	// Set optional fields
	if kb.ChunkConfig != nil {
		builder = builder.SetChunkConfig(kb.ChunkConfig)
	}
	if kb.Documents != nil {
		builder = builder.SetDocuments(kb.Documents)
	}
	if kb.Metadata != nil {
		builder = builder.SetMetadata(kb.Metadata)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating knowledge base: %w", err)
	}

	return created, nil
}

// Get retrieves a knowledge base by ID
func (r *KnowledgeBaseRepository) Get(ctx context.Context, id string) (*ent.KnowledgeBase, error) {
	kb, err := r.client.KnowledgeBase.
		Query().
		Where(knowledgebase.ID(id)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("knowledge base not found: %s", id)
		}
		return nil, fmt.Errorf("failed querying knowledge base: %w", err)
	}

	return kb, nil
}

// List retrieves knowledge bases with pagination and filters
func (r *KnowledgeBaseRepository) List(ctx context.Context, page, pageSize int32, kbType, createdBy string) ([]*ent.KnowledgeBase, int, error) {
	query := r.client.KnowledgeBase.Query()

	// Apply filters
	if kbType != "" {
		query = query.Where(knowledgebase.Type(kbType))
	}
	if createdBy != "" {
		query = query.Where(knowledgebase.CreatedBy(createdBy))
	}

	// Get total count
	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed counting knowledge bases: %w", err)
	}

	// Apply pagination
	offset := int((page - 1) * pageSize)
	kbs, err := query.
		Order(ent.Desc(knowledgebase.FieldCreatedAt)).
		Offset(offset).
		Limit(int(pageSize)).
		All(ctx)

	if err != nil {
		return nil, 0, fmt.Errorf("failed listing knowledge bases: %w", err)
	}

	return kbs, total, nil
}

// Update updates an existing knowledge base
func (r *KnowledgeBaseRepository) Update(ctx context.Context, id string, updates map[string]interface{}) (*ent.KnowledgeBase, error) {
	updateQuery := r.client.KnowledgeBase.UpdateOneID(id)

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
		case "embedding_model":
			if v, ok := value.(string); ok {
				updateQuery = updateQuery.SetEmbeddingModel(v)
			}
		case "chunk_config":
			if v, ok := value.(map[string]interface{}); ok {
				updateQuery = updateQuery.SetChunkConfig(v)
			}
		case "documents":
			if v, ok := value.([]interface{}); ok {
				updateQuery = updateQuery.SetDocuments(v)
			}
		case "metadata":
			if v, ok := value.(map[string]interface{}); ok {
				updateQuery = updateQuery.SetMetadata(v)
			}
		case "document_count":
			if v, ok := value.(int); ok {
				updateQuery = updateQuery.SetDocumentCount(v)
			}
		case "vector_count":
			if v, ok := value.(int); ok {
				updateQuery = updateQuery.SetVectorCount(v)
			}
		}
	}

	updated, err := updateQuery.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("knowledge base not found: %s", id)
		}
		return nil, fmt.Errorf("failed updating knowledge base: %w", err)
	}

	return updated, nil
}

// Delete deletes a knowledge base by ID
func (r *KnowledgeBaseRepository) Delete(ctx context.Context, id string) error {
	err := r.client.KnowledgeBase.
		DeleteOneID(id).
		Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("knowledge base not found: %s", id)
		}
		return fmt.Errorf("failed deleting knowledge base: %w", err)
	}

	return nil
}
