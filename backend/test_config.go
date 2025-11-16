package main

import (
	"agent-platform/internal/config"
	"fmt"
	"log"
	"os"
)

func main() {
	// Change to backend directory
	if err := os.Chdir("/Users/xuzelu/workspace/agent-opus/backend"); err != nil {
		log.Fatal(err)
	}

	// Load .env.test file
	os.Setenv("ENV_FILE", ".env.test")

	cfg := config.Load()

	fmt.Println("=== AI Configuration Test ===\n")

	// Print AI configuration
	fmt.Printf("Default Provider: %s\n", cfg.AI.DefaultProvider)
	fmt.Printf("Embedding Provider: %s\n", cfg.AI.EmbeddingProvider)
	fmt.Printf("Embedding Model: %s\n", cfg.AI.EmbeddingModel)
	fmt.Printf("Embedding Dimension: %d\n\n", cfg.AI.EmbeddingDimension)

	fmt.Println("Configured Providers:")
	for name, provider := range cfg.AI.Providers {
		fmt.Printf("\n  Provider: %s\n", name)
		fmt.Printf("    Enabled: %v\n", provider.Enabled)
		fmt.Printf("    API Base: %s\n", provider.APIBase)
		fmt.Printf("    Default Model: %s\n", provider.DefaultModel)
		fmt.Printf("    Supported Models: %v\n", provider.Models)
		// Don't print API key for security
		fmt.Printf("    API Key: %s...\n", maskAPIKey(provider.APIKey))
	}

	fmt.Println("\n=== Configuration Loaded Successfully ===")
}

func maskAPIKey(key string) string {
	if len(key) <= 10 {
		return "***"
	}
	return key[:7] + "***" + key[len(key)-3:]
}
