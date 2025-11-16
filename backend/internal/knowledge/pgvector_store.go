package knowledge

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"agent-platform/internal/model/ent"

	entsql "entgo.io/ent/dialect/sql"
	"go.uber.org/zap"
)

// PgVectorStore implements vector storage using PostgreSQL with pgvector extension
type PgVectorStore struct {
	client *ent.Client
	db     *sql.DB
	logger *zap.Logger
}

// NewPgVectorStore creates a new pgvector-based vector store
func NewPgVectorStore(client *ent.Client, dsn string, logger *zap.Logger) (*PgVectorStore, error) {
	// Open a separate database connection for raw SQL queries
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	store := &PgVectorStore{
		client: client,
		db:     db,
		logger: logger,
	}

	// Ensure pgvector extension is enabled
	if err := store.ensurePgVector(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ensure pgvector: %w", err)
	}

	return store, nil
}

// ensurePgVector ensures the pgvector extension is installed
func (s *PgVectorStore) ensurePgVector(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, "CREATE EXTENSION IF NOT EXISTS vector")
	if err != nil {
		return fmt.Errorf("failed to create vector extension: %w", err)
	}

	s.logger.Info("pgvector extension ensured")
	return nil
}

// AddChunks implements VectorStore interface - adds multiple chunks
func (s *PgVectorStore) AddChunks(kbID string, chunks []*Chunk) error {
	ctx := context.Background()

	for _, chunk := range chunks {
		if err := s.Store(ctx, kbID, chunk); err != nil {
			return fmt.Errorf("failed to add chunk %s: %w", chunk.ID, err)
		}
	}

	s.logger.Info("Added chunks to pgvector store",
		zap.String("kb_id", kbID),
		zap.Int("count", len(chunks)),
	)

	return nil
}

// Store stores a single chunk with its embedding
func (s *PgVectorStore) Store(ctx context.Context, kbID string, chunk *Chunk) error {
	// Store in database using ent
	_, err := s.client.DocumentChunk.Create().
		SetID(chunk.ID).
		SetKnowledgeBaseID(kbID).
		SetDocumentID(chunk.DocumentID).
		SetChunkIndex(chunk.Index).
		SetContent(chunk.Content).
		SetEmbedding(chunk.Embedding).
		SetMetadata(chunk.Metadata).
		SetCreatedAt(time.Now()).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to store chunk: %w", err)
	}

	s.logger.Debug("Stored chunk with embedding",
		zap.String("chunk_id", chunk.ID),
		zap.String("kb_id", kbID),
		zap.Int("embedding_dim", len(chunk.Embedding)),
	)

	return nil
}

// Search performs vector similarity search using pgvector
func (s *PgVectorStore) Search(kbID string, queryEmbedding []float32, topK int, threshold float64) ([]*SearchResult, error) {
	ctx := context.Background()

	// Convert embedding to vector format for PostgreSQL
	embeddingJSON, err := json.Marshal(queryEmbedding)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query embedding: %w", err)
	}

	// Use cosine distance (<=> operator in pgvector)
	// Note: We're using embedding as JSON, so we need to cast it
	query := `
		SELECT
			id,
			knowledge_base_id,
			document_id,
			chunk_index,
			content,
			embedding,
			metadata,
			created_at,
			1 - (embedding::vector <=> $1::vector) as similarity
		FROM document_chunks
		WHERE knowledge_base_id = $2
		  AND 1 - (embedding::vector <=> $1::vector) >= $3
		ORDER BY embedding::vector <=> $1::vector
		LIMIT $4
	`

	rows, err := s.db.QueryContext(ctx, query, string(embeddingJSON), kbID, threshold, topK)
	if err != nil {
		return nil, fmt.Errorf("failed to execute similarity search: %w", err)
	}
	defer rows.Close()

	var results []*SearchResult
	for rows.Next() {
		var (
			id              string
			knowledgeBaseID string
			documentID      string
			chunkIndex      int
			content         string
			embeddingBytes  []byte
			metadataBytes   []byte
			createdAt       time.Time
			similarity      float64
		)

		err := rows.Scan(
			&id,
			&knowledgeBaseID,
			&documentID,
			&chunkIndex,
			&content,
			&embeddingBytes,
			&metadataBytes,
			&createdAt,
			&similarity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Parse embedding
		var embedding []float32
		if err := json.Unmarshal(embeddingBytes, &embedding); err != nil {
			s.logger.Warn("Failed to unmarshal embedding, skipping", zap.Error(err))
			continue
		}

		// Parse metadata
		var metadata map[string]interface{}
		if len(metadataBytes) > 0 {
			if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
				s.logger.Warn("Failed to unmarshal metadata", zap.Error(err))
				metadata = make(map[string]interface{})
			}
		}

		chunk := &Chunk{
			ID:         id,
			DocumentID: documentID,
			Index:      chunkIndex,
			Content:    content,
			Embedding:  embedding,
			Metadata:   metadata,
		}

		results = append(results, &SearchResult{
			Chunk: chunk,
			Score: similarity,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	s.logger.Debug("Vector similarity search completed",
		zap.String("kb_id", kbID),
		zap.Int("results", len(results)),
		zap.Int("top_k", topK),
		zap.Float64("threshold", threshold),
	)

	return results, nil
}

// DeleteKnowledgeBase deletes all chunks for a knowledge base
func (s *PgVectorStore) DeleteKnowledgeBase(kbID string) error {
	ctx := context.Background()
	_, err := s.client.DocumentChunk.Delete().
		Where(func(s *entsql.Selector) {
			s.Where(entsql.EQ("knowledge_base_id", kbID))
		}).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete chunks: %w", err)
	}

	s.logger.Info("Deleted all chunks for knowledge base", zap.String("kb_id", kbID))
	return nil
}

// DeleteDocument deletes all chunks for a specific document
func (s *PgVectorStore) DeleteDocument(ctx context.Context, kbID, documentID string) error {
	_, err := s.client.DocumentChunk.Delete().
		Where(func(s *entsql.Selector) {
			s.Where(entsql.And(
				entsql.EQ("knowledge_base_id", kbID),
				entsql.EQ("document_id", documentID),
			))
		}).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete document chunks: %w", err)
	}

	s.logger.Info("Deleted chunks for document",
		zap.String("kb_id", kbID),
		zap.String("document_id", documentID),
	)
	return nil
}

// GetStats returns the number of chunks in a knowledge base
func (s *PgVectorStore) GetStats(kbID string) (int, error) {
	ctx := context.Background()
	count, err := s.client.DocumentChunk.Query().
		Where(func(s *entsql.Selector) {
			s.Where(entsql.EQ("knowledge_base_id", kbID))
		}).
		Count(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to count chunks: %w", err)
	}

	return count, nil
}
