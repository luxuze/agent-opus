# 后端重构说明文档

## 重构概述

本次重构主要完成了以下工作：

1. **统一使用 Protocol Buffers 定义接口**
2. **改造为纯 gRPC 服务架构**
3. **配置 Envoy 作为 gRPC-Web 代理**
4. **建立标准化的项目结构和开发流程**
5. **移除 MongoDB 依赖，简化数据库架构**
6. **将 Proto 文件移至根目录，前后端共享接口定义**
7. **✨ 新增：集成 gRPC-Gateway，同时支持 HTTP REST API**

## 主要改动

### 1. gRPC 服务架构

**从 HTTP REST 改造为纯 gRPC 服务：**

```
Frontend (Browser) → Envoy Proxy (gRPC-Web) → Backend (gRPC)
    :5173               :8000                    :9000
```

**关键组件：**

- **Backend**: 纯 gRPC 服务器，监听端口 9000
- **Envoy**: gRPC-Web 代理，将 HTTP/1.1 转换为 gRPC HTTP/2
- **Frontend**: 使用 gRPC-Web 客户端库调用服务

### 2. Proto 定义（项目根目录 `proto/`）

**重要变更：Proto 文件现已移至项目根目录，实现前后端共享**

创建了以下 proto 文件（包含 message 和 service 定义）：

- **`common.proto`** - 通用响应结构和分页定义
- **`agent.proto`** - Agent 服务（AgentService）及消息定义
- **`conversation.proto`** - Conversation 服务（ConversationService）及消息定义
- **`tool.proto`** - Tool 服务（ToolService）及消息定义
- **`knowledge_base.proto`** - KnowledgeBase 服务（KnowledgeBaseService）及消息定义

**每个 proto 文件包含：**

- Message 定义（请求/响应）
- Service 定义（gRPC 方法）

**生成的代码位置：**

- **后端 Go 代码**: `backend/gen/go/` - 使用 protoc-gen-go 和 protoc-gen-go-grpc 生成
- **前端 TypeScript 代码**: `frontend/src/proto/` - 使用 protoc-gen-grpc-web 生成

### 3. gRPC 服务实现

**新增 gRPC 服务实现（`internal/grpc/`）：**

- `agent_service.go` - AgentService 实现
- `conversation_service.go` - ConversationService 实现
- `tool_service.go` - ToolService 实现
- `knowledge_base_service.go` - KnowledgeBaseService 实现

**特性：**

- 使用 Protobuf 进行序列化/反序列化
- 支持 gRPC Reflection（便于调试）
- 统一的错误处理（gRPC status codes）
- 类型安全的 API

### 4. Envoy gRPC-Web 代理配置

**配置文件：** `backend/envoy.yaml`

**功能：**

- gRPC-Web 到 gRPC 协议转换
- 自动处理 CORS
- HTTP/1.1 到 HTTP/2 转换
- 路由配置和负载均衡

**端口：**

- 8000: gRPC-Web HTTP 接口（对外）
- 9901: 管理和监控接口

### 5. 服务启动方式变更

**原架构（HTTP REST）：**

```go
router := gin.Default()
router.Run(":8000")
```

**新架构（gRPC）：**

```go
grpcServer := grpc.NewServer()
pb.RegisterAgentServiceServer(grpcServer, agentService)
// ...
grpcServer.Serve(lis)
```

**启动文件：** `cmd/server/main.go` - 完全重写为启动 gRPC 服务器

### 6. 数据库架构简化

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

**从项目根目录执行（推荐）：**

```bash
# 生成前后端所有 proto 代码
make proto

# 仅生成后端 Go 代码
make proto-backend

# 仅生成前端 TypeScript 代码
make proto-frontend
```

**或者在后端目录单独生成：**

```bash
cd backend
make proto
```

**生成的文件位置：**

- 后端 Go 代码：`backend/gen/go/`
  - `*.pb.go` - Protobuf 消息定义
  - `*_grpc.pb.go` - gRPC 服务接口和桩代码
- 前端 TypeScript 代码：`frontend/src/proto/`
  - `*_pb.js` - Protobuf 消息定义（JavaScript）
  - `*_grpc_web_pb.js` - gRPC-Web 客户端代码

### 清理生成的文件

**从根目录清理所有生成的文件：**

```bash
make proto-clean
```

**或在后端目录单独清理：**

```bash
cd backend
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
agent-opus/                 # 项目根目录
├── proto/                  # 共享的 Proto 定义文件（前后端共享）
│   ├── common.proto
│   ├── agent.proto
│   ├── conversation.proto
│   ├── tool.proto
│   └── knowledge_base.proto
├── backend/
│   ├── gen/
│   │   └── go/             # 生成的 Go 代码
│   │       ├── *.pb.go          # Protobuf 消息
│   │       └── *_grpc.pb.go     # gRPC 服务
│   ├── cmd/
│   │   └── server/
│   │       └── main.go     # gRPC 服务器主入口
│   ├── internal/
│   │   ├── grpc/           # gRPC 服务实现
│   │   │   ├── agent_service.go
│   │   │   ├── conversation_service.go
│   │   │   ├── tool_service.go
│   │   │   └── knowledge_base_service.go
│   │   ├── config/         # 配置管理
│   │   ├── model/schema/   # 数据模型
│   │   ├── repository/     # 数据访问层（待实现）
│   │   └── service/        # 业务逻辑层（待实现）
│   ├── envoy.yaml         # Envoy 代理配置
│   ├── Makefile           # 后端构建脚本
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   └── proto/          # 生成的 TypeScript/JavaScript 代码
│   │       ├── *_pb.js          # Protobuf 消息
│   │       └── *_grpc_web_pb.js # gRPC-Web 客户端
│   ├── package.json
│   └── ...
├── Makefile               # 根目录构建脚本（统一 proto 生成）
└── docker-compose.yml
```

