package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env.test file
	if err := godotenv.Load(".env.test"); err != nil {
		log.Printf("Warning: %v", err)
	}

	// Check environment variables
	fmt.Println("=== Environment Variables ===")
	fmt.Printf("OPENAI_API_KEY: %s\n", os.Getenv("OPENAI_API_KEY"))
	fmt.Printf("SILICONFLOW_API_KEY: %s\n", os.Getenv("SILICONFLOW_API_KEY"))
	fmt.Printf("ANTHROPIC_API_KEY: %s\n", os.Getenv("ANTHROPIC_API_KEY"))
	fmt.Printf("AI_DEFAULT_PROVIDER: %s\n", os.Getenv("AI_DEFAULT_PROVIDER"))
}
