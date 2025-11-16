package knowledge

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"go.uber.org/zap"
)

// VectorStore stores and retrieves vectors
type VectorStore interface {
	AddChunks(kbID string, chunks []*Chunk) error
	Search(kbID string, queryEmbedding []float32, topK int, threshold float64) ([]*SearchResult, error)
	DeleteKnowledgeBase(kbID string) error
	GetStats(kbID string) (int, error)
}

// InMemoryVectorStore is a simple in-memory vector store
type InMemoryVectorStore struct {
	mu     sync.RWMutex
	chunks map[string][]*Chunk // kbID -> chunks
	logger *zap.Logger
}

// NewInMemoryVectorStore creates a new in-memory vector store
func NewInMemoryVectorStore(logger *zap.Logger) *InMemoryVectorStore {
	return &InMemoryVectorStore{
		chunks: make(map[string][]*Chunk),
		logger: logger,
	}
}

// AddChunks adds chunks to the store
func (s *InMemoryVectorStore) AddChunks(kbID string, chunks []*Chunk) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.chunks[kbID] == nil {
		s.chunks[kbID] = []*Chunk{}
	}

	s.chunks[kbID] = append(s.chunks[kbID], chunks...)

	s.logger.Info("Added chunks to vector store",
		zap.String("kb_id", kbID),
		zap.Int("count", len(chunks)),
		zap.Int("total", len(s.chunks[kbID])),
	)

	return nil
}

// Search performs vector similarity search
func (s *InMemoryVectorStore) Search(kbID string, queryEmbedding []float32, topK int, threshold float64) ([]*SearchResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	chunks, exists := s.chunks[kbID]
	if !exists {
		return []*SearchResult{}, nil
	}

	// Calculate similarities
	results := make([]*SearchResult, 0)
	for _, chunk := range chunks {
		if chunk.Embedding == nil {
			continue
		}

		similarity := CosineSimilarity(queryEmbedding, chunk.Embedding)
		if similarity >= threshold {
			results = append(results, &SearchResult{
				Chunk:      chunk,
				Score:      similarity,
				DocumentID: chunk.DocumentID,
			})
		}
	}

	// Sort by score descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// Return top K
	if len(results) > topK {
		results = results[:topK]
	}

	s.logger.Info("Vector search completed",
		zap.String("kb_id", kbID),
		zap.Int("results", len(results)),
		zap.Int("top_k", topK),
	)

	return results, nil
}

// DeleteKnowledgeBase removes all chunks for a knowledge base
func (s *InMemoryVectorStore) DeleteKnowledgeBase(kbID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.chunks, kbID)

	s.logger.Info("Deleted knowledge base from vector store",
		zap.String("kb_id", kbID),
	)

	return nil
}

// GetStats returns statistics for a knowledge base
func (s *InMemoryVectorStore) GetStats(kbID string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	chunks, exists := s.chunks[kbID]
	if !exists {
		return 0, nil
	}

	return len(chunks), nil
}

// SerializeChunks converts chunks to JSON for database storage
func SerializeChunks(chunks []*Chunk) ([]byte, error) {
	return json.Marshal(chunks)
}

// DeserializeChunks converts JSON back to chunks
func DeserializeChunks(data []byte) ([]*Chunk, error) {
	var chunks []*Chunk
	err := json.Unmarshal(data, &chunks)
	return chunks, err
}

// DocumentStore manages documents and their chunks
type DocumentStore struct {
	documents map[string]*Document // docID -> document
	chunks    map[string][]*Chunk  // docID -> chunks
	mu        sync.RWMutex
	logger    *zap.Logger
}

// NewDocumentStore creates a new document store
func NewDocumentStore(logger *zap.Logger) *DocumentStore {
	return &DocumentStore{
		documents: make(map[string]*Document),
		chunks:    make(map[string][]*Chunk),
		logger:    logger,
	}
}

// AddDocument adds a document and its chunks
func (ds *DocumentStore) AddDocument(kbID string, doc *Document, chunks []*Chunk) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	docKey := fmt.Sprintf("%s:%s", kbID, doc.ID)
	ds.documents[docKey] = doc
	ds.chunks[docKey] = chunks

	ds.logger.Info("Added document to store",
		zap.String("kb_id", kbID),
		zap.String("doc_id", doc.ID),
		zap.Int("chunks", len(chunks)),
	)

	return nil
}

// GetDocument retrieves a document
func (ds *DocumentStore) GetDocument(kbID, docID string) (*Document, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	docKey := fmt.Sprintf("%s:%s", kbID, docID)
	doc, exists := ds.documents[docKey]
	if !exists {
		return nil, fmt.Errorf("document not found: %s", docID)
	}

	return doc, nil
}

// GetChunks retrieves chunks for a document
func (ds *DocumentStore) GetChunks(kbID, docID string) ([]*Chunk, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	docKey := fmt.Sprintf("%s:%s", kbID, docID)
	chunks, exists := ds.chunks[docKey]
	if !exists {
		return []*Chunk{}, nil
	}

	return chunks, nil
}

// ListDocuments lists all documents in a knowledge base
func (ds *DocumentStore) ListDocuments(kbID string) ([]*Document, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	docs := make([]*Document, 0)
	prefix := kbID + ":"

	for key, doc := range ds.documents {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			docs = append(docs, doc)
		}
	}

	return docs, nil
}

// DeleteDocument removes a document and its chunks
func (ds *DocumentStore) DeleteDocument(kbID, docID string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	docKey := fmt.Sprintf("%s:%s", kbID, docID)
	delete(ds.documents, docKey)
	delete(ds.chunks, docKey)

	ds.logger.Info("Deleted document from store",
		zap.String("kb_id", kbID),
		zap.String("doc_id", docID),
	)

	return nil
}