## gRPC 服务使用示例

### 使用 grpcurl 测试

**列出所有服务：**

```bash
grpcurl -plaintext localhost:9000 list
```

**创建 Agent：**

```bash
grpcurl -plaintext -d '{
  "name": "Customer Service Agent",
  "description": "AI agent for customer service",
  "type": "single"
}' localhost:9000 api.AgentService/CreateAgent
```

**获取 Agent 列表：**

```bash
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10
}' localhost:9000 api.AgentService/ListAgents
```

### 前端 gRPC-Web 调用示例

```typescript
import { AgentServiceClient } from "./proto/agent_grpc_web_pb";
import { CreateAgentRequest } from "./proto/agent_pb";

// 创建客户端（连接到 Envoy 代理）
const client = new AgentServiceClient("http://localhost:8000");

// 创建请求
const request = new CreateAgentRequest();
request.setName("Customer Service Agent");
request.setDescription("AI agent for customer service");
request.setType("single");

// 调用 gRPC 方法
client.createAgent(request, {}, (err, response) => {
  if (err) {
    console.error("Error:", err);
    return;
  }

  console.log("Created agent:", response.toObject());
});
```

## 添加新 gRPC 方法的步骤

1. **在对应的 `.proto` 文件中定义消息和 RPC 方法**

   ```protobuf
   // 定义消息
   message NewFeatureRequest {
     string name = 1;
     string description = 2;
   }

   message NewFeature {
     string id = 1;
     string name = 2;
     string description = 3;
   }

   // 在 service 中添加 RPC 方法
   service FeatureService {
     rpc CreateFeature(NewFeatureRequest) returns (NewFeature);
   }
   ```

2. **生成 Go 代码**

   ```bash
   make proto
   ```

3. **实现 gRPC 服务方法**

   ```go
   func (s *FeatureServer) CreateFeature(ctx context.Context, req *pb.NewFeatureRequest) (*pb.NewFeature, error) {
       if req.Name == "" {
           return nil, status.Error(codes.InvalidArgument, "name is required")
       }

       feature := &pb.NewFeature{
           Id:          uuid.New().String(),
           Name:        req.Name,
           Description: req.Description,
       }

       return feature, nil
   }
   ```

4. **在 main.go 中注册服务**
   ```go
   pb.RegisterFeatureServiceServer(grpcServer, grpcserver.NewFeatureServer())
   ```

## 注意事项

1. **所有接口都使用 proto 定义的消息类型**
2. **使用 gRPC status codes 进行错误处理**
3. **修改 proto 文件后必须重新运行 `make proto`**
4. **gRPC 服务默认端口为 9000（可配置）**
5. **前端通过 Envoy 代理（端口 8000）访问服务**
6. **使用 grpcurl 或 Postman 测试 gRPC 服务**

## 工具安装

### 从项目根目录一键安装（推荐）

```bash
cd ..  # 回到项目根目录
make install-proto-tools
```

这会自动安装所有必需的工具：

- **后端工具**:
  - `protoc-gen-go` - Proto 到 Go 的编译器插件
  - `protoc-gen-go-grpc` - gRPC 相关的编译器插件
- **前端工具**:
  - `protoc-gen-grpc-web` - gRPC-Web 客户端生成插件
  - `protoc-gen-js` - JavaScript/TypeScript 代码生成插件

### 仅安装后端工具

```bash
make install-tools
```

### 手动安装

**后端:**

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

**前端 (macOS):**

```bash
brew install protoc-gen-grpc-web protoc-gen-js
```

更多安装选项请参考 `proto/README.md`

## 后续优化建议

1. **实现 Service 层** - 将业务逻辑从 Handler 中抽离
2. **实现 Repository 层** - 实现真实的数据库访问
3. **添加单元测试** - 为 Handler 和 Service 添加测试
4. **API 文档生成** - 使用 proto 定义生成 API 文档
5. **错误码规范** - 定义统一的错误码体系
6. **日志规范** - 统一日志格式和级别
7. **性能优化** - 减少 proto 到 map 的转换开销

## 相关文档

- [gRPC 服务架构详细说明](./GRPC.md)
- [数据库配置指南](./DATABASE.md)
- [Protocol Buffers 文档](https://protobuf.dev/)
- [Go gRPC 文档](https://grpc.io/docs/languages/go/)
- [gRPC-Web](https://github.com/grpc/grpc-web)
- [Envoy Proxy](https://www.envoyproxy.io/)
