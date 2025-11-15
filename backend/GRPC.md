# gRPC 服务架构说明

## 架构概述

后端已完全改造为**纯 gRPC 服务**，前端通过 **Envoy gRPC-Web 代理**调用服务。

```
┌─────────────┐      gRPC-Web       ┌────────────┐      gRPC       ┌─────────────┐
│   Frontend  │ ───────────────────> │   Envoy    │ ──────────────> │   Backend   │
│   (Browser) │  HTTP/1.1 + Proto   │   Proxy    │  HTTP/2 + Proto │  (gRPC)     │
└─────────────┘                      └────────────┘                 └─────────────┘
    :5173                                :80                          :9000
```

### 组件说明

1. **Frontend (浏览器)**

   - 使用 gRPC-Web 客户端库
   - 通过 HTTP/1.1 与 Envoy 通信
   - 使用 Protobuf 序列化数据

2. **Envoy Proxy**

   - 作为 gRPC-Web 到 gRPC 的转换层
   - 监听端口 8000（HTTP/gRPC-Web）
   - 转发请求到后端 gRPC 服务（端口 9000）
   - 自动处理 CORS
   - 管理端口 9901（可查看配置和统计）

3. **Backend (gRPC Server)**
   - 纯 gRPC 服务器
   - 监听端口 9000
   - 实现所有业务逻辑
   - 支持 gRPC Reflection（便于调试）

## gRPC 服务列表

### 1. AgentService

管理 AI Agent 的 CRUD 操作。

**方法：**

- `CreateAgent` - 创建 Agent
- `ListAgents` - 获取 Agent 列表（支持分页和筛选）
- `GetAgent` - 获取 Agent 详情
- `UpdateAgent` - 更新 Agent
- `DeleteAgent` - 删除 Agent

### 2. ConversationService

管理对话和消息。

**方法：**

- `CreateConversation` - 创建对话
- `GetConversation` - 获取对话详情
- `ListConversations` - 获取对话列表
- `SendMessage` - 发送消息并获取 AI 响应

### 3. ToolService

管理工具。

**方法：**

- `CreateTool` - 创建工具
- `ListTools` - 获取工具列表
- `GetTool` - 获取工具详情
- `DeleteTool` - 删除工具

### 4. KnowledgeBaseService

管理知识库和文档。

**方法：**

- `CreateKnowledgeBase` - 创建知识库
- `ListKnowledgeBases` - 获取知识库列表
- `GetKnowledgeBase` - 获取知识库详情
- `UploadDocument` - 上传文档到知识库
- `DeleteKnowledgeBase` - 删除知识库

## 开发流程

### 修改 Proto 定义

**重要：Proto 文件现已统一放在项目根目录 `proto/` 下，前后端共享**

1. 编辑 proto 文件（项目根目录 `proto/*.proto`）
2. 重新生成代码：

   ```bash
   # 从项目根目录
   make proto              # 生成前后端所有代码
   make proto-backend      # 仅生成后端 Go 代码
   make proto-frontend     # 仅生成前端 TypeScript 代码

   # 或从后端目录
   cd backend && make proto
   ```

### 添加新的 gRPC 方法

1. 在项目根目录的 proto 文件中定义新的 RPC 方法（`proto/*.proto`）
2. 运行 `make proto` 生成前后端代码
3. 在对应的 service 文件中实现方法（`backend/internal/grpc/*_service.go`）
4. 在前端使用生成的客户端代码调用新方法

### 本地运行

```bash
# 从项目根目录
make proto              # 生成 proto 代码

# 进入后端目录
cd backend

# 编译
make build

# 运行 gRPC 服务器
./bin/agent-platform
```

服务将在 `0.0.0.0:8000`（配置的 SERVER_PORT）启动。

## Docker 部署

### 构建并启动所有服务

```bash
docker-compose up -d
```

### 服务端口

- **Backend gRPC**: `9000` (仅容器内部访问)
- **Envoy Proxy**: `8000` (gRPC-Web HTTP 接口)
- **Envoy Admin**: `9901` (管理和监控)
- **Frontend**: `5173`

### 查看日志

```bash
# 后端 gRPC 服务日志
docker logs -f agent-platform-backend

# Envoy 代理日志
docker logs -f agent-platform-envoy

# 前端服务日志
docker logs -f agent-platform-frontend
```

## 测试 gRPC 服务

### 使用 grpcurl

grpcurl 是一个命令行工具，类似于 curl 但用于 gRPC。

**安装：**

```bash
brew install grpcurl
```

**列出所有服务：**

```bash
grpcurl -plaintext localhost:9000 list
```

**列出服务的方法：**

```bash
grpcurl -plaintext localhost:9000 list api.AgentService
```

**调用方法：**

```bash
# 创建 Agent
grpcurl -plaintext -d '{
  "name": "Test Agent",
  "description": "A test agent",
  "type": "single"
}' localhost:9000 api.AgentService/CreateAgent

# 获取 Agent 列表
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10
}' localhost:9000 api.AgentService/ListAgents

# 获取 Agent 详情
grpcurl -plaintext -d '{
  "id": "agent-id-here"
}' localhost:9000 api.AgentService/GetAgent
```

### 使用 Postman

Postman 支持 gRPC 调用：

