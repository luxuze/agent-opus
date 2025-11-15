package main

import (
	"agent-platform/internal/config"
	grpcserver "agent-platform/internal/grpc"
	"fmt"
	"log"
	"net"

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

	// Create listener for gRPC
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register services
	pb.RegisterAgentServiceServer(grpcServer, grpcserver.NewAgentServer())
	pb.RegisterConversationServiceServer(grpcServer, grpcserver.NewConversationServer())
	pb.RegisterToolServiceServer(grpcServer, grpcserver.NewToolServer())
	pb.RegisterKnowledgeBaseServiceServer(grpcServer, grpcserver.NewKnowledgeBaseServer())

	// Register reflection service (for grpcurl and other tools)
	reflection.Register(grpcServer)

	// Start HTTP Gateway (REST API)
	if err := setupGateway(addr, cfg.Server.HTTPPort, logger); err != nil {
		logger.Fatal("Failed to setup HTTP Gateway", zap.Error(err))
	}

	logger.Info("gRPC Server listening", zap.String("address", addr))
	logger.Info("HTTP Gateway listening", zap.String("port", cfg.Server.HTTPPort))

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
