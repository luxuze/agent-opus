# HTTP REST API Gateway 配置说明

## 概述

项目现在同时支持 **gRPC** 和 **HTTP REST API** 两种调用方式，通过 gRPC-Gateway 实现自动转换。

## 架构

```
┌──────────────────────────────────────────┐
│           客户端选择                      │
│                                          │
│  ┌──────────────┐    ┌────────────────┐ │
│  │ gRPC Client  │    │  HTTP Client   │ │
│  │  (Protobuf)  │    │    (JSON)      │ │
│  └──────┬───────┘    └────────┬───────┘ │
└─────────┼──────────────────────┼─────────┘
          │                      │
          │                      │
          ▼                      ▼
     Port 9000             Port 8000
          │                      │
          │                      │
          ▼                      ▼
    ┌─────────┐          ┌──────────────┐
    │  gRPC   │◄─────────│ gRPC-Gateway │
    │ Server  │          │  (HTTP→gRPC) │
    └─────────┘          └──────────────┘
```

## 端口分配

| 服务         | 端口 | 环境变量  | 协议       | 用途                 |
| ------------ | ---- | --------- | ---------- | -------------------- |
| gRPC Server  | 9000 | GRPC_PORT | gRPC/HTTP2 | 直接 gRPC 调用       |
| HTTP Gateway | 8000 | HTTP_PORT | HTTP/1.1   | REST API (JSON 格式) |

### 配置端口

端口可通过环境变量配置（`.env` 文件或环境变量）：

```bash
# .env 文件
GRPC_PORT=9000          # gRPC 服务端口
HTTP_PORT=8000           # HTTP REST API 端口
SERVER_HOST=0.0.0.0      # 监听地址
```

或通过环境变量：

```bash
export GRPC_PORT=9000
export HTTP_PORT=8000
go run ./cmd/server/main.go
```

## 统一的 API 路径

所有 HTTP REST API 遵循统一规范：`/api/v1/{resource}`

### Agent API

```
POST   /api/v1/agents           创建 Agent
GET    /api/v1/agents           获取列表（支持分页）
GET    /api/v1/agents/{id}      获取详情
PUT    /api/v1/agents/{id}      更新
DELETE /api/v1/agents/{id}      删除
```

### Conversation API

```
POST   /api/v1/conversations                          创建对话
GET    /api/v1/conversations                          获取列表
GET    /api/v1/conversations/{id}                     获取详情
POST   /api/v1/conversations/{conversation_id}/messages   发送消息
```

### Tool API

```
POST   /api/v1/tools            创建工具
GET    /api/v1/tools            获取列表
GET    /api/v1/tools/{id}       获取详情
DELETE /api/v1/tools/{id}       删除
```

### KnowledgeBase API

```
POST   /api/v1/knowledge-bases                              创建知识库
GET    /api/v1/knowledge-bases                              获取列表
GET    /api/v1/knowledge-bases/{id}                         获取详情
POST   /api/v1/knowledge-bases/{knowledge_base_id}/documents   上传文档
DELETE /api/v1/knowledge-bases/{id}                         删除
```

## 快速测试

### 启动服务

```bash
cd backend
go run ./cmd/server/main.go
```

服务启动后会看到：

```
INFO    Starting Agent Platform gRPC Server    {"port": "9000", "mode": "debug"}
INFO    gRPC Server listening   {"address": "0.0.0.0:9000"}
INFO    HTTP Gateway listening  {"port": "8000"}
```

### 测试 HTTP REST API

#### 创建 Agent

```bash
curl -X POST http://localhost:8000/api/v1/agents \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Agent",
    "description": "A test agent",
    "type": "single"
  }'
```

#### 获取 Agent 列表

```bash
curl http://localhost:8000/api/v1/agents
```

#### 获取指定 Agent

```bash
curl http://localhost:8000/api/v1/agents/{agent-id}
```

### 测试 gRPC（使用 grpcurl）

```bash
# 创建 Agent
grpcurl -plaintext -d '{
  "name": "Test Agent",
  "description": "A test agent",
  "type": "single"
}' localhost:9000 api.AgentService/CreateAgent

# 获取列表
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10
}' localhost:9000 api.AgentService/ListAgents
```

## 前端集成示例

### 使用 Axios (HTTP REST API)

```typescript
import axios from "axios";

const apiClient = axios.create({
  baseURL: "http://localhost:8000/api/v1",
  headers: {
    "Content-Type": "application/json",
  },
});

// 创建 Agent
export const createAgent = async (data) => {
  const response = await apiClient.post("/agents", data);
  return response.data;
};

// 获取列表
export const listAgents = async (page = 1, pageSize = 10) => {
  const response = await apiClient.get("/agents", {
    params: { page, page_size: pageSize },
  });
  return response.data;
};
```

### 使用 gRPC-Web (Protobuf)

