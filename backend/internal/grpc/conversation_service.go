package grpc

import (
	"context"
	"time"

	pb "agent-platform/gen/go"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ConversationServer gRPC Conversation 服务实现
type ConversationServer struct {
	pb.UnimplementedConversationServiceServer
}

// NewConversationServer 创建 Conversation 服务实例
func NewConversationServer() *ConversationServer {
	return &ConversationServer{}
}

// CreateConversation 创建对话
func (s *ConversationServer) CreateConversation(ctx context.Context, req *pb.CreateConversationRequest) (*pb.Conversation, error) {
	if req.AgentId == "" {
		return nil, status.Error(codes.InvalidArgument, "agent_id is required")
	}

	now := timestamppb.New(time.Now())
	conversation := &pb.Conversation{
		Id:            uuid.New().String(),
		AgentId:       req.AgentId,
		UserId:        "anonymous", // TODO: 从 context 获取
		Title:         req.Title,
		Messages:      []*pb.Message{},
		Context:       req.Context,
		Status:        "active",
		CreatedAt:     now,
		UpdatedAt:     now,
		LastMessageAt: now,
	}

	return conversation, nil
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

	now := timestamppb.New(time.Now())

	// 用户消息
	userMessage := &pb.Message{
		Id:        uuid.New().String(),
		Role:      "user",
		Content:   req.Content,
		Timestamp: now,
	}

	// Mock AI 响应
	assistantMessage := &pb.Message{
		Id:        uuid.New().String(),
		Role:      "assistant",
		Content:   "This is a mock response from the agent. In production, this would call the AI model.",
		Timestamp: timestamppb.New(time.Now().Add(1 * time.Second)),
	}

	return &pb.SendMessageResponse{
		ConversationId: req.ConversationId,
		Messages:       []*pb.Message{userMessage, assistantMessage},
	}, nil
}
