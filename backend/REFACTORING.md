# 后端重构说明文档

## 重构概述

本次重构主要完成了以下工作：

1. **统一使用 Protocol Buffers 定义接口**
2. **统一 HTTP 响应格式为 `{code, message, data}` 结构**
3. **建立标准化的项目结构和开发流程**
4. **移除 MongoDB 依赖，简化数据库架构**

## 主要改动

### 1. Proto 定义（`backend/api/proto/`）

创建了以下 proto 文件：

- **`common.proto`** - 通用响应结构和分页定义
- **`agent.proto`** - Agent 相关的请求/响应定义
- **`conversation.proto`** - 对话相关的请求/响应定义
- **`tool.proto`** - 工具相关的请求/响应定义
- **`knowledge_base.proto`** - 知识库相关的请求/响应定义

### 2. 统一响应格式

所有 HTTP 接口返回统一的 JSON 格式：

```json
{
  "code": 0,           // 状态码：0表示成功，非0表示错误
  "message": "success", // 消息描述
  "data": {            // 实际数据
    // ... 业务数据
  }
}
```

### 3. 响应工具包（`backend/pkg/response/`）

创建了统一的响应处理工具函数：

- `Success(c, data)` - 返回成功响应
- `SuccessWithMessage(c, message, data)` - 返回带自定义消息的成功响应
- `Created(c, data)` - 返回创建成功响应（HTTP 201）
- `BadRequest(c, message)` - 返回 400 错误
- `Unauthorized(c, message)` - 返回 401 错误
- `NotFound(c, message)` - 返回 404 错误
- `InternalServerError(c, message)` - 返回 500 错误

### 4. Handler 重构

所有 Handler 已重构为使用 proto 定义的消息类型：

- `AgentHandler` - 使用 `pb.Agent`, `pb.CreateAgentRequest` 等
- `ConversationHandler` - 使用 `pb.Conversation`, `pb.Message` 等
- `ToolHandler` - 使用 `pb.Tool`, `pb.CreateToolRequest` 等
- `KnowledgeBaseHandler` - 使用 `pb.KnowledgeBase`, `pb.Document` 等

### 5. 数据库架构简化

**移除了 MongoDB 依赖：**

- 从 `internal/config/config.go` 中移除 `MongoDBConfig` 结构体
- 从 `.env` 和 `.env.example` 中移除 MongoDB 相关配置项
- 简化了数据存储架构，使用单一的 MySQL 数据库

**当前数据库配置：**

- **MySQL** - 主数据库（用于结构化数据存储）
- **Redis** - 缓存和会话存储
- **Milvus** - 向量数据库（用于知识库和语义搜索）

这种架构更加简洁，减少了运维复杂度，同时保持了系统的核心功能。

## 开发流程

### 编译 Proto 文件

```bash
cd backend
make proto
```

这会在 `api/proto/gen/` 目录下生成 Go 代码。

### 清理生成的文件

```bash
make proto-clean
```

### 构建项目

```bash
make build
```

或者

```bash
go build -o bin/agent-platform ./cmd/server
```

### 运行项目

```bash
make run
```

或者

```bash
go run ./cmd/server/main.go
```

### 完整开发流程

```bash
# 生成 proto、构建并运行
make dev
```

## 项目结构

```
backend/
├── api/
│   └── proto/              # Proto 定义文件
│       ├── common.proto
│       ├── agent.proto
│       ├── conversation.proto
│       ├── tool.proto
│       ├── knowledge_base.proto
│       └── gen/            # 生成的 Go 代码
├── cmd/
│   └── server/
│       ├── main.go         # 主入口
│       └── router.go       # 路由配置
├── internal/
│   ├── handler/            # HTTP 处理器（使用 proto 定义）
│   │   ├── agent_handler.go
│   │   ├── conversation_handler.go
│   │   ├── tool_handler.go
│   │   └── knowledge_base_handler.go
│   ├── middleware/         # 中间件
│   ├── model/schema/       # 数据模型
│   ├── repository/         # 数据访问层（待实现）
│   └── service/            # 业务逻辑层（待实现）
├── pkg/
│   └── response/           # 统一响应工具包
│       └── response.go
├── Makefile               # 构建脚本
├── go.mod
└── go.sum
```

## API 示例

### 创建 Agent

**请求：**
```bash
POST /api/v1/agents
Content-Type: application/json

{
  "name": "Customer Service Agent",
  "description": "AI agent for customer service",
  "type": "single",
  "model_config": {
    "model": "gpt-4",
    "temperature": 0.7
  },
  "tools": ["search", "email"],
  "knowledge_bases": ["kb-001"],
  "prompt_template": "You are a helpful customer service agent."
}
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Customer Service Agent",
    "description": "AI agent for customer service",
    "type": "single",
    "status": "draft",
    "version": "1.0.0",
    "created_by": "system",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "model_config": {
      "model": "gpt-4",
      "temperature": 0.7
    },
    "tools": ["search", "email"],
    "knowledge_bases": ["kb-001"],
    "prompt_template": "You are a helpful customer service agent."
  }
}
```

### 获取 Agent 列表

**请求：**
```bash
GET /api/v1/agents?page=1&page_size=10&status=published
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "Customer Service Agent",
        "description": "AI agent for customer service",
        "type": "single",
        "status": "published",
        "version": "1.0.0",
        "created_by": "admin",
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
      }
    ],
    "page": 1,
    "page_size": 10,
    "total": 1
  }
}
```

## 添加新接口的步骤

1. **在对应的 `.proto` 文件中定义消息类型**
   ```protobuf
   message NewFeatureRequest {
     string name = 1;
     string description = 2;
   }

   message NewFeature {
     string id = 1;
     string name = 2;
     string description = 3;
   }
   ```

2. **生成 Go 代码**
   ```bash
   make proto
   ```

3. **在 Handler 中实现接口**
   ```go
   func (h *Handler) CreateNewFeature(c *gin.Context) {
       var req pb.NewFeatureRequest
       if err := c.ShouldBindJSON(&req); err != nil {
           response.BadRequest(c, err.Error())
           return
       }

       // 业务逻辑...

       response.Created(c, data)
   }
   ```

4. **在 router.go 中注册路由**
   ```go
   api.POST("/features", handler.CreateNewFeature)
   ```

## 注意事项

1. **所有接口都使用 proto 定义的消息类型**
2. **所有响应都使用 `pkg/response` 包中的工具函数**
3. **修改 proto 文件后必须重新运行 `make proto`**
4. **保持响应格式的一致性：`{code, message, data}`**
5. **错误处理统一使用 response 包提供的错误函数**

## 工具安装

如果需要安装 protoc 编译器插件：

```bash
make install-tools
```

这会安装：
- `protoc-gen-go` - Proto 到 Go 的编译器插件
- `protoc-gen-go-grpc` - gRPC 相关的编译器插件

## 后续优化建议

1. **实现 Service 层** - 将业务逻辑从 Handler 中抽离
2. **实现 Repository 层** - 实现真实的数据库访问
3. **添加单元测试** - 为 Handler 和 Service 添加测试
4. **API 文档生成** - 使用 proto 定义生成 API 文档
5. **错误码规范** - 定义统一的错误码体系
6. **日志规范** - 统一日志格式和级别
7. **性能优化** - 减少 proto 到 map 的转换开销

## 相关文档

- [Protocol Buffers 文档](https://protobuf.dev/)
- [Gin Web Framework](https://gin-gonic.com/)
- [Go gRPC 文档](https://grpc.io/docs/languages/go/)
