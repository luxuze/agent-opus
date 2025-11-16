package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server      ServerConfig
	Postgres    PostgresConfig
	Redis       RedisConfig
	JWT         JWTConfig
	OpenAI      OpenAIConfig
	SiliconFlow SiliconFlowConfig
	Embedding   EmbeddingConfig
	CORS        CORSConfig
	Log         LogConfig
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

type OpenAIConfig struct {
	APIKey  string
	APIBase string
}

type SiliconFlowConfig struct {
	APIKey  string
	APIBase string
	Model   string
}

type EmbeddingConfig struct {
	Model     string
	Dimension int
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
		OpenAI: OpenAIConfig{
			APIKey:  getEnv("OPENAI_API_KEY", ""),
			APIBase: getEnv("OPENAI_API_BASE", "https://api.openai.com/v1"),
		},
		SiliconFlow: SiliconFlowConfig{
			APIKey:  getEnv("SILICONFLOW_API_KEY", ""),
			APIBase: getEnv("SILICONFLOW_API_BASE", "https://api.siliconflow.cn/v1"),
			Model:   getEnv("SILICONFLOW_MODEL", "deepseek-ai/DeepSeek-V3"),
		},
		Embedding: EmbeddingConfig{
			Model:     getEnv("EMBEDDING_MODEL", "text-embedding-ada-002"),
			Dimension: embeddingDim,
		},
		CORS: CORSConfig{
			Origins: strings.Split(getEnv("CORS_ORIGINS", "http://localhost:5173"), ","),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
			File:  getEnv("LOG_FILE", "logs/app.log"),
		},
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
