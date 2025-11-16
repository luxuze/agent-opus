package knowledge

import "time"

// Document represents a document in the knowledge base
type Document struct {
	ID           string                 `json:"id"`
	Title        string                 `json:"title"`
	Content      string                 `json:"content"`
	ContentType  string                 `json:"content_type"` // text, markdown, pdf, etc.
	Source       string                 `json:"source"`       // file path, URL, etc.
	Metadata     map[string]interface{} `json:"metadata"`
	UploadedAt   time.Time              `json:"uploaded_at"`
	ChunkCount   int                    `json:"chunk_count"`
}

// Chunk represents a text chunk from a document
type Chunk struct {
	ID         string                 `json:"id"`
	DocumentID string                 `json:"document_id"`
	Content    string                 `json:"content"`
	Index      int                    `json:"index"`
	Metadata   map[string]interface{} `json:"metadata"`
	Embedding  []float32              `json:"embedding,omitempty"`
}

// SearchRequest represents a vector search request
type SearchRequest struct {
	Query      string
	TopK       int
	Threshold  float64
	Metadata   map[string]interface{}
}

// SearchResult represents a search result with similarity score
type SearchResult struct {
	Chunk      *Chunk
	Score      float64
	DocumentID string
}

// ChunkConfig defines how to chunk documents
type ChunkConfig struct {
	ChunkSize    int // Characters per chunk
	ChunkOverlap int // Overlap between chunks
	Separator    string
}

// DefaultChunkConfig returns default chunking configuration
func DefaultChunkConfig() ChunkConfig {
	return ChunkConfig{
		ChunkSize:    1000,
		ChunkOverlap: 200,
		Separator:    "\n\n",
	}
}
