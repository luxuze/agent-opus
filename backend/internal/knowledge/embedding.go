package knowledge

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

// EmbeddingService generates embeddings for text
type EmbeddingService struct {
	client *openai.Client
	model  openai.EmbeddingModel
	logger *zap.Logger
}

// NewEmbeddingService creates a new embedding service
func NewEmbeddingService(apiKey, apiBase, model string, logger *zap.Logger) *EmbeddingService {
	var client *openai.Client

	if apiBase != "" && apiBase != "https://api.openai.com/v1" {
		// Use custom base URL
		config := openai.DefaultConfig(apiKey)
		config.BaseURL = apiBase
		client = openai.NewClientWithConfig(config)
	} else {
		// Use default OpenAI client
		client = openai.NewClient(apiKey)
	}

	embModel := openai.EmbeddingModel(model)
	if model == "" {
		embModel = openai.AdaEmbeddingV2
	}

	return &EmbeddingService{
		client: client,
		model:  embModel,
		logger: logger,
	}
}

// GenerateEmbedding generates an embedding for a single text
func (s *EmbeddingService) GenerateEmbedding(text string) ([]float32, error) {
	s.logger.Debug("Generating embedding",
		zap.String("model", string(s.model)),
		zap.Int("text_length", len(text)),
	)

	req := openai.EmbeddingRequest{
		Input: []string{text},
		Model: s.model,
	}

	resp, err := s.client.CreateEmbeddings(context.Background(), req)
	if err != nil {
		s.logger.Error("Failed to generate embedding", zap.Error(err))
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	return resp.Data[0].Embedding, nil
}

// GenerateEmbeddings generates embeddings for multiple texts in batch
func (s *EmbeddingService) GenerateEmbeddings(texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return [][]float32{}, nil
	}

	s.logger.Debug("Generating embeddings batch",
		zap.String("model", string(s.model)),
		zap.Int("count", len(texts)),
	)

	req := openai.EmbeddingRequest{
		Input: texts,
		Model: s.model,
	}

	resp, err := s.client.CreateEmbeddings(context.Background(), req)
	if err != nil {
		s.logger.Error("Failed to generate embeddings", zap.Error(err))
		return nil, fmt.Errorf("failed to generate embeddings: %w", err)
	}

	embeddings := make([][]float32, len(resp.Data))
	for i, data := range resp.Data {
		embeddings[i] = data.Embedding
	}

	return embeddings, nil
}

// EmbedChunks generates embeddings for all chunks
func (s *EmbeddingService) EmbedChunks(chunks []*Chunk) error {
	if len(chunks) == 0 {
		return nil
	}

	// Extract texts
	texts := make([]string, len(chunks))
	for i, chunk := range chunks {
		texts[i] = chunk.Content
	}

	// Generate embeddings in batch
	embeddings, err := s.GenerateEmbeddings(texts)
	if err != nil {
		return err
	}

	// Assign embeddings to chunks
	for i, embedding := range embeddings {
		chunks[i].Embedding = embedding
	}

	return nil
}

// CosineSimilarity calculates cosine similarity between two vectors
func CosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (normA * normB)
}
