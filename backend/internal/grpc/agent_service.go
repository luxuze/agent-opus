package grpc

import (
	"context"

	pb "agent-platform/gen/go"
	"agent-platform/internal/model/ent"
	"agent-platform/internal/repository"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AgentServer gRPC Agent 服务实现
type AgentServer struct {
	pb.UnimplementedAgentServiceServer
	client *ent.Client
	repo   *repository.AgentRepository
}

// NewAgentServer 创建 Agent 服务实例
func NewAgentServer(client *ent.Client) *AgentServer {
	return &AgentServer{
		client: client,
		repo:   repository.NewAgentRepository(client),
	}
}

// CreateAgent 创建 Agent
func (s *AgentServer) CreateAgent(ctx context.Context, req *pb.CreateAgentRequest) (*pb.Agent, error) {
	// 验证请求
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	// 准备数据
	agentID := uuid.New().String()
	agentType := req.Type
	if agentType == "" {
		agentType = "single"
	}

	// 创建 ent entity
	entAgent := &ent.Agent{
		ID:          agentID,
		Name:        req.Name,
		Description: req.Description,
		Type:        agentType,
		Status:      "draft",
		Version:     "1.0.0",
		CreatedBy:   "system", // TODO: 从 context 中获取用户信息
	}

	// 设置可选字段
	if req.ModelConfig != nil {
		entAgent.ModelConfig = req.ModelConfig.AsMap()
	}
	if req.Tools != nil {
		entAgent.Tools = req.Tools
	}
	if req.KnowledgeBases != nil {
		entAgent.KnowledgeBases = req.KnowledgeBases
	}
	if req.PromptTemplate != "" {
		entAgent.PromptTemplate = req.PromptTemplate
	}
	if req.Parameters != nil {
		entAgent.Parameters = req.Parameters.AsMap()
	}
	if req.Tags != nil {
		entAgent.Tags = req.Tags
	}
	if req.Folder != "" {
		entAgent.Folder = req.Folder
	}

	// 保存到数据库
	created, err := s.repo.Create(ctx, entAgent)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create agent: %v", err)
	}

	// 转换为 protobuf
	return entAgentToProto(created), nil
}

// ListAgents 获取 Agent 列表
func (s *AgentServer) ListAgents(ctx context.Context, req *pb.ListAgentsRequest) (*pb.ListAgentsResponse, error) {
	// 设置默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 从数据库查询
	agents, total, err := s.repo.List(ctx, page, pageSize, req.Status, req.Type, "")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list agents: %v", err)
	}

	// 转换为 protobuf
	pbAgents := make([]*pb.Agent, len(agents))
	for i, agent := range agents {
		pbAgents[i] = entAgentToProto(agent)
	}

	return &pb.ListAgentsResponse{
		Items:    pbAgents,
		Page:     page,
		PageSize: pageSize,
		Total:    int64(total),
	}, nil
}

// GetAgent 获取 Agent 详情
func (s *AgentServer) GetAgent(ctx context.Context, req *pb.GetAgentRequest) (*pb.Agent, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// 从数据库查询
	agent, err := s.repo.Get(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "agent not found: %v", err)
	}

	// 转换为 protobuf
	return entAgentToProto(agent), nil
}

// UpdateAgent 更新 Agent
func (s *AgentServer) UpdateAgent(ctx context.Context, req *pb.UpdateAgentRequest) (*pb.Agent, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// 准备更新字段
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.ModelConfig != nil {
		updates["model_config"] = req.ModelConfig.AsMap()
	}
	if req.Tools != nil {
		updates["tools"] = req.Tools
	}
	if req.KnowledgeBases != nil {
		updates["knowledge_bases"] = req.KnowledgeBases
	}
	if req.PromptTemplate != "" {
		updates["prompt_template"] = req.PromptTemplate
	}
	if req.Parameters != nil {
		updates["parameters"] = req.Parameters.AsMap()
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	// 更新数据库
	updated, err := s.repo.Update(ctx, req.Id, updates)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update agent: %v", err)
	}

	// 转换为 protobuf
	return entAgentToProto(updated), nil
}

// DeleteAgent 删除 Agent
func (s *AgentServer) DeleteAgent(ctx context.Context, req *pb.DeleteAgentRequest) (*emptypb.Empty, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// 从数据库删除
	err := s.repo.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete agent: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// Helper function to convert ent.Agent to pb.Agent
func entAgentToProto(agent *ent.Agent) *pb.Agent {
	pbAgent := &pb.Agent{
		Id:          agent.ID,
		Name:        agent.Name,
		Description: agent.Description,
		Type:        agent.Type,
		Status:      agent.Status,
		Version:     agent.Version,
		CreatedBy:   agent.CreatedBy,
		CreatedAt:   timestamppb.New(agent.CreatedAt),
		UpdatedAt:   timestamppb.New(agent.UpdatedAt),
	}

	// 设置可选字段
	if agent.PromptTemplate != "" {
		pbAgent.PromptTemplate = agent.PromptTemplate
	}
	if agent.Folder != "" {
		pbAgent.Folder = agent.Folder
	}
	if agent.Tools != nil {
		pbAgent.Tools = agent.Tools
	}
	if agent.KnowledgeBases != nil {
		pbAgent.KnowledgeBases = agent.KnowledgeBases
	}
	if agent.Tags != nil {
		pbAgent.Tags = agent.Tags
	}

	return pbAgent
}
