package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	JWT      JWTConfig
	AI       AIConfig
	CORS     CORSConfig
	Log      LogConfig
}

type ServerConfig struct {
	Port     string // Deprecated: Use GRPCPort instead
	GRPCPort string
	HTTPPort string
	Mode     string
	Host     string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

// AIProviderConfig represents configuration for a single AI provider
type AIProviderConfig struct {
	Name         string   // Provider name (e.g., "openai", "siliconflow", "anthropic")
	APIKey       string   // API key
	APIBase      string   // Base URL for API
	DefaultModel string   // Default model to use
	Models       []string // Supported models
	Enabled      bool     // Whether this provider is enabled
}

// AIConfig contains configuration for all AI providers
type AIConfig struct {
	Providers          map[string]*AIProviderConfig // Map of provider name to config
	DefaultProvider    string                       // Default provider to use
	EmbeddingModel     string                       // Model to use for embeddings
	EmbeddingDimension int                          // Dimension of embeddings
	EmbeddingProvider  string                       // Provider to use for embeddings
}

type CORSConfig struct {
	Origins []string
}

type LogConfig struct {
	Level string
	File  string
}

func Load() *Config {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	jwtExpireHours, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
	embeddingDim, _ := strconv.Atoi(getEnv("EMBEDDING_DIMENSION", "1536"))

	// Load AI configuration
	aiConfig := loadAIConfig(embeddingDim)

	return &Config{
		Server: ServerConfig{
			Port:     getEnv("SERVER_PORT", "9000"), // For backward compatibility
			GRPCPort: getEnv("GRPC_PORT", getEnv("SERVER_PORT", "9000")),
			HTTPPort: getEnv("HTTP_PORT", "8000"),
			Mode:     getEnv("SERVER_MODE", "debug"),
			Host:     getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", ""),
			Database: getEnv("POSTGRES_DATABASE", "agent_platform"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "your-secret-key"),
			ExpireHours: jwtExpireHours,
		},
		AI: aiConfig,
		CORS: CORSConfig{
			Origins: strings.Split(getEnv("CORS_ORIGINS", "http://localhost:5173"), ","),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
			File:  getEnv("LOG_FILE", "logs/app.log"),
		},
	}
}

// loadAIConfig loads AI provider configurations from environment variables
func loadAIConfig(embeddingDim int) AIConfig {
	providers := make(map[string]*AIProviderConfig)

	// Load OpenAI configuration
	openaiKey := getEnv("OPENAI_API_KEY", "")
	if openaiKey != "" && openaiKey != "your-openai-api-key" {
		providers["openai"] = &AIProviderConfig{
			Name:         "openai",
			APIKey:       openaiKey,
			APIBase:      getEnv("OPENAI_API_BASE", "https://api.openai.com/v1"),
			DefaultModel: getEnv("OPENAI_DEFAULT_MODEL", "gpt-4o"),
			Models:       strings.Split(getEnv("OPENAI_MODELS", "gpt-4o,gpt-4,gpt-3.5-turbo"), ","),
			Enabled:      true,
		}
		log.Printf("OpenAI provider configured with key: %s...", openaiKey[:min(10, len(openaiKey))])
	} else {
		log.Printf("OpenAI provider not configured (key: %q)", openaiKey)
	}

	// Load SiliconFlow (DeepSeek) configuration
	siliconflowKey := getEnv("SILICONFLOW_API_KEY", "")
	if siliconflowKey != "" && siliconflowKey != "your-siliconflow-api-key" {
		providers["siliconflow"] = &AIProviderConfig{
			Name:         "siliconflow",
			APIKey:       siliconflowKey,
			APIBase:      getEnv("SILICONFLOW_API_BASE", "https://api.siliconflow.cn/v1"),
			DefaultModel: getEnv("SILICONFLOW_DEFAULT_MODEL", "deepseek-ai/DeepSeek-V3"),
			Models:       strings.Split(getEnv("SILICONFLOW_MODELS", "deepseek-ai/DeepSeek-V3,deepseek-chat"), ","),
			Enabled:      true,
		}
	}

	// Load Anthropic configuration
	anthropicKey := getEnv("ANTHROPIC_API_KEY", "")
	if anthropicKey != "" && anthropicKey != "your-anthropic-api-key" {
		providers["anthropic"] = &AIProviderConfig{
			Name:         "anthropic",
			APIKey:       anthropicKey,
			APIBase:      getEnv("ANTHROPIC_API_BASE", "https://api.anthropic.com"),
			DefaultModel: getEnv("ANTHROPIC_DEFAULT_MODEL", "claude-3-5-sonnet-20241022"),
			Models:       strings.Split(getEnv("ANTHROPIC_MODELS", "claude-3-5-sonnet-20241022,claude-3-opus-20240229,claude-3-haiku-20240307"), ","),
			Enabled:      true,
		}
	}

	// Determine default provider
	defaultProvider := getEnv("AI_DEFAULT_PROVIDER", "")
	if defaultProvider == "" {
		// Auto-detect default provider based on what's configured
		if _, ok := providers["openai"]; ok {
			defaultProvider = "openai"
		} else if _, ok := providers["siliconflow"]; ok {
			defaultProvider = "siliconflow"
		} else if _, ok := providers["anthropic"]; ok {
			defaultProvider = "anthropic"
		}
	}

	return AIConfig{
		Providers:          providers,
		DefaultProvider:    defaultProvider,
		EmbeddingModel:     getEnv("EMBEDDING_MODEL", "text-embedding-ada-002"),
		EmbeddingDimension: embeddingDim,
		EmbeddingProvider:  getEnv("EMBEDDING_PROVIDER", "openai"),
	}
}

func (c *PostgresConfig) DSN() string {
	return "host=" + c.Host + " port=" + c.Port + " user=" + c.User + " password=" + c.Password + " dbname=" + c.Database + " sslmode=" + c.SSLMode
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
