package knowledge

import (
	"fmt"
	"time"

	"agent-platform/internal/model/ent"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Manager manages knowledge bases, documents, and vector search
type Manager struct {
	chunker       *Chunker
	embedder      *EmbeddingService
	vectorStore   VectorStore
	documentStore *DocumentStore
	logger        *zap.Logger
}

// NewManager creates a new knowledge base manager
func NewManager(client *ent.Client, dsn, apiKey, apiBase, embeddingModel string, logger *zap.Logger) (*Manager, error) {
	chunker := NewChunker(DefaultChunkConfig())

	var embedder *EmbeddingService
	if apiKey != "" && apiKey != "your-openai-api-key" {
		embedder = NewEmbeddingService(apiKey, apiBase, embeddingModel, logger)
		logger.Info("Embedding service initialized",
			zap.String("model", embeddingModel),
			zap.String("api_base", apiBase),
		)
	} else {
		logger.Warn("Embedding API key not configured, embedding features will be limited")
	}

	// Use PostgreSQL with pgvector for vector storage
	vectorStore, err := NewPgVectorStore(client, dsn, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgvector store: %w", err)
	}

	documentStore := NewDocumentStore(logger)

	return &Manager{
		chunker:       chunker,
		embedder:      embedder,
		vectorStore:   vectorStore,
		documentStore: documentStore,
		logger:        logger,
	}, nil
}

// AddDocument adds a document to a knowledge base
func (m *Manager) AddDocument(kbID, title, content, contentType, source string, metadata map[string]interface{}) (*Document, error) {
	// Create document
	doc := &Document{
		ID:          uuid.New().String(),
		Title:       title,
		Content:     content,
		ContentType: contentType,
		Source:      source,
		Metadata:    metadata,
		UploadedAt:  time.Now(),
	}

	// Chunk the document
	chunks, err := m.chunker.ChunkDocument(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to chunk document: %w", err)
	}

	doc.ChunkCount = len(chunks)

	// Generate embeddings if available
	if m.embedder != nil {
		m.logger.Info("Generating embeddings for document",
			zap.String("doc_id", doc.ID),
			zap.Int("chunks", len(chunks)),
		)

		err = m.embedder.EmbedChunks(chunks)
		if err != nil {
			m.logger.Error("Failed to generate embeddings", zap.Error(err))
			// Continue without embeddings - they can be generated later
		}
	}

	// Store document and chunks
	err = m.documentStore.AddDocument(kbID, doc, chunks)
	if err != nil {
		return nil, fmt.Errorf("failed to store document: %w", err)
	}

	// Add chunks to vector store
	if m.embedder != nil {
		err = m.vectorStore.AddChunks(kbID, chunks)
		if err != nil {
			return nil, fmt.Errorf("failed to add chunks to vector store: %w", err)
		}
	}

	m.logger.Info("Document added successfully",
		zap.String("kb_id", kbID),
		zap.String("doc_id", doc.ID),
		zap.Int("chunks", len(chunks)),
	)

	return doc, nil
}

// Search performs semantic search in a knowledge base
func (m *Manager) Search(kbID, query string, topK int, threshold float64) ([]*SearchResult, error) {
	if m.embedder == nil {
		return nil, fmt.Errorf("embedding service not available - please configure OPENAI_API_KEY")
	}

	// Generate embedding for query
	m.logger.Info("Searching knowledge base",
		zap.String("kb_id", kbID),
		zap.String("query", query),
		zap.Int("top_k", topK),
	)

	queryEmbedding, err := m.embedder.GenerateEmbedding(query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Search vector store
	results, err := m.vectorStore.Search(kbID, queryEmbedding, topK, threshold)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	m.logger.Info("Search completed",
		zap.String("kb_id", kbID),
		zap.Int("results", len(results)),
	)

	return results, nil
}

// GetDocument retrieves a document
func (m *Manager) GetDocument(kbID, docID string) (*Document, error) {
	return m.documentStore.GetDocument(kbID, docID)
}

// ListDocuments lists all documents in a knowledge base
func (m *Manager) ListDocuments(kbID string) ([]*Document, error) {
	return m.documentStore.ListDocuments(kbID)
}

// DeleteDocument removes a document from a knowledge base
func (m *Manager) DeleteDocument(kbID, docID string) error {
	return m.documentStore.DeleteDocument(kbID, docID)
}

// DeleteKnowledgeBase removes all data for a knowledge base
func (m *Manager) DeleteKnowledgeBase(kbID string) error {
	// Delete from vector store
	err := m.vectorStore.DeleteKnowledgeBase(kbID)
	if err != nil {
		return err
	}

	// Delete all documents
	docs, err := m.documentStore.ListDocuments(kbID)
	if err != nil {
		return err
	}

	for _, doc := range docs {
		err = m.documentStore.DeleteDocument(kbID, doc.ID)
		if err != nil {
			m.logger.Error("Failed to delete document",
				zap.String("doc_id", doc.ID),
				zap.Error(err),
			)
		}
	}

	m.logger.Info("Knowledge base deleted",
		zap.String("kb_id", kbID),
	)

	return nil
}

// GetStats returns statistics for a knowledge base
func (m *Manager) GetStats(kbID string) (docCount, chunkCount int, err error) {
	docs, err := m.documentStore.ListDocuments(kbID)
	if err != nil {
		return 0, 0, err
	}

	docCount = len(docs)

	chunkCount, err = m.vectorStore.GetStats(kbID)
	if err != nil {
		return 0, 0, err
	}

	return docCount, chunkCount, nil
}

// GetRelevantContext retrieves relevant chunks for a query (useful for RAG)
func (m *Manager) GetRelevantContext(kbID, query string, topK int) (string, error) {
	results, err := m.Search(kbID, query, topK, 0.7)
	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "", nil
	}

	// Combine top results into context
	context := ""
	for i, result := range results {
		if i > 0 {
			context += "\n\n---\n\n"
		}
		context += result.Chunk.Content
	}

	return context, nil
}
