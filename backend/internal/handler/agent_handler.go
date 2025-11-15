package handler

import (
	"strconv"
	"time"

	pb "agent-platform/api/proto/gen"
	"agent-platform/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AgentHandler struct {
	// service will be injected
}

func NewAgentHandler() *AgentHandler {
	return &AgentHandler{}
}

// CreateAgent creates a new agent
func (h *AgentHandler) CreateAgent(c *gin.Context) {
	var req pb.CreateAgentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Get user from context (set by auth middleware)
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "system" // default for demo
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
		CreatedBy:      userID,
		Tags:           req.Tags,
		Folder:         req.Folder,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// 将proto message转换为map以便JSON序列化
	agentMap := map[string]interface{}{
		"id":              agent.Id,
		"name":            agent.Name,
		"description":     agent.Description,
		"type":            agent.Type,
		"status":          agent.Status,
		"version":         agent.Version,
		"created_by":      agent.CreatedBy,
		"tags":            agent.Tags,
		"folder":          agent.Folder,
		"prompt_template": agent.PromptTemplate,
		"created_at":      agent.CreatedAt.AsTime(),
		"updated_at":      agent.UpdatedAt.AsTime(),
	}

	// 添加可选字段
	if agent.ModelConfig != nil {
		agentMap["model_config"] = agent.ModelConfig.AsMap()
	}
	if agent.Parameters != nil {
		agentMap["parameters"] = agent.Parameters.AsMap()
	}
	if len(agent.Tools) > 0 {
		agentMap["tools"] = agent.Tools
	}
	if len(agent.KnowledgeBases) > 0 {
		agentMap["knowledge_bases"] = agent.KnowledgeBases
	}

	response.Created(c, agentMap)
}

// ListAgents returns a list of agents
func (h *AgentHandler) ListAgents(c *gin.Context) {
	// Parse query parameters
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	status := c.Query("status")
	agentType := c.Query("type")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	// Mock response with proto
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

	// 转换为map数组
	agentMaps := make([]map[string]interface{}, len(agents))
	for i, agent := range agents {
		agentMaps[i] = map[string]interface{}{
			"id":          agent.Id,
			"name":        agent.Name,
			"description": agent.Description,
			"type":        agent.Type,
			"status":      agent.Status,
			"version":     agent.Version,
			"created_by":  agent.CreatedBy,
			"created_at":  agent.CreatedAt.AsTime(),
			"updated_at":  agent.UpdatedAt.AsTime(),
		}
	}

	responseData := map[string]interface{}{
		"items":     agentMaps,
		"page":      page,
		"page_size": pageSize,
		"total":     int64(1),
	}

	// 添加筛选参数到响应（如果存在）
	if status != "" {
		responseData["status"] = status
	}
	if agentType != "" {
		responseData["type"] = agentType
	}

	response.Success(c, responseData)
}

// GetAgent returns agent details
func (h *AgentHandler) GetAgent(c *gin.Context) {
	id := c.Param("id")

	// Mock model config and parameters
	modelConfig, _ := structpb.NewStruct(map[string]interface{}{
		"model":       "gpt-4",
		"temperature": 0.7,
	})

	now := timestamppb.New(time.Now())
	agent := &pb.Agent{
		Id:             id,
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

	agentMap := map[string]interface{}{
		"id":              agent.Id,
		"name":            agent.Name,
		"description":     agent.Description,
		"type":            agent.Type,
		"status":          agent.Status,
		"version":         agent.Version,
		"created_by":      agent.CreatedBy,
		"prompt_template": agent.PromptTemplate,
		"tools":           agent.Tools,
		"knowledge_bases": agent.KnowledgeBases,
		"created_at":      agent.CreatedAt.AsTime(),
		"updated_at":      agent.UpdatedAt.AsTime(),
	}

	if agent.ModelConfig != nil {
		agentMap["model_config"] = agent.ModelConfig.AsMap()
	}

	response.Success(c, agentMap)
}

// UpdateAgent updates an existing agent
func (h *AgentHandler) UpdateAgent(c *gin.Context) {
	id := c.Param("id")

	var req pb.UpdateAgentRequest
	req.Id = id

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	now := timestamppb.New(time.Now())
	agent := &pb.Agent{
		Id:             id,
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

	agentMap := map[string]interface{}{
		"id":              agent.Id,
		"name":            agent.Name,
		"description":     agent.Description,
		"status":          agent.Status,
		"prompt_template": agent.PromptTemplate,
		"updated_at":      agent.UpdatedAt.AsTime(),
	}

	if agent.ModelConfig != nil {
		agentMap["model_config"] = agent.ModelConfig.AsMap()
	}
	if agent.Parameters != nil {
		agentMap["parameters"] = agent.Parameters.AsMap()
	}
	if len(agent.Tools) > 0 {
		agentMap["tools"] = agent.Tools
	}
	if len(agent.KnowledgeBases) > 0 {
		agentMap["knowledge_bases"] = agent.KnowledgeBases
	}

	response.Success(c, agentMap)
}

// DeleteAgent deletes an agent
func (h *AgentHandler) DeleteAgent(c *gin.Context) {
	id := c.Param("id")

	deleteResp := &pb.DeleteAgentResponse{
		Id: id,
	}

	response.SuccessWithMessage(c, "Agent deleted successfully", map[string]interface{}{
		"id": deleteResp.Id,
	})
}
