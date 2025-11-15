package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server    ServerConfig
	MySQL     MySQLConfig
	Redis     RedisConfig
	Milvus    MilvusConfig
	JWT       JWTConfig
	OpenAI    OpenAIConfig
	Embedding EmbeddingConfig
	CORS      CORSConfig
	Log       LogConfig
}

type ServerConfig struct {
	Port string
	Mode string
	Host string
}

type MySQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type MilvusConfig struct {
	Host string
	Port string
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

type OpenAIConfig struct {
	APIKey  string
	APIBase string
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
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("SERVER_MODE", "debug"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		MySQL: MySQLConfig{
			Host:     getEnv("MYSQL_HOST", "localhost"),
			Port:     getEnv("MYSQL_PORT", "3306"),
			User:     getEnv("MYSQL_USER", "root"),
			Password: getEnv("MYSQL_PASSWORD", ""),
			Database: getEnv("MYSQL_DATABASE", "agent_platform"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		Milvus: MilvusConfig{
			Host: getEnv("MILVUS_HOST", "localhost"),
			Port: getEnv("MILVUS_PORT", "19530"),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "your-secret-key"),
			ExpireHours: jwtExpireHours,
		},
		OpenAI: OpenAIConfig{
			APIKey:  getEnv("OPENAI_API_KEY", ""),
			APIBase: getEnv("OPENAI_API_BASE", "https://api.openai.com/v1"),
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

func (c *MySQLConfig) DSN() string {
	return c.User + ":" + c.Password + "@tcp(" + c.Host + ":" + c.Port + ")/" + c.Database + "?parseTime=true&charset=utf8mb4"
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
