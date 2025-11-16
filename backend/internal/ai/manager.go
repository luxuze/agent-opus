package ai

import (
	"agent-platform/internal/config"
	"fmt"

	"go.uber.org/zap"
)

// Manager manages different AI service providers
type Manager struct {
	openai AIService
	logger *zap.Logger
}

// NewManager creates a new AI manager
func NewManager(cfg *config.Config, logger *zap.Logger) (*Manager, error) {
	var openaiService AIService

	// Initialize OpenAI service if API key is provided
	if cfg.OpenAI.APIKey != "" && cfg.OpenAI.APIKey != "your-openai-api-key" {
		openaiService = NewOpenAIService(cfg.OpenAI.APIKey, logger)
		logger.Info("OpenAI service initialized")
	} else {
		logger.Warn("OpenAI API key not configured, AI features will be limited")
	}

	return &Manager{
		openai: openaiService,
		logger: logger,
	}, nil
}

// GetService returns the appropriate AI service based on model name
func (m *Manager) GetService(model string) (AIService, error) {
	// For now, we only support OpenAI
	// In the future, we can add logic to route to different providers based on model name
	if m.openai == nil {
		return nil, fmt.Errorf("OpenAI service not initialized - please configure OPENAI_API_KEY")
	}

	return m.openai, nil
}

// Chat is a convenience method that routes to the appropriate service
func (m *Manager) Chat(request ChatRequest) (*ChatResponse, error) {
	service, err := m.GetService(request.Model)
	if err != nil {
		return nil, err
	}

	return service.Chat(request)
}

// ChatStream is a convenience method that routes to the appropriate service
func (m *Manager) ChatStream(request ChatRequest) (<-chan string, <-chan error) {
	service, err := m.GetService(request.Model)
	if err != nil {
		errChan := make(chan error, 1)
		contentChan := make(chan string)
		errChan <- err
		close(errChan)
		close(contentChan)
		return contentChan, errChan
	}

	return service.ChatStream(request)
}
