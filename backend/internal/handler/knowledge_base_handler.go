package handler

import (
	"strconv"
	"time"

	pb "agent-platform/gen/go"
	"agent-platform/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type KnowledgeBaseHandler struct{}

func NewKnowledgeBaseHandler() *KnowledgeBaseHandler {
	return &KnowledgeBaseHandler{}
}

// CreateKnowledgeBase creates a new knowledge base
func (h *KnowledgeBaseHandler) CreateKnowledgeBase(c *gin.Context) {
	var req pb.CreateKnowledgeBaseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		userID = "system"
	}

	now := timestamppb.New(time.Now())
	kb := &pb.KnowledgeBase{
		Id:             uuid.New().String(),
		Name:           req.Name,
		Description:    req.Description,
		Type:           req.Type,
		EmbeddingModel: req.EmbeddingModel,
		ChunkConfig:    req.ChunkConfig,
		CreatedBy:      userID,
		DocumentCount:  0,
		VectorCount:    0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	kbMap := map[string]interface{}{
		"id":              kb.Id,
		"name":            kb.Name,
		"description":     kb.Description,
		"type":            kb.Type,
		"embedding_model": kb.EmbeddingModel,
		"created_by":      kb.CreatedBy,
		"document_count":  kb.DocumentCount,
		"vector_count":    kb.VectorCount,
		"created_at":      kb.CreatedAt.AsTime(),
		"updated_at":      kb.UpdatedAt.AsTime(),
	}

	if kb.ChunkConfig != nil {
		kbMap["chunk_config"] = kb.ChunkConfig.AsMap()
	}

	response.Created(c, kbMap)
}

// ListKnowledgeBases lists knowledge bases
func (h *KnowledgeBaseHandler) ListKnowledgeBases(c *gin.Context) {
	kbType := c.Query("type")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

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

	// Convert to map
	kbMaps := make([]map[string]interface{}, len(kbs))
	for i, kb := range kbs {
		kbMaps[i] = map[string]interface{}{
			"id":             kb.Id,
			"name":           kb.Name,
			"description":    kb.Description,
			"type":           kb.Type,
			"document_count": kb.DocumentCount,
			"vector_count":   kb.VectorCount,
		}
	}

	responseData := map[string]interface{}{
		"items":     kbMaps,
		"page":      page,
		"page_size": pageSize,
		"total":     int64(1),
	}

	if kbType != "" {
		responseData["type"] = kbType
	}

	response.Success(c, responseData)
}

// GetKnowledgeBase returns knowledge base details
func (h *KnowledgeBaseHandler) GetKnowledgeBase(c *gin.Context) {
	id := c.Param("id")

	// Mock chunk config
	chunkConfig, _ := structpb.NewStruct(map[string]interface{}{
		"chunk_size":    1000,
		"chunk_overlap": 200,
	})

	now := timestamppb.New(time.Now())
	kb := &pb.KnowledgeBase{
		Id:             id,
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

	kbMap := map[string]interface{}{
		"id":              kb.Id,
		"name":            kb.Name,
		"description":     kb.Description,
		"type":            kb.Type,
		"embedding_model": kb.EmbeddingModel,
		"document_count":  kb.DocumentCount,
		"vector_count":    kb.VectorCount,
		"created_at":      kb.CreatedAt.AsTime(),
		"updated_at":      kb.UpdatedAt.AsTime(),
	}

	if kb.ChunkConfig != nil {
		kbMap["chunk_config"] = kb.ChunkConfig.AsMap()
	}

	response.Success(c, kbMap)
}

// UploadDocument uploads a document to knowledge base
func (h *KnowledgeBaseHandler) UploadDocument(c *gin.Context) {
	id := c.Param("id")

	var req pb.UploadDocumentRequest
	req.KnowledgeBaseId = id

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	now := timestamppb.New(time.Now())
	document := &pb.Document{
		Id:              uuid.New().String(),
		KnowledgeBaseId: id,
		Title:           req.Title,
		Content:         req.Content,
		Metadata:        req.Metadata,
		Status:          "processing",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	docMap := map[string]interface{}{
		"id":                document.Id,
		"knowledge_base_id": document.KnowledgeBaseId,
		"title":             document.Title,
		"content":           document.Content,
		"status":            document.Status,
		"created_at":        document.CreatedAt.AsTime(),
		"updated_at":        document.UpdatedAt.AsTime(),
	}

	if document.Metadata != nil {
		docMap["metadata"] = document.Metadata.AsMap()
	}

	response.SuccessWithMessage(c, "Document uploaded successfully", docMap)
}

// DeleteKnowledgeBase deletes a knowledge base
func (h *KnowledgeBaseHandler) DeleteKnowledgeBase(c *gin.Context) {
	id := c.Param("id")

	deleteResp := &pb.DeleteKnowledgeBaseResponse{
		Id: id,
	}

	response.SuccessWithMessage(c, "Knowledge base deleted successfully", map[string]interface{}{
		"id": deleteResp.Id,
	})
}
