package ai

import (
	"agent-platform/internal/config"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

// Manager manages different AI service providers
type Manager struct {
	services map[string]AIService // Map of provider name to service
	config   *config.AIConfig
	logger   *zap.Logger
}

// NewManager creates a new AI manager with unified configuration
func NewManager(cfg *config.Config, logger *zap.Logger) (*Manager, error) {
	services := make(map[string]AIService)

	// Initialize services based on configuration
	for name, providerCfg := range cfg.AI.Providers {
		if !providerCfg.Enabled {
			continue
		}

		var service AIService

		switch name {
		case "openai":
			service = NewOpenAIService(providerCfg.APIKey, logger)
			logger.Info("OpenAI service initialized",
				zap.String("default_model", providerCfg.DefaultModel),
			)

		case "siliconflow":
			service = NewDeepSeekService(
				providerCfg.APIKey,
				providerCfg.APIBase,
				providerCfg.DefaultModel,
				logger,
			)
			logger.Info("SiliconFlow (DeepSeek) service initialized",
				zap.String("default_model", providerCfg.DefaultModel),
			)

		case "anthropic":
			// TODO: Implement Anthropic service
			logger.Warn("Anthropic provider configured but not yet implemented",
				zap.String("provider", name),
			)
			continue

		default:
			logger.Warn("Unknown AI provider",
				zap.String("provider", name),
			)
			continue
		}

		if service != nil {
			services[name] = service
		}
	}

	// Check if at least one service is configured
	if len(services) == 0 {
		logger.Warn("No AI services configured, AI features will be limited")
	}

	return &Manager{
		services: services,
		config:   &cfg.AI,
		logger:   logger,
	}, nil
}

// GetService returns the appropriate AI service based on model name
func (m *Manager) GetService(model string) (AIService, error) {
	if model == "" {
		// Use default provider
		if m.config.DefaultProvider != "" {
			if service, ok := m.services[m.config.DefaultProvider]; ok {
				return service, nil
			}
		}

		// Fallback: return first available service
		for _, service := range m.services {
			return service, nil
		}

		return nil, fmt.Errorf("no AI service available")
	}

	// Route based on model name patterns
	modelLower := strings.ToLower(model)

	// Check DeepSeek models
	if strings.Contains(modelLower, "deepseek") {
		if service, ok := m.services["siliconflow"]; ok {
			return service, nil
		}
		return nil, fmt.Errorf("DeepSeek model requested but SiliconFlow provider not configured")
	}

	// Check Claude models
	if strings.Contains(modelLower, "claude") {
		if service, ok := m.services["anthropic"]; ok {
			return service, nil
		}
		return nil, fmt.Errorf("Claude model requested but Anthropic provider not configured")
	}

	// Check GPT models (OpenAI)
	if strings.Contains(modelLower, "gpt") || strings.Contains(modelLower, "text-") {
		if service, ok := m.services["openai"]; ok {
			return service, nil
		}
		return nil, fmt.Errorf("GPT model requested but OpenAI provider not configured")
	}

	// Try to match model name with provider's supported models
	for providerName, providerCfg := range m.config.Providers {
		for _, supportedModel := range providerCfg.Models {
			if supportedModel == model {
				if service, ok := m.services[providerName]; ok {
					return service, nil
				}
			}
		}
	}

	// Use default provider as fallback
	if m.config.DefaultProvider != "" {
		if service, ok := m.services[m.config.DefaultProvider]; ok {
			m.logger.Info("Model not matched, using default provider",
				zap.String("model", model),
				zap.String("provider", m.config.DefaultProvider),
			)
			return service, nil
		}
	}

	// Last resort: return first available service
	for name, service := range m.services {
		m.logger.Warn("Model not matched, using first available provider",
			zap.String("model", model),
			zap.String("provider", name),
		)
		return service, nil
	}

	return nil, fmt.Errorf("no AI service available for model: %s", model)
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
