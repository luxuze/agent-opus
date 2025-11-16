package ai

import (
	"agent-platform/internal/config"
	"fmt"

	"go.uber.org/zap"
)

// Manager manages different AI service providers
type Manager struct {
	openai   AIService
	deepseek AIService
	logger   *zap.Logger
}

// NewManager creates a new AI manager
func NewManager(cfg *config.Config, logger *zap.Logger) (*Manager, error) {
	var openaiService AIService
	var deepseekService AIService

	// Initialize OpenAI service if API key is provided
	if cfg.OpenAI.APIKey != "" && cfg.OpenAI.APIKey != "your-openai-api-key" {
		openaiService = NewOpenAIService(cfg.OpenAI.APIKey, logger)
		logger.Info("OpenAI service initialized")
	} else {
		logger.Warn("OpenAI API key not configured")
	}

	// Initialize DeepSeek service (via SiliconFlow) if API key is provided
	if cfg.SiliconFlow.APIKey != "" && cfg.SiliconFlow.APIKey != "your-siliconflow-api-key" {
		deepseekService = NewDeepSeekService(
			cfg.SiliconFlow.APIKey,
			cfg.SiliconFlow.APIBase,
			cfg.SiliconFlow.Model,
			logger,
		)
		logger.Info("DeepSeek service (via SiliconFlow) initialized",
			zap.String("model", cfg.SiliconFlow.Model),
		)
	} else {
		logger.Warn("SiliconFlow API key not configured")
	}

	// At least one service should be configured
	if openaiService == nil && deepseekService == nil {
		logger.Warn("No AI services configured, AI features will be limited")
	}

	return &Manager{
		openai:   openaiService,
		deepseek: deepseekService,
		logger:   logger,
	}, nil
}

// GetService returns the appropriate AI service based on model name
func (m *Manager) GetService(model string) (AIService, error) {
	// Route to DeepSeek service if model name contains "deepseek"
	if len(model) >= 8 && (model[:8] == "deepseek" || model == "DeepSeek-V3" || model[:11] == "deepseek-ai") {
		if m.deepseek == nil {
			return nil, fmt.Errorf("DeepSeek service not initialized - please configure SILICONFLOW_API_KEY")
		}
		return m.deepseek, nil
	}

	// Default to OpenAI service
	if m.openai == nil {
		// If OpenAI is not available, try DeepSeek as fallback
		if m.deepseek != nil {
			m.logger.Info("OpenAI not configured, using DeepSeek as fallback")
			return m.deepseek, nil
		}
		return nil, fmt.Errorf("no AI service initialized - please configure OPENAI_API_KEY or SILICONFLOW_API_KEY")
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
