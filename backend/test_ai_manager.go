package main

import (
	"agent-platform/internal/ai"
	"agent-platform/internal/config"
	"fmt"
	"log"

	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	// Initialize logger
	logger, err := initTestLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	fmt.Println("=== AI Manager Test ===\n")

	// Initialize AI Manager
	manager, err := ai.NewManager(cfg, logger)
	if err != nil {
		log.Fatalf("Failed to create AI manager: %v", err)
	}

	// Test model routing
	testCases := []string{
		"gpt-4o",
		"gpt-3.5-turbo",
		"deepseek-ai/DeepSeek-V3",
		"deepseek-chat",
		"claude-3-5-sonnet-20241022",
		"claude-3-haiku-20240307",
		"", // Should use default provider
	}

	fmt.Println("Testing Model Routing:")
	for _, model := range testCases {
		displayModel := model
		if displayModel == "" {
			displayModel = "(empty - use default)"
		}

		service, err := manager.GetService(model)
		if err != nil {
			fmt.Printf("  ✗ Model: %-35s -> Error: %v\n", displayModel, err)
		} else {
			fmt.Printf("  ✓ Model: %-35s -> Service Found\n", displayModel)
			_ = service // Use the service to avoid unused variable warning
		}
	}

	fmt.Println("\n=== AI Manager Test Completed ===")
}

func initTestLogger() (*zap.Logger, error) {
	return zap.NewDevelopment()
}
