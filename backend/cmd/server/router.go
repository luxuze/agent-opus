package main

import (
	"agent-platform/internal/config"
	"agent-platform/internal/handler"
	"agent-platform/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func setupRouter(cfg *config.Config, logger *zap.Logger) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware(logger))
	r.Use(middleware.CORSMiddleware(cfg.CORS.Origins))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Public routes
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			// Agent routes
			agentHandler := handler.NewAgentHandler()
			agents := protected.Group("/agents")
			{
				agents.POST("", agentHandler.CreateAgent)
				agents.GET("", agentHandler.ListAgents)
				agents.GET("/:id", agentHandler.GetAgent)
				agents.PUT("/:id", agentHandler.UpdateAgent)
				agents.DELETE("/:id", agentHandler.DeleteAgent)
			}

			// Conversation routes
			conversationHandler := handler.NewConversationHandler()
			conversations := protected.Group("/conversations")
			{
				conversations.POST("", conversationHandler.CreateConversation)
				conversations.GET("", conversationHandler.ListConversations)
				conversations.GET("/:id", conversationHandler.GetConversation)
				conversations.POST("/:id/messages", conversationHandler.SendMessage)
			}

			// Tool routes
			toolHandler := handler.NewToolHandler()
			tools := protected.Group("/tools")
			{
				tools.POST("", toolHandler.CreateTool)
				tools.GET("", toolHandler.ListTools)
				tools.GET("/:id", toolHandler.GetTool)
				tools.DELETE("/:id", toolHandler.DeleteTool)
			}

			// Knowledge base routes
			kbHandler := handler.NewKnowledgeBaseHandler()
			knowledgeBases := protected.Group("/knowledge-bases")
			{
				knowledgeBases.POST("", kbHandler.CreateKnowledgeBase)
				knowledgeBases.GET("", kbHandler.ListKnowledgeBases)
				knowledgeBases.GET("/:id", kbHandler.GetKnowledgeBase)
				knowledgeBases.POST("/:id/documents", kbHandler.UploadDocument)
				knowledgeBases.DELETE("/:id", kbHandler.DeleteKnowledgeBase)
			}
		}
	}

	return r
}
