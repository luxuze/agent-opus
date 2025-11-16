package grpc

import (
	"context"
	"time"

	pb "agent-platform/gen/go"
	"agent-platform/internal/ai"
	"agent-platform/internal/model/ent"
	"agent-platform/internal/repository"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ConversationServer gRPC Conversation 服务实现
type ConversationServer struct {
	pb.UnimplementedConversationServiceServer
	client      *ent.Client
	aiManager   *ai.Manager
	convRepo    *repository.ConversationRepository
	agentRepo   *repository.AgentRepository
}

// NewConversationServer 创建 Conversation 服务实例
func NewConversationServer(client *ent.Client, aiManager *ai.Manager) *ConversationServer {
	return &ConversationServer{
		client:    client,
		aiManager: aiManager,
		convRepo:  repository.NewConversationRepository(client),
		agentRepo: repository.NewAgentRepository(client),
	}
}

// CreateConversation 创建对话
func (s *ConversationServer) CreateConversation(ctx context.Context, req *pb.CreateConversationRequest) (*pb.Conversation, error) {
	if req.AgentId == "" {
		return nil, status.Error(codes.InvalidArgument, "agent_id is required")
	}

	// Verify agent exists
	_, err := s.agentRepo.Get(ctx, req.AgentId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "agent not found: %v", err)
	}

	// Create conversation entity
	convID := uuid.New().String()
	title := req.Title
	if title == "" {
		title = "New Conversation"
	}

	entConv := &ent.Conversation{
		ID:       convID,
		AgentID:  req.AgentId,
		UserID:   "anonymous", // TODO: get from context
		Title:    title,
		Messages: []interface{}{},
		Status:   "active",
	}

	// Set context if provided
	if req.Context != nil {
		entConv.Context = req.Context.AsMap()
	}

	// Save to database
	created, err := s.convRepo.Create(ctx, entConv)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create conversation: %v", err)
	}

	// Convert to protobuf
	return entConversationToProto(created), nil
}

// GetConversation 获取对话详情
func (s *ConversationServer) GetConversation(ctx context.Context, req *pb.GetConversationRequest) (*pb.Conversation, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	now := timestamppb.New(time.Now())
	conversation := &pb.Conversation{
		Id:      req.Id,
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

	return conversation, nil
}

// ListConversations 获取对话列表
func (s *ConversationServer) ListConversations(ctx context.Context, req *pb.ListConversationsRequest) (*pb.ListConversationsResponse, error) {
	now := timestamppb.New(time.Now())
	conversations := []*pb.Conversation{
		{
			Id:            uuid.New().String(),
			AgentId:       req.AgentId,
			UserId:        "user-001",
			Title:         "Conversation 1",
			Status:        "active",
			CreatedAt:     now,
			UpdatedAt:     now,
			LastMessageAt: now,
		},
	}

	return &pb.ListConversationsResponse{
		Items:    conversations,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    1,
	}, nil
}

// SendMessage 发送消息
func (s *ConversationServer) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	if req.ConversationId == "" {
		return nil, status.Error(codes.InvalidArgument, "conversation_id is required")
	}
	if req.Content == "" {
		return nil, status.Error(codes.InvalidArgument, "content is required")
	}

	// Get conversation from database
	conv, err := s.convRepo.Get(ctx, req.ConversationId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "conversation not found: %v", err)
	}

	// Get agent
	agent, err := s.agentRepo.Get(ctx, conv.AgentID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "agent not found: %v", err)
	}

	// Create user message
	now := time.Now()
	userMessage := &pb.Message{
		Id:        uuid.New().String(),
		Role:      "user",
		Content:   req.Content,
		Timestamp: timestamppb.New(now),
	}

	// Build message history for AI
	messages := []ai.Message{}

	// Add system prompt if agent has one
	if agent.PromptTemplate != "" {
		messages = append(messages, ai.Message{
			Role:    "system",
			Content: agent.PromptTemplate,
		})
	}

	// Add conversation history
	if conv.Messages != nil {
		for _, msg := range conv.Messages {
			if msgMap, ok := msg.(map[string]interface{}); ok {
				role, _ := msgMap["role"].(string)
				content, _ := msgMap["content"].(string)
				if role != "" && content != "" {
					messages = append(messages, ai.Message{
						Role:    role,
						Content: content,
					})
				}
			}
		}
	}

	// Add new user message
	messages = append(messages, ai.Message{
		Role:    "user",
		Content: req.Content,
	})

	// Get model config from agent
	model := "gpt-4o"
	temperature := float32(0.7)
	if agent.ModelConfig != nil {
		if m, ok := agent.ModelConfig["model"].(string); ok && m != "" {
			model = m
		}
		if t, ok := agent.ModelConfig["temperature"].(float64); ok {
			temperature = float32(t)
		}
	}

	// Call AI service
	aiResp, err := s.aiManager.Chat(ai.ChatRequest{
		Model:       model,
		Messages:    messages,
		Temperature: temperature,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "AI service error: %v", err)
	}

	// Create assistant message
	assistantMessage := &pb.Message{
		Id:        uuid.New().String(),
		Role:      "assistant",
		Content:   aiResp.Content,
		Timestamp: timestamppb.New(time.Now()),
	}

	// Save messages to conversation
	userMsgMap := map[string]interface{}{
		"id":        userMessage.Id,
		"role":      userMessage.Role,
		"content":   userMessage.Content,
		"timestamp": userMessage.Timestamp.AsTime().Unix(),
	}
	assistantMsgMap := map[string]interface{}{
		"id":        assistantMessage.Id,
		"role":      assistantMessage.Role,
		"content":   assistantMessage.Content,
		"timestamp": assistantMessage.Timestamp.AsTime().Unix(),
	}

	_, err = s.convRepo.AddMessage(ctx, req.ConversationId, userMsgMap)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save user message: %v", err)
	}

	_, err = s.convRepo.AddMessage(ctx, req.ConversationId, assistantMsgMap)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save assistant message: %v", err)
	}

	return &pb.SendMessageResponse{
		ConversationId: req.ConversationId,
		Messages:       []*pb.Message{userMessage, assistantMessage},
	}, nil
}
