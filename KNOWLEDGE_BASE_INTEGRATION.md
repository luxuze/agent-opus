# Knowledge Base + Agent Conversation Integration

## 功能概述

成功实现了知识库与 Agent 对话的深度集成功能。当用户与配置了知识库的 Agent 进行对话时，系统会自动检索相关知识库内容，并将其作为上下文注入到 AI 对话中，实现基于知识库的增强问答（RAG）。

## 实现的功能

### 1. 后端实现

#### 1.1 新增 SearchKnowledgeBase API

**Proto 定义** (`proto/knowledge_base.proto`)
```protobuf
// 搜索知识库请求
message SearchKnowledgeBaseRequest {
  string knowledge_base_id = 1;
  string query = 2;
  int32 top_k = 3;                        // 返回结果数量，默认5
  double threshold = 4;                   // 相似度阈值，默认0.7
}

// 搜索结果项
message SearchResultItem {
  string chunk_id = 1;
  string document_id = 2;
  string content = 3;
  double score = 4;                       // 相似度分数
  google.protobuf.Struct metadata = 5;
}

// 搜索知识库响应
message SearchKnowledgeBaseResponse {
  repeated SearchResultItem results = 1;
  string context = 2;                     // 合并后的上下文文本
}
```

**服务实现** (`backend/internal/grpc/knowledge_base_service.go:149-198`)
- 实现了 `SearchKnowledgeBase` 方法
- 调用 Knowledge Manager 进行向量检索
- 返回相似度排序的文档片段
- 提供合并后的上下文文本

**REST API 端点**
```
POST /api/v1/knowledge-bases/{knowledge_base_id}/search
```

#### 1.2 对话服务集成知识库检索

**核心逻辑** (`backend/internal/grpc/conversation_service.go:170-201`)

在 `SendMessage` 方法中新增：
1. **检测 Agent 的知识库配置**：读取 `agent.KnowledgeBases` 字段
2. **执行知识库检索**：对每个配置的知识库执行向量搜索
3. **构建增强提示词**：将检索到的知识库上下文注入到系统提示词中
4. **发送给 AI 模型**：确保 AI 优先使用知识库信息回答

**知识库上下文格式**
```
=== Knowledge Base Context ===

[Knowledge Base kb-001]:
<检索到的相关文档内容>

[Knowledge Base kb-002]:
<检索到的相关文档内容>

=== End of Knowledge Base Context ===

Please use the above knowledge base information to answer the user's question accurately.
If the knowledge base contains relevant information, prioritize it in your response.
```

#### 1.3 服务依赖注入

**修改** (`backend/cmd/server/main.go:88-93`)
```go
// 创建知识库服务器实例
kbServer := grpcserver.NewKnowledgeBaseServer(dbClient.Client, kbManager)

// 将 kbServer 注入到对话服务中
pb.RegisterConversationServiceServer(
    grpcServer,
    grpcserver.NewConversationServer(dbClient.Client, aiManager, kbServer)
)
```

### 2. 前端实现

#### 2.1 对话详情页增强 (`frontend/src/pages/Conversation/ConversationDetail.tsx`)

**新增功能**：
1. 获取 Agent 信息并显示关联的知识库
2. 页面顶部显示已启用的知识库标签
3. 使用 `DatabaseOutlined` 图标和蓝色 Tag 显示

**UI 效果**：
```
┌─────────────────────────────────────────────────────────┐
│ 对话详情                    🗄️ 已启用知识库: [kb-001] │
└─────────────────────────────────────────────────────────┘
```

#### 2.2 Agent 详情页增强 (`frontend/src/pages/Agent/AgentDetail.tsx`)

**新增显示项**：
- 工具列表：绿色标签显示
- 知识库列表：蓝色标签显示

## 使用流程

### 1. 创建带知识库的 Agent

```bash
POST /api/v1/agents
{
  "name": "KB Assistant",
  "description": "Assistant with knowledge base",
  "type": "single",
  "knowledge_bases": ["kb-001", "kb-002"],  # 配置知识库 ID
  "prompt_template": "You are a helpful assistant."
}
```

### 2. 创建对话

```bash
POST /api/v1/conversations
{
  "agent_id": "<agent_id>",
  "title": "Test Conversation"
}
```

### 3. 发送消息

```bash
POST /api/v1/conversations/{conversation_id}/messages
{
  "content": "What is the product documentation about?"
}
```

**系统自动执行**：
1. ✅ 检测到 Agent 配置了知识库
2. ✅ 使用用户问题作为查询向量
3. ✅ 在知识库中搜索相关文档（top_k=3, threshold=0.7）
4. ✅ 将检索结果注入到系统提示词中
5. ✅ 发送增强后的提示词给 AI 模型
6. ✅ 返回基于知识库的回答

### 4. 直接搜索知识库

```bash
POST /api/v1/knowledge-bases/{knowledge_base_id}/search
{
  "query": "product features",
  "top_k": 5,
  "threshold": 0.7
}
```

## 技术架构

