package grpc

import (
	"context"
	"time"

	pb "agent-platform/gen/go"
	"agent-platform/internal/model/ent"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToolServer gRPC Tool 服务实现
type ToolServer struct {
	pb.UnimplementedToolServiceServer
	client *ent.Client
}

// NewToolServer 创建 Tool 服务实例
func NewToolServer(client *ent.Client) *ToolServer {
	return &ToolServer{
		client: client,
	}
}

// CreateTool 创建工具
func (s *ToolServer) CreateTool(ctx context.Context, req *pb.CreateToolRequest) (*pb.Tool, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	now := timestamppb.New(time.Now())
	tool := &pb.Tool{
		Id:             uuid.New().String(),
		Name:           req.Name,
		Description:    req.Description,
		Type:           req.Type,
		Schema:         req.Schema,
		Implementation: req.Implementation,
		Version:        "1.0.0",
		IsPublic:       false,
		CreatedBy:      "system", // TODO: 从 context 获取
		Category:       req.Category,
		Tags:           req.Tags,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	return tool, nil
}

// ListTools 获取工具列表
func (s *ToolServer) ListTools(ctx context.Context, req *pb.ListToolsRequest) (*pb.ListToolsResponse, error) {
	now := timestamppb.New(time.Now())
	tools := []*pb.Tool{
		{
			Id:          "tool-001",
			Name:        "Web Search",
			Description: "Search the web for information",
			Type:        "api",
			Category:    "search",
			IsPublic:    true,
			Version:     "1.0.0",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Id:          "tool-002",
			Name:        "Send Email",
			Description: "Send emails to users",
			Type:        "api",
			Category:    "communication",
			IsPublic:    true,
			Version:     "1.0.0",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	return &pb.ListToolsResponse{
		Items:    tools,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    2,
	}, nil
}

// GetTool 获取工具详情
func (s *ToolServer) GetTool(ctx context.Context, req *pb.GetToolRequest) (*pb.Tool, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// Mock schema
	schema, _ := structpb.NewStruct(map[string]interface{}{
		"type": "function",
		"function": map[string]interface{}{
			"name":        "web_search",
			"description": "Search the web",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Search query",
					},
				},
				"required": []string{"query"},
			},
		},
	})

	now := timestamppb.New(time.Now())
	tool := &pb.Tool{
		Id:          req.Id,
		Name:        "Web Search",
		Description: "Search the web for information",
		Type:        "api",
		Schema:      schema,
		Category:    "search",
		IsPublic:    true,
		Version:     "1.0.0",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return tool, nil
}

// DeleteTool 删除工具
func (s *ToolServer) DeleteTool(ctx context.Context, req *pb.DeleteToolRequest) (*emptypb.Empty, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// TODO: 实际删除逻辑

	return &emptypb.Empty{}, nil
}
