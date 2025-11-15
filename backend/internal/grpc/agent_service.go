package grpc

import (
	"context"
	"time"

	pb "agent-platform/gen/go"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AgentServer gRPC Agent 服务实现
type AgentServer struct {
	pb.UnimplementedAgentServiceServer
	// 这里可以注入 service 层
}

// NewAgentServer 创建 Agent 服务实例
func NewAgentServer() *AgentServer {
	return &AgentServer{}
}

// CreateAgent 创建 Agent
func (s *AgentServer) CreateAgent(ctx context.Context, req *pb.CreateAgentRequest) (*pb.Agent, error) {
	// 验证请求
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	now := timestamppb.New(time.Now())
	agent := &pb.Agent{
		Id:             uuid.New().String(),
		Name:           req.Name,
		Description:    req.Description,
		Type:           req.Type,
		ModelConfig:    req.ModelConfig,
		Tools:          req.Tools,
		KnowledgeBases: req.KnowledgeBases,
		PromptTemplate: req.PromptTemplate,
		Parameters:     req.Parameters,
		Status:         "draft",
		Version:        "1.0.0",
		CreatedBy:      "system", // TODO: 从 context 中获取用户信息
		Tags:           req.Tags,
		Folder:         req.Folder,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	return agent, nil
}

// ListAgents 获取 Agent 列表
func (s *AgentServer) ListAgents(ctx context.Context, req *pb.ListAgentsRequest) (*pb.ListAgentsResponse, error) {
	// Mock 数据
	now := timestamppb.New(time.Now())
	agents := []*pb.Agent{
		{
			Id:          uuid.New().String(),
			Name:        "Customer Service Agent",
			Description: "AI agent for customer service",
			Type:        "single",
			Status:      "published",
			Version:     "1.0.0",
			CreatedBy:   "admin",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	return &pb.ListAgentsResponse{
		Items:    agents,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    1,
	}, nil
}

// GetAgent 获取 Agent 详情
func (s *AgentServer) GetAgent(ctx context.Context, req *pb.GetAgentRequest) (*pb.Agent, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// Mock model config
	modelConfig, _ := structpb.NewStruct(map[string]interface{}{
		"model":       "gpt-4",
		"temperature": 0.7,
	})

	now := timestamppb.New(time.Now())
	agent := &pb.Agent{
		Id:             req.Id,
		Name:           "Customer Service Agent",
		Description:    "AI agent for customer service",
		Type:           "single",
		ModelConfig:    modelConfig,
		Tools:          []string{"search", "email"},
		KnowledgeBases: []string{"kb-001"},
		PromptTemplate: "You are a helpful customer service agent.",
		Status:         "published",
		Version:        "1.0.0",
		CreatedBy:      "admin",
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	return agent, nil
}

// UpdateAgent 更新 Agent
func (s *AgentServer) UpdateAgent(ctx context.Context, req *pb.UpdateAgentRequest) (*pb.Agent, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	now := timestamppb.New(time.Now())
	agent := &pb.Agent{
		Id:             req.Id,
		Name:           req.Name,
		Description:    req.Description,
		ModelConfig:    req.ModelConfig,
		Tools:          req.Tools,
		KnowledgeBases: req.KnowledgeBases,
		PromptTemplate: req.PromptTemplate,
		Parameters:     req.Parameters,
		Status:         req.Status,
		UpdatedAt:      now,
	}

	return agent, nil
}

// DeleteAgent 删除 Agent
func (s *AgentServer) DeleteAgent(ctx context.Context, req *pb.DeleteAgentRequest) (*emptypb.Empty, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// TODO: 实际删除逻辑

	return &emptypb.Empty{}, nil
}