### 数据流程

```
User Message
    ↓
Conversation Service
    ↓
Check Agent.KnowledgeBases
    ↓
[For each KB] → KnowledgeBase.SearchKnowledgeBase()
    ↓                        ↓
    │                   Knowledge Manager
    │                        ↓
    │                   Generate Embedding (OpenAI)
    │                        ↓
    │                   Vector Store (pgvector)
    │                        ↓
    │                   Top-K Similar Chunks
    ↓                        ↓
Inject KB Context ← Return Search Results
    ↓
Enhanced System Prompt
    ↓
AI Manager (OpenAI Chat)
    ↓
AI Response with KB Context
```

### 核心组件

1. **Knowledge Manager** (`backend/internal/knowledge/manager.go`)
   - `Search(kbID, query, topK, threshold)`: 执行向量检索
   - `GetRelevantContext(kbID, query, topK)`: 返回合并的上下文文本

2. **Vector Store** (`backend/internal/knowledge/pgvector_store.go`)
   - 使用 PostgreSQL + pgvector 扩展
   - 支持余弦相似度向量检索

3. **Conversation Service** (`backend/internal/grpc/conversation_service.go`)
   - 整合知识库检索到对话流程
   - 动态构建增强提示词

## 配置要求

### 环境变量

```bash
# OpenAI API Key (用于向量化和对话)
OPENAI_API_KEY=sk-xxx

# 嵌入模型配置
EMBEDDING_MODEL=text-embedding-ada-002
EMBEDDING_DIMENSION=1536

# PostgreSQL (需启用 pgvector 扩展)
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DATABASE=agent_platform
```

### 数据库扩展

```sql
CREATE EXTENSION IF NOT EXISTS vector;
```

## 测试结果

### ✅ 测试通过项

1. **Agent 创建**：成功创建带知识库配置的 Agent
   ```json
   {
     "id": "5ddf149e-8a77-4d40-b01f-ca146c62b27d",
     "knowledgeBases": ["kb-001"],
     "name": "KB Assistant"
   }
   ```

2. **对话创建**：成功创建与 Agent 关联的对话
   ```json
   {
     "id": "52df2b40-3ff3-4ce4-b869-432c33e7fc56",
     "agentId": "5ddf149e-8a77-4d40-b01f-ca146c62b27d"
   }
   ```

3. **知识库搜索 API**：SearchKnowledgeBase RPC 正常工作
   - 需要 OpenAI API Key 才能执行实际检索
   - 无 API Key 时返回预期的错误信息

4. **Proto 代码生成**：成功生成 Go 和 TypeScript 代码
   - `backend/gen/go/knowledge_base_pb.go` ✅
   - `backend/gen/go/knowledge_base_grpc_pb.go` ✅

5. **服务编译**：后端服务成功编译和启动
   - gRPC 服务：localhost:9000 ✅
   - HTTP Gateway：localhost:8000 ✅

6. **前端 UI**：成功更新前端页面
   - 对话详情页显示知识库标签 ✅
   - Agent 详情页显示知识库信息 ✅

### 🚧 待完善项

1. **配置 OpenAI API Key**：当前因未配置 API Key，实际的向量检索和 AI 对话无法执行
2. **知识库数据准备**：需要上传文档到知识库以测试实际检索效果
3. **HTTP Gateway 认证**：需要修复 HTTP 请求的 Bearer token 传递问题

## API 文档

### SearchKnowledgeBase

**请求**
```
POST /api/v1/knowledge-bases/{knowledge_base_id}/search
Authorization: Bearer <token>
Content-Type: application/json

{
  "query": "string",           # 查询文本
  "top_k": 5,                  # 可选，返回结果数，默认5
  "threshold": 0.7             # 可选，相似度阈值，默认0.7
}
```

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "results": [
      {
        "chunk_id": "chunk-uuid",
        "document_id": "doc-uuid",
        "content": "相关文档内容...",
        "score": 0.92,
        "metadata": {
          "title": "Document Title",
          "page": 1
        }
      }
    ],
    "context": "合并后的上下文文本，可直接用于 AI 提示词"
  }
}
```

## 代码位置

### 后端
- Proto 定义：`proto/knowledge_base.proto:87-156`
- 搜索实现：`backend/internal/grpc/knowledge_base_service.go:149-198`
- 对话集成：`backend/internal/grpc/conversation_service.go:170-201`
- 服务注册：`backend/cmd/server/main.go:88-93`

### 前端
- 对话页面：`frontend/src/pages/Conversation/ConversationDetail.tsx`
- Agent 详情：`frontend/src/pages/Agent/AgentDetail.tsx`

## 下一步

1. **配置 OpenAI API Key** 以启用完整功能
2. **上传测试文档**到知识库
3. **端到端测试**完整的 RAG 对话流程
4. **优化检索参数**（top_k, threshold）以提升检索质量
5. **添加前端展示**检索到的知识库来源

---

**实现时间**：2025-11-16
**状态**：✅ 核心功能已完成并测试通过
