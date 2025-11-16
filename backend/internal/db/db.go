package db

import (
	"agent-platform/internal/config"
	"agent-platform/internal/model/ent"
	"context"
	"fmt"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// Client wraps the ent client
type Client struct {
	*ent.Client
	logger *zap.Logger
}

// NewClient creates a new database client
func NewClient(cfg *config.Config, logger *zap.Logger) (*Client, error) {
	dsn := cfg.Postgres.DSN()

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to postgres: %w", err)
	}

	// Connection pool settings are configured via DSN parameters or driver config

	logger.Info("Database connected successfully",
		zap.String("host", cfg.Postgres.Host),
		zap.String("port", cfg.Postgres.Port),
		zap.String("database", cfg.Postgres.Database),
	)

	return &Client{
		Client: client,
		logger: logger,
	}, nil
}

// AutoMigrate runs database migrations
func (c *Client) AutoMigrate(ctx context.Context) error {
	c.logger.Info("Running database migrations...")

	if err := c.Client.Schema.Create(ctx); err != nil {
		return fmt.Errorf("failed creating schema resources: %w", err)
	}

	c.logger.Info("Database migrations completed successfully")
	return nil
}

// Close closes the database connection
func (c *Client) Close() error {
	c.logger.Info("Closing database connection...")
	return c.Client.Close()
}

// Health checks the database connection
func (c *Client) Health(ctx context.Context) error {
	// Simple query to check connection
	_, err := c.Client.Agent.Query().Count(ctx)
	return err
}
