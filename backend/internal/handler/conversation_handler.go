package handler

import (
	"strconv"
	"time"

	pb "agent-platform/gen/go"
	"agent-platform/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ConversationHandler struct{}

func NewConversationHandler() *ConversationHandler {
	return &ConversationHandler{}
}

// CreateConversation creates a new conversation
func (h *ConversationHandler) CreateConversation(c *gin.Context) {
	var req pb.CreateConversationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		userID = "anonymous"
	}

	now := timestamppb.New(time.Now())
	conversation := &pb.Conversation{
		Id:            uuid.New().String(),
		AgentId:       req.AgentId,
		UserId:        userID,
		Title:         req.Title,
		Messages:      []*pb.Message{},
		Context:       req.Context,
		Status:        "active",
		CreatedAt:     now,
		UpdatedAt:     now,
		LastMessageAt: now,
	}

	conversationMap := map[string]interface{}{
		"id":              conversation.Id,
		"agent_id":        conversation.AgentId,
		"user_id":         conversation.UserId,
		"title":           conversation.Title,
		"messages":        []interface{}{},
		"status":          conversation.Status,
		"created_at":      conversation.CreatedAt.AsTime(),
		"updated_at":      conversation.UpdatedAt.AsTime(),
		"last_message_at": conversation.LastMessageAt.AsTime(),
	}

	if conversation.Context != nil {
		conversationMap["context"] = conversation.Context.AsMap()
	}

	response.Created(c, conversationMap)
}

// SendMessage sends a message in a conversation
func (h *ConversationHandler) SendMessage(c *gin.Context) {
	conversationID := c.Param("id")

	var req pb.SendMessageRequest
	req.ConversationId = conversationID

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	now := timestamppb.New(time.Now())

	// Create user message
	userMessage := &pb.Message{
		Id:        uuid.New().String(),
		Role:      "user",
		Content:   req.Content,
		Timestamp: now,
	}

	// Mock AI response
	assistantMessage := &pb.Message{
		Id:        uuid.New().String(),
		Role:      "assistant",
		Content:   "This is a mock response from the agent. In production, this would call the AI model.",
		Timestamp: timestamppb.New(time.Now().Add(1 * time.Second)),
	}

	messages := []*pb.Message{userMessage, assistantMessage}

	// Convert to map
	messageMaps := make([]map[string]interface{}, len(messages))
	for i, msg := range messages {
		messageMaps[i] = map[string]interface{}{
			"id":        msg.Id,
			"role":      msg.Role,
			"content":   msg.Content,
			"timestamp": msg.Timestamp.AsTime(),
		}
		if msg.Metadata != nil {
			messageMaps[i]["metadata"] = msg.Metadata.AsMap()
		}
	}

	responseData := map[string]interface{}{
		"conversation_id": conversationID,
		"messages":        messageMaps,
	}

	response.Success(c, responseData)
}

// GetConversation retrieves a conversation
func (h *ConversationHandler) GetConversation(c *gin.Context) {
	id := c.Param("id")

	now := timestamppb.New(time.Now())
	conversation := &pb.Conversation{
		Id:      id,
		AgentId: "agent-001",
		UserId:  "user-001",
		Title:   "Sample Conversation",
		Messages: []*pb.Message{
			{
				Id:        "msg-001",
				Role:      "user",
				Content:   "Hello!",
				Timestamp: timestamppb.New(time.Now().Add(-10 * time.Minute)),
			},
			{
				Id:        "msg-002",
				Role:      "assistant",
				Content:   "Hi! How can I help you today?",
				Timestamp: timestamppb.New(time.Now().Add(-9 * time.Minute)),
			},
		},
		Status:        "active",
		CreatedAt:     timestamppb.New(time.Now().Add(-1 * time.Hour)),
		UpdatedAt:     now,
		LastMessageAt: timestamppb.New(time.Now().Add(-9 * time.Minute)),
	}

	// Convert messages to map
	messageMaps := make([]map[string]interface{}, len(conversation.Messages))
	for i, msg := range conversation.Messages {
		messageMaps[i] = map[string]interface{}{
			"id":        msg.Id,
			"role":      msg.Role,
			"content":   msg.Content,
			"timestamp": msg.Timestamp.AsTime(),
		}
	}

	conversationMap := map[string]interface{}{
		"id":              conversation.Id,
		"agent_id":        conversation.AgentId,
		"user_id":         conversation.UserId,
		"title":           conversation.Title,
		"messages":        messageMaps,
		"status":          conversation.Status,
		"created_at":      conversation.CreatedAt.AsTime(),
		"updated_at":      conversation.UpdatedAt.AsTime(),
		"last_message_at": conversation.LastMessageAt.AsTime(),
	}

	response.Success(c, conversationMap)
}

// ListConversations lists conversations
func (h *ConversationHandler) ListConversations(c *gin.Context) {
	agentID := c.Query("agent_id")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	now := timestamppb.New(time.Now())
	conversations := []*pb.Conversation{
		{
			Id:            uuid.New().String(),
			AgentId:       agentID,
			UserId:        "user-001",
			Title:         "Conversation 1",
			Status:        "active",
			CreatedAt:     now,
			UpdatedAt:     now,
			LastMessageAt: now,
		},
	}

	// Convert to map
	conversationMaps := make([]map[string]interface{}, len(conversations))
	for i, conv := range conversations {
		conversationMaps[i] = map[string]interface{}{
			"id":              conv.Id,
			"agent_id":        conv.AgentId,
			"user_id":         conv.UserId,
			"title":           conv.Title,
			"status":          conv.Status,
			"last_message_at": conv.LastMessageAt.AsTime(),
			"created_at":      conv.CreatedAt.AsTime(),
		}
	}

	responseData := map[string]interface{}{
		"items":     conversationMaps,
		"page":      page,
		"page_size": pageSize,
		"total":     int64(1),
	}

	response.Success(c, responseData)
}
