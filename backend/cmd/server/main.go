package main

import (
	"agent-platform/internal/ai"
	"agent-platform/internal/auth"
	"agent-platform/internal/config"
	"agent-platform/internal/db"
	grpcserver "agent-platform/internal/grpc"
	"agent-platform/internal/knowledge"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "agent-platform/gen/go"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger, err := initLogger(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting Agent Platform Server",
		zap.String("grpc_port", cfg.Server.GRPCPort),
		zap.String("http_port", cfg.Server.HTTPPort),
		zap.String("mode", cfg.Server.Mode),
	)

	// Initialize database
	ctx := context.Background()
	dbClient, err := db.NewClient(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer dbClient.Close()

	// Run database migrations
	if err := dbClient.AutoMigrate(ctx); err != nil {
		logger.Fatal("Failed to run database migrations", zap.Error(err))
	}

	// Initialize AI manager
	aiManager, err := ai.NewManager(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize AI manager", zap.Error(err))
	}

	// Initialize Knowledge Base manager with pgvector
	// Get embedding provider configuration
	embeddingProvider := cfg.AI.EmbeddingProvider
	if embeddingProvider == "" {
		embeddingProvider = "openai" // default
	}

	var embeddingAPIKey, embeddingAPIBase string
	if providerCfg, ok := cfg.AI.Providers[embeddingProvider]; ok {
		embeddingAPIKey = providerCfg.APIKey
		embeddingAPIBase = providerCfg.APIBase
	} else {
		logger.Warn("Embedding provider not configured, using first available AI provider",
			zap.String("requested_provider", embeddingProvider),
		)
		// Fallback to first available provider
		for _, providerCfg := range cfg.AI.Providers {
			embeddingAPIKey = providerCfg.APIKey
			embeddingAPIBase = providerCfg.APIBase
			break
		}
	}

	kbManager, err := knowledge.NewManager(
		dbClient.Client,
		cfg.Postgres.DSN(),
		embeddingAPIKey,
		embeddingAPIBase,
		cfg.AI.EmbeddingModel,
		logger,
	)
	if err != nil {
		logger.Fatal("Failed to initialize KB manager", zap.Error(err))
	}

	// Initialize JWT service
	jwtService := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpireHours)

	// Create listener for gRPC
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	// Create gRPC server with auth interceptors
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			auth.UnaryAuthInterceptor(jwtService),
		),
		grpc.ChainStreamInterceptor(
			auth.StreamAuthInterceptor(jwtService),
		),
	)

	// Register services with database client, AI manager, and KB manager
	kbServer := grpcserver.NewKnowledgeBaseServer(dbClient.Client, kbManager)
	pb.RegisterAgentServiceServer(grpcServer, grpcserver.NewAgentServer(dbClient.Client))
	pb.RegisterConversationServiceServer(grpcServer, grpcserver.NewConversationServer(dbClient.Client, aiManager, kbServer))
	pb.RegisterToolServiceServer(grpcServer, grpcserver.NewToolServer(dbClient.Client))
	pb.RegisterKnowledgeBaseServiceServer(grpcServer, kbServer)
	pb.RegisterUserServiceServer(grpcServer, grpcserver.NewUserServer(dbClient.Client, jwtService))

	// Register reflection service (for grpcurl and other tools)
	reflection.Register(grpcServer)

	// Start HTTP Gateway (REST API)
	if err := setupGateway(addr, cfg.Server.HTTPPort, logger); err != nil {
		logger.Fatal("Failed to setup HTTP Gateway", zap.Error(err))
	}

	logger.Info("gRPC Server listening", zap.String("address", addr))
	logger.Info("HTTP Gateway listening", zap.String("port", cfg.Server.HTTPPort))

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		logger.Info("Shutting down gracefully...")
		grpcServer.GracefulStop()
	}()

	// Start gRPC server (blocking)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("Failed to serve", zap.Error(err))
	}
}

func initLogger(cfg *config.Config) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	if cfg.Server.Mode == "release" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	return logger, err
}
