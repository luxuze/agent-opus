package ai

import (
	"context"
	"fmt"
	"io"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

// DeepSeekService implements AIService using SiliconFlow's DeepSeek API
// SiliconFlow API is compatible with OpenAI API format
type DeepSeekService struct {
	client       *openai.Client
	logger       *zap.Logger
	defaultModel string
}

// NewDeepSeekService creates a new DeepSeek service via SiliconFlow
func NewDeepSeekService(apiKey, apiBase, defaultModel string, logger *zap.Logger) *DeepSeekService {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = apiBase

	client := openai.NewClientWithConfig(config)
	return &DeepSeekService{
		client:       client,
		logger:       logger,
		defaultModel: defaultModel,
	}
}

// Chat sends a chat request to DeepSeek and returns the response
func (s *DeepSeekService) Chat(request ChatRequest) (*ChatResponse, error) {
	// Convert messages to OpenAI format
	messages := make([]openai.ChatCompletionMessage, len(request.Messages))
	for i, msg := range request.Messages {
		messages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Set default model if not specified
	model := request.Model
	if model == "" {
		model = s.defaultModel
	}

	// Set default temperature if not specified
	temperature := request.Temperature
	if temperature == 0 {
		temperature = 0.7
	}

	// Create chat completion request
	req := openai.ChatCompletionRequest{
		Model:       model,
		Messages:    messages,
		Temperature: temperature,
	}

	if request.MaxTokens > 0 {
		req.MaxTokens = request.MaxTokens
	}

	s.logger.Info("Sending chat request to DeepSeek via SiliconFlow",
		zap.String("model", model),
		zap.Int("messages", len(messages)),
	)

	// Send request
	resp, err := s.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		s.logger.Error("Failed to create chat completion", zap.Any("dsclient", s), zap.Error(err))
		return nil, fmt.Errorf("failed to create chat completion: %w", err)
	}

	// Extract response
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from DeepSeek")
	}

	choice := resp.Choices[0]
	return &ChatResponse{
		Content:      choice.Message.Content,
		FinishReason: string(choice.FinishReason),
		TokensUsed:   resp.Usage.TotalTokens,
		Model:        resp.Model,
	}, nil
}

// ChatStream sends a streaming chat request to DeepSeek
func (s *DeepSeekService) ChatStream(request ChatRequest) (<-chan string, <-chan error) {
	contentChan := make(chan string)
	errChan := make(chan error, 1)

	go func() {
		defer close(contentChan)
		defer close(errChan)

		// Convert messages to OpenAI format
		messages := make([]openai.ChatCompletionMessage, len(request.Messages))
		for i, msg := range request.Messages {
			messages[i] = openai.ChatCompletionMessage{
				Role:    msg.Role,
				Content: msg.Content,
			}
		}

		// Set default model if not specified
		model := request.Model
		if model == "" {
			model = s.defaultModel
		}

		// Set default temperature if not specified
		temperature := request.Temperature
		if temperature == 0 {
			temperature = 0.7
		}

		// Create streaming chat completion request
		req := openai.ChatCompletionRequest{
			Model:       model,
			Messages:    messages,
			Temperature: temperature,
			Stream:      true,
		}

		if request.MaxTokens > 0 {
			req.MaxTokens = request.MaxTokens
		}

		s.logger.Info("Sending streaming chat request to DeepSeek via SiliconFlow",
			zap.String("model", model),
			zap.Int("messages", len(messages)),
		)

		// Send streaming request
		stream, err := s.client.CreateChatCompletionStream(context.Background(), req)
		if err != nil {
			s.logger.Error("Failed to create chat completion stream", zap.Error(err))
			errChan <- fmt.Errorf("failed to create chat completion stream: %w", err)
			return
		}
		defer stream.Close()

		// Read stream
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				s.logger.Error("Error receiving stream", zap.Error(err))
				errChan <- fmt.Errorf("error receiving stream: %w", err)
				return
			}

			if len(response.Choices) > 0 {
				delta := response.Choices[0].Delta.Content
				if delta != "" {
					contentChan <- delta
				}
			}
		}
	}()

	return contentChan, errChan
}