```typescript
import { AgentServiceClient } from "./proto/AgentServiceClientPb";
import { CreateAgentRequest } from "./proto/agent_pb";

const client = new AgentServiceClient("http://localhost:8000");

const request = new CreateAgentRequest();
request.setName("Test Agent");
request.setDescription("A test agent");
request.setType("single");

client.createAgent(request, {}, (err, response) => {
  if (err) {
    console.error("Error:", err);
    return;
  }
  console.log("Created:", response.toObject());
});
```

## 实现细节

### Proto 定义

在 proto 文件中使用 Google API 注解定义 HTTP 路由：

```protobuf
import "google/api/annotations.proto";

service AgentService {
  rpc CreateAgent(CreateAgentRequest) returns (Agent) {
    option (google.api.http) = {
      post: "/api/v1/agents"
      body: "*"
    };
  }

  rpc GetAgent(GetAgentRequest) returns (Agent) {
    option (google.api.http) = {
      get: "/api/v1/agents/{id}"
    };
  }
}
```

### 代码生成

Makefile 已更新，自动生成 gRPC-Gateway 代码：

```bash
make proto-backend
```

生成的文件：

- `*.pb.go` - Protobuf 消息定义
- `*_grpc.pb.go` - gRPC 服务端代码
- `*.pb.gw.go` - gRPC-Gateway HTTP 处理代码 ✨ 新增

### 后端集成

`cmd/server/gateway.go` 提供 HTTP Gateway 服务：

```go
func setupGateway(grpcAddress string, httpPort string, logger *zap.Logger) error {
    mux := runtime.NewServeMux()

    // 注册所有服务
    pb.RegisterAgentServiceHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
    pb.RegisterConversationServiceHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
    // ...

    // 启动 HTTP 服务器
    http.ListenAndServe(httpAddr, handler)
}
```

## 对比选择

| 特性         | gRPC               | HTTP REST              |
| ------------ | ------------------ | ---------------------- |
| 传输协议     | HTTP/2             | HTTP/1.1               |
| 数据格式     | Protobuf (二进制)  | JSON (文本)            |
| 性能         | ⚡ 更快            | 🐢 较慢                |
| 浏览器支持   | 需要 gRPC-Web      | ✅ 原生支持            |
| 开发调试     | 需要特殊工具       | ✅ curl/Postman        |
| 类型安全     | ✅ 强类型          | ⚠️ 弱类型              |
| 流式传输     | ✅ 支持            | ❌ 不支持              |
| **推荐场景** | **微服务、高性能** | **公共 API、快速集成** |

## 开发工具推荐

### HTTP REST API

- **curl** - 命令行测试
- **Postman** - 图形界面测试
- **HTTPie** - 友好的命令行工具
- **Swagger/OpenAPI** - API 文档（可选配置）

### gRPC

- **grpcurl** - 命令行测试工具
- **BloomRPC** - 图形界面测试工具
- **Postman** - 支持 gRPC 测试

## 生产环境配置

### 添加认证

```go
// 在 gateway.go 中添加认证中间件
func authMiddleware(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if !validateToken(token) {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        h.ServeHTTP(w, r)
    })
}
```

### CORS 配置

当前已在 `gateway.go` 中配置了基本的 CORS 支持，生产环境建议：

```go
func cors(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 限制允许的来源
        origin := r.Header.Get("Origin")
        if isAllowedOrigin(origin) {
            w.Header().Set("Access-Control-Allow-Origin", origin)
        }
        // ... 其他 CORS 配置
    })
}
```

### 限流和监控

建议添加：

- 速率限制（rate limiting）
- 请求日志
- 性能监控
- 错误追踪

## 相关文档

- **[HTTP API 详细文档](./proto/HTTP_API.md)** - 完整的 REST API 使用说明
- **[gRPC 文档](./backend/GRPC.md)** - gRPC 服务架构说明
- **[Proto 文件说明](./proto/README.md)** - Protocol Buffers 定义
- **[gRPC-Gateway 官方文档](https://grpc-ecosystem.github.io/grpc-gateway/)** - 官方参考

## 故障排查

### HTTP Gateway 启动失败

```bash
# 检查端口是否被占用
lsof -i :8000

# 检查 gRPC 服务是否正常
grpcurl -plaintext localhost:9000 list
```

### 跨域错误

确保在 `gateway.go` 中正确配置了 CORS：

```go
w.Header().Set("Access-Control-Allow-Origin", "*")
```

### JSON 格式问题

gRPC-Gateway 自动处理 Protobuf 到 JSON 的转换，字段名会自动转换为 snake_case。

## 总结

✅ **已完成：**

- 添加 HTTP REST API 支持（gRPC-Gateway）
- 统一的 API 路径规范 (`/api/v1/*`)
- 前后端共享 proto 定义
- 自动代码生成
- 完整的文档和示例

🎯 **建议：**

- 根据实际需求选择 gRPC 或 HTTP REST API
- 内部服务优先使用 gRPC（高性能）
- 公共 API 和前端优先使用 HTTP REST（易用性）
- 两种方式可以同时使用，互不影响
