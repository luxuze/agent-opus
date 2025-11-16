package knowledge

import (
	"strings"

	"github.com/google/uuid"
)

// Chunker splits documents into smaller chunks
type Chunker struct {
	config ChunkConfig
}

// NewChunker creates a new chunker with the given configuration
func NewChunker(config ChunkConfig) *Chunker {
	return &Chunker{config: config}
}

// ChunkDocument splits a document into chunks
func (c *Chunker) ChunkDocument(doc *Document) ([]*Chunk, error) {
	chunks := []*Chunk{}

	// Split by separator first
	parts := strings.Split(doc.Content, c.config.Separator)

	currentChunk := ""
	chunkIndex := 0

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// If adding this part would exceed chunk size, save current chunk
		if len(currentChunk)+len(part) > c.config.ChunkSize && currentChunk != "" {
			chunk := c.createChunk(doc.ID, currentChunk, chunkIndex)
			chunks = append(chunks, chunk)
			chunkIndex++

			// Start new chunk with overlap
			if c.config.ChunkOverlap > 0 && len(currentChunk) > c.config.ChunkOverlap {
				currentChunk = currentChunk[len(currentChunk)-c.config.ChunkOverlap:]
			} else {
				currentChunk = ""
			}
		}

		// Add part to current chunk
		if currentChunk != "" {
			currentChunk += c.config.Separator + part
		} else {
			currentChunk = part
		}

		// If current chunk is too large, split it further
		for len(currentChunk) > c.config.ChunkSize {
			splitPoint := c.config.ChunkSize
			chunk := c.createChunk(doc.ID, currentChunk[:splitPoint], chunkIndex)
			chunks = append(chunks, chunk)
			chunkIndex++

			if c.config.ChunkOverlap > 0 {
				currentChunk = currentChunk[splitPoint-c.config.ChunkOverlap:]
			} else {
				currentChunk = currentChunk[splitPoint:]
			}
		}
	}

	// Add remaining chunk
	if currentChunk != "" {
		chunk := c.createChunk(doc.ID, currentChunk, chunkIndex)
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}

func (c *Chunker) createChunk(documentID, content string, index int) *Chunk {
	return &Chunk{
		ID:         uuid.New().String(),
		DocumentID: documentID,
		Content:    strings.TrimSpace(content),
		Index:      index,
		Metadata: map[string]interface{}{
			"chunk_size": len(content),
		},
	}
}

// ChunkText is a helper function to chunk text directly
func ChunkText(text string, chunkSize, overlap int) []string {
	config := ChunkConfig{
		ChunkSize:    chunkSize,
		ChunkOverlap: overlap,
		Separator:    "\n\n",
	}

	chunker := NewChunker(config)
	doc := &Document{
		ID:      "temp",
		Content: text,
	}

	chunks, _ := chunker.ChunkDocument(doc)
	result := make([]string, len(chunks))
	for i, chunk := range chunks {
		result[i] = chunk.Content
	}

	return result
}