1. 打开 Postman
2. 创建新的 gRPC 请求
3. 输入服务地址：`localhost:9000`
4. 选择方法并发送请求

### 通过 Envoy 测试（gRPC-Web）

使用浏览器或任何 HTTP 客户端通过 Envoy 代理访问：

```bash
# 注意：需要使用 gRPC-Web 格式的请求
curl -X POST http://localhost:8000/api.AgentService/ListAgents \
  -H "Content-Type: application/grpc-web+proto" \
  --data-binary @request.bin
```

## 前端集成

### 安装 gRPC-Web 客户端

```bash
cd frontend
npm install grpc-web
npm install google-protobuf
```

### 生成前端代码

**从项目根目录（推荐）：**

```bash
# 生成前端 TypeScript/JavaScript 代码
make proto-frontend
```

**或手动生成：**

```bash
# 使用 protoc-gen-grpc-web 插件
protoc -I=./proto \
  --js_out=import_style=commonjs:./frontend/src/proto \
  --grpc-web_out=import_style=typescript,mode=grpcwebtext:./frontend/src/proto \
  ./proto/*.proto
```

**或使用 npm script：**

```bash
cd frontend
npm run proto
```

### 使用示例

```typescript
import { AgentServiceClient } from "./proto/agent_grpc_web_pb";
import { ListAgentsRequest } from "./proto/agent_pb";

// 创建客户端（指向 Envoy 代理）
const client = new AgentServiceClient("http://localhost:8000");

// 调用 gRPC 方法
const request = new ListAgentsRequest();
request.setPage(1);
request.setPageSize(10);

client.listAgents(request, {}, (err, response) => {
  if (err) {
    console.error("Error:", err);
    return;
  }

  console.log("Agents:", response.getItemsList());
});
```

## Envoy 配置

Envoy 配置文件位于 `backend/envoy.yaml`。

### 重要配置项

1. **CORS 配置**

   - 允许所有来源（生产环境需要限制）
   - 支持必要的 gRPC-Web 头

2. **gRPC-Web 过滤器**

   - 自动转换 gRPC-Web 到 gRPC

3. **超时设置**
   - 路由超时设置为 0（无超时限制）
   - 根据实际需求调整

### 查看 Envoy 统计

```bash
curl http://localhost:9901/stats
```

### 查看 Envoy 配置

```bash
curl http://localhost:9901/config_dump
```

## 性能优化

### gRPC 优势

1. **HTTP/2 多路复用**

   - 单个连接上并发多个请求
   - 减少连接开销

2. **二进制协议**

   - Protobuf 序列化比 JSON 更高效
   - 更小的消息大小

3. **流式传输**
   - 支持服务器流、客户端流、双向流
   - 适合实时数据传输

### 优化建议

1. **连接池**

   - 复用 gRPC 连接
   - 避免频繁创建销毁连接

2. **压缩**

   - 启用 gRPC 消息压缩
   - 减少网络传输数据量

3. **并发控制**
   - 合理设置最大并发请求数
   - 避免服务器过载

## 故障排查

### 常见问题

1. **连接被拒绝**

   ```
   Error: failed to connect to all addresses
   ```

   - 检查后端 gRPC 服务是否正常运行
   - 验证端口号是否正确

2. **Envoy 转发失败**

   ```
   Error: upstream connect error or disconnect/reset before headers
   ```

   - 检查 Envoy 配置中的 backend 地址
   - 确认 backend 容器名和端口正确

3. **CORS 错误**
   ```
   Access to fetch at ... has been blocked by CORS policy
   ```
   - 检查 Envoy 的 CORS 配置
   - 确认允许的来源包含前端地址

### 调试技巧

1. **启用详细日志**

   ```go
   grpcServer := grpc.NewServer(
       grpc.UnaryInterceptor(loggingInterceptor),
   )
   ```

2. **使用 gRPC 健康检查**

   ```bash
   grpcurl -plaintext localhost:9000 grpc.health.v1.Health/Check
   ```

3. **查看 Envoy 日志**
   ```bash
   docker logs -f agent-platform-envoy
   ```

## 安全考虑

### 生产环境配置

1. **TLS 加密**

   - 启用 gRPC TLS
   - 配置 Envoy TLS 终止

2. **认证授权**

   - 实现 JWT 认证拦截器
   - 在 gRPC metadata 中传递 token

3. **速率限制**

   - 配置 Envoy 速率限制
   - 防止 DDoS 攻击

4. **CORS 限制**
   - 生产环境中限制允许的来源
   - 不要使用通配符 `*`

## 相关资源

- [gRPC 官方文档](https://grpc.io/)
- [gRPC-Web 文档](https://github.com/grpc/grpc-web)
- [Envoy 代理文档](https://www.envoyproxy.io/docs/envoy/latest/)
- [Protocol Buffers](https://protobuf.dev/)
- [grpcurl](https://github.com/fullstorydev/grpcurl)

## 后续计划

- [ ] 实现 Service 层和 Repository 层
- [ ] 添加 JWT 认证中间件
- [ ] 实现流式 RPC（用于实时对话）
- [ ] 添加 gRPC 健康检查
- [ ] 实现分布式追踪（OpenTelemetry）
- [ ] 添加 Prometheus 指标导出
