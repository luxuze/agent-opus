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
	kbServer    *KnowledgeBaseServer
}

// NewConversationServer 创建 Conversation 服务实例
func NewConversationServer(client *ent.Client, aiManager *ai.Manager, kbServer *KnowledgeBaseServer) *ConversationServer {
	return &ConversationServer{
		client:    client,
		aiManager: aiManager,
		convRepo:  repository.NewConversationRepository(client),
		agentRepo: repository.NewAgentRepository(client),
		kbServer:  kbServer,
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

	// Get conversation from database
	conv, err := s.convRepo.Get(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "conversation not found: %v", err)
	}

	// Convert to protobuf
	return entConversationToProto(conv), nil
}

// ListConversations 获取对话列表
func (s *ConversationServer) ListConversations(ctx context.Context, req *pb.ListConversationsRequest) (*pb.ListConversationsResponse, error) {
	// Set default pagination
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	// List conversations from database
	// Note: proto only has agent_id, so we pass empty strings for user_id and status
	conversations, total, err := s.convRepo.List(ctx, page, pageSize, req.AgentId, "", "")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list conversations: %v", err)
	}

	// Convert to protobuf
	items := make([]*pb.Conversation, len(conversations))
	for i, conv := range conversations {
		items[i] = entConversationToProto(conv)
	}

	return &pb.ListConversationsResponse{
		Items:    items,
		Page:     page,
		PageSize: pageSize,
		Total:    int64(total),
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

	// Add system prompt with knowledge base context if agent has one
	systemPrompt := ""
	if agent.PromptTemplate != "" {
		systemPrompt = agent.PromptTemplate
	}

	// Retrieve knowledge base context if agent has knowledge bases configured
	if len(agent.KnowledgeBases) > 0 && s.kbServer != nil {
		kbContexts := []string{}
		for _, kbID := range agent.KnowledgeBases {
			// Search each knowledge base
			searchResp, err := s.kbServer.SearchKnowledgeBase(ctx, &pb.SearchKnowledgeBaseRequest{
				KnowledgeBaseId: kbID,
				Query:           req.Content,
				TopK:            3,
				Threshold:       0.7,
			})
			if err == nil && searchResp.Context != "" {
				kbContexts = append(kbContexts, searchResp.Context)
			}
		}

		// Append knowledge base context to system prompt
		if len(kbContexts) > 0 {
			kbContext := "\n\n=== Knowledge Base Context ===\n"
			for i, ctx := range kbContexts {
				kbContext += "\n[Knowledge Base " + agent.KnowledgeBases[i] + "]:\n" + ctx
			}
			kbContext += "\n=== End of Knowledge Base Context ===\n\n"
			kbContext += "Please use the above knowledge base information to answer the user's question accurately. If the knowledge base contains relevant information, prioritize it in your response."

			if systemPrompt != "" {
				systemPrompt = systemPrompt + kbContext
			} else {
				systemPrompt = kbContext
			}
		}
	}

	// Add system prompt to messages
	if systemPrompt != "" {
		messages = append(messages, ai.Message{
			Role:    "system",
			Content: systemPrompt,
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
	model := "deepseek-ai/DeepSeek-V3" // Default to DeepSeek (更经济实惠)
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
