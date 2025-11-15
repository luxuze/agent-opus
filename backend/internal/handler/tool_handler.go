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

type ToolHandler struct{}

func NewToolHandler() *ToolHandler {
	return &ToolHandler{}
}

// CreateTool creates a new tool
func (h *ToolHandler) CreateTool(c *gin.Context) {
	var req pb.CreateToolRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		userID = "system"
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
		CreatedBy:      userID,
		Category:       req.Category,
		Tags:           req.Tags,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	toolMap := map[string]interface{}{
		"id":             tool.Id,
		"name":           tool.Name,
		"description":    tool.Description,
		"type":           tool.Type,
		"implementation": tool.Implementation,
		"version":        tool.Version,
		"is_public":      tool.IsPublic,
		"created_by":     tool.CreatedBy,
		"category":       tool.Category,
		"tags":           tool.Tags,
		"created_at":     tool.CreatedAt.AsTime(),
		"updated_at":     tool.UpdatedAt.AsTime(),
	}

	if tool.Schema != nil {
		toolMap["schema"] = tool.Schema.AsMap()
	}

	response.Created(c, toolMap)
}

// ListTools returns a list of tools
func (h *ToolHandler) ListTools(c *gin.Context) {
	toolType := c.Query("type")
	category := c.Query("category")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

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

	// Convert to map
	toolMaps := make([]map[string]interface{}, len(tools))
	for i, tool := range tools {
		toolMaps[i] = map[string]interface{}{
			"id":          tool.Id,
			"name":        tool.Name,
			"description": tool.Description,
			"type":        tool.Type,
			"category":    tool.Category,
			"is_public":   tool.IsPublic,
			"version":     tool.Version,
		}
	}

	responseData := map[string]interface{}{
		"items":     toolMaps,
		"page":      page,
		"page_size": pageSize,
		"total":     int64(2),
	}

	if toolType != "" {
		responseData["type"] = toolType
	}
	if category != "" {
		responseData["category"] = category
	}

	response.Success(c, responseData)
}

// GetTool returns tool details
func (h *ToolHandler) GetTool(c *gin.Context) {
	id := c.Param("id")

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
		Id:          id,
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

	toolMap := map[string]interface{}{
		"id":          tool.Id,
		"name":        tool.Name,
		"description": tool.Description,
		"type":        tool.Type,
		"category":    tool.Category,
		"is_public":   tool.IsPublic,
		"version":     tool.Version,
		"created_at":  tool.CreatedAt.AsTime(),
		"updated_at":  tool.UpdatedAt.AsTime(),
	}

	if tool.Schema != nil {
		toolMap["schema"] = tool.Schema.AsMap()
	}

	response.Success(c, toolMap)
}

// DeleteTool deletes a tool
func (h *ToolHandler) DeleteTool(c *gin.Context) {
	id := c.Param("id")

	deleteResp := &pb.DeleteToolResponse{
		Id: id,
	}

	response.SuccessWithMessage(c, "Tool deleted successfully", map[string]interface{}{
		"id": deleteResp.Id,
	})
}
