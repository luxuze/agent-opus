package grpc

import (
	"context"
	"time"

	pb "agent-platform/gen/go"
	"agent-platform/internal/knowledge"
	"agent-platform/internal/model/ent"
	"agent-platform/internal/repository"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// KnowledgeBaseServer gRPC KnowledgeBase 服务实现
type KnowledgeBaseServer struct {
	pb.UnimplementedKnowledgeBaseServiceServer
	client   *ent.Client
	repo     *repository.KnowledgeBaseRepository
	kbMgr    *knowledge.Manager
}

// NewKnowledgeBaseServer 创建 KnowledgeBase 服务实例
func NewKnowledgeBaseServer(client *ent.Client, kbMgr *knowledge.Manager) *KnowledgeBaseServer {
	return &KnowledgeBaseServer{
		client: client,
		repo:   repository.NewKnowledgeBaseRepository(client),
		kbMgr:  kbMgr,
	}
}

// CreateKnowledgeBase 创建知识库
func (s *KnowledgeBaseServer) CreateKnowledgeBase(ctx context.Context, req *pb.CreateKnowledgeBaseRequest) (*pb.KnowledgeBase, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	now := timestamppb.New(time.Now())
	kb := &pb.KnowledgeBase{
		Id:             uuid.New().String(),
		Name:           req.Name,
		Description:    req.Description,
		Type:           req.Type,
		EmbeddingModel: req.EmbeddingModel,
		ChunkConfig:    req.ChunkConfig,
		CreatedBy:      "system", // TODO: 从 context 获取
		DocumentCount:  0,
		VectorCount:    0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	return kb, nil
}

// ListKnowledgeBases 获取知识库列表
func (s *KnowledgeBaseServer) ListKnowledgeBases(ctx context.Context, req *pb.ListKnowledgeBasesRequest) (*pb.ListKnowledgeBasesResponse, error) {
	now := timestamppb.New(time.Now())
	kbs := []*pb.KnowledgeBase{
		{
			Id:            "kb-001",
			Name:          "Product Documentation",
			Description:   "Product documentation knowledge base",
			Type:          "document",
			DocumentCount: 150,
			VectorCount:   5000,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}

	return &pb.ListKnowledgeBasesResponse{
		Items:    kbs,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    1,
	}, nil
}

// GetKnowledgeBase 获取知识库详情
func (s *KnowledgeBaseServer) GetKnowledgeBase(ctx context.Context, req *pb.GetKnowledgeBaseRequest) (*pb.KnowledgeBase, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// Mock chunk config
	chunkConfig, _ := structpb.NewStruct(map[string]interface{}{
		"chunk_size":    1000,
		"chunk_overlap": 200,
	})

	now := timestamppb.New(time.Now())
	kb := &pb.KnowledgeBase{
		Id:             req.Id,
		Name:           "Product Documentation",
		Description:    "Product documentation knowledge base",
		Type:           "document",
		EmbeddingModel: "text-embedding-ada-002",
		ChunkConfig:    chunkConfig,
		DocumentCount:  150,
		VectorCount:    5000,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	return kb, nil
}

// UploadDocument 上传文档
func (s *KnowledgeBaseServer) UploadDocument(ctx context.Context, req *pb.UploadDocumentRequest) (*pb.Document, error) {
	if req.KnowledgeBaseId == "" {
		return nil, status.Error(codes.InvalidArgument, "knowledge_base_id is required")
	}
	if req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	now := timestamppb.New(time.Now())
	document := &pb.Document{
		Id:              uuid.New().String(),
		KnowledgeBaseId: req.KnowledgeBaseId,
		Title:           req.Title,
		Content:         req.Content,
		Metadata:        req.Metadata,
		Status:          "processing",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	return document, nil
}

// DeleteKnowledgeBase 删除知识库
func (s *KnowledgeBaseServer) DeleteKnowledgeBase(ctx context.Context, req *pb.DeleteKnowledgeBaseRequest) (*emptypb.Empty, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// TODO: 实际删除逻辑

	return &emptypb.Empty{}, nil
}
