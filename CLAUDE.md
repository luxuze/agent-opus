# CLAUDE.md

Claude Code 工作指南 - Agent Platform 智能代理平台

## 项目概述

Agent Platform 是一个智能代理管理和编排系统，采用前后端分离架构：

**技术栈**

- 后端：Go + gRPC + HTTP Gateway (grpc-gateway)
- 前端：React 18 + TypeScript + Vite + Ant Design
- 数据库：PostgreSQL (主库 + pgvector 向量扩展)、Redis (缓存)
- AI：OpenAI API、SiliconFlow DeepSeek API

**核心功能**

- Agent 管理（创建、配置、执行智能代理）
- 对话管理（多轮对话、上下文记忆）
- 工具集成（代理可调用的外部工具）
- 知识库（文档上传、向量检索、RAG）
- 用户认证（JWT）

---

## 快速开始

### 本地开发（推荐）

**启动后端**

```bash
cd backend
go run cmd/server/main.go
# gRPC: http://localhost:9000
# HTTP REST API: http://localhost:8000
```

**启动前端**

```bash
cd frontend
npm install
npm run dev
# 访问 http://localhost:5173
```

**测试 API**

```bash
# 健康检查
curl http://localhost:8000/health

# 查看 API 列表
./test-api.sh
```

### Docker 部署

**配置环境变量**（在项目根目录创建 `.env` 文件）：
```bash
# 项目根目录的 .env 文件（供 docker-compose 使用）
SILICONFLOW_API_KEY=sk-your-key-here
OPENAI_API_KEY=sk-your-key-here  # 可选
```

**启动服务**：
```bash
# 启动所有服务
make start
# 或
docker-compose up -d

# 查看日志（检查 AI 服务是否初始化成功）
make logs
# 或
docker-compose logs -f backend

# 停止服务
make stop
```

---

## 架构设计

### 后端架构：gRPC + HTTP Gateway 双协议

```
┌─────────────┐
│   Frontend  │
│  (React)    │
└──────┬──────┘
       │ HTTP REST
       ▼
┌─────────────────────┐
│   HTTP Gateway      │ ← grpc-gateway 自动生成 REST API
│   :8000             │
└──────┬──────────────┘
       │ 内部调用
       ▼
┌─────────────────────┐
│   gRPC Services     │ ← 核心业务逻辑
│   :9000             │
└──────┬──────────────┘
       │
       ▼
┌─────────────────────┐
│   Database Layer    │
│  (PostgreSQL/Redis) │
└─────────────────────┘
```

**关键实现**

- `proto/*.proto` - 服务定义（单一数据源）
- `backend/internal/grpc/*_service.go` - gRPC 服务实现
- `backend/cmd/server/gateway.go` - HTTP Gateway 配置
- `backend/cmd/server/marshaler.go` - 统一响应格式 + snake_case 命名

**响应格式**（所有 HTTP API）

```json
{
  "code": 0,
  "message": "success",
  "data": {
    /* 实际数据，使用 snake_case 字段命名 */
    "conversation_id": "xxx",
    "agent_id": "yyy",
    "created_at": "2025-01-01T00:00:00Z"
  },
  "timestamp": 1634567890,
  "request_id": "uuid"
}
```

**重要**：所有 API 响应使用 **snake_case** 命名（与 proto 定义一致），不是 camelCase

### 前端架构

```
src/
├── pages/          # 页面组件（Agent, Conversation, Tool, KnowledgeBase）
├── components/     # 公共组件
├── services/       # API 调用层（使用统一的 api.ts 客户端）
├── store/          # Redux 状态管理
└── types/          # TypeScript 类型定义
```

**API 调用规范**

- 使用 `src/services/api.ts` 配置的 axios 实例
- baseURL: `/api/v1` (通过 Vite proxy 转发到后端)
- 自动添加 JWT token (从 localStorage)
- 统一错误处理（401 自动跳转登录）

---

## 常用命令

### 后端开发

```bash
# 运行
cd backend
go run cmd/server/main.go

# 构建
go build -o server ./cmd/server
# 或
make backend

# 测试
go test ./...

# 格式化
go fmt ./...
```

### 前端开发

```bash
cd frontend

# 安装依赖
npm install

# 开发模式
npm run dev

# 构建生产版本
npm run build
# 或
make frontend

# 代码检查
npm run lint
```

### Proto 代码生成

**何时需要生成 proto 代码？**
修改了 `proto/*.proto` 文件后必须执行：

```bash
# 生成 Go 代码（后端必须）
make proto-backend

# 生成 TypeScript 代码（前端可选）
make proto-frontend

# 同时生成
make proto

# 清理生成的文件
make proto-clean
```

**首次使用需安装工具**

```bash
make install-proto-tools
```

生成的代码位置：

- 后端：`backend/gen/go/`
- 前端：`frontend/src/proto/`

---

## 开发工作流

### 添加新的 API 功能

1. **定义 Proto 文件** - `proto/my_service.proto`

```protobuf
service MyService {
  rpc CreateItem(CreateItemRequest) returns (CreateItemResponse) {
    option (google.api.http) = {
      post: "/api/v1/items"
      body: "*"
    };
  }
}
```

2. **生成代码**

```bash
make proto-backend
```

3. **实现 gRPC 服务** - `backend/internal/grpc/my_service.go`

```go
type MyServer struct {
    pb.UnimplementedMyServiceServer
    db *ent.Client
}

func (s *MyServer) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*pb.CreateItemResponse, error) {
    // 业务逻辑实现
}
```

4. **注册服务** - `backend/cmd/server/main.go`

```go
pb.RegisterMyServiceServer(grpcServer, grpcserver.NewMyServer(dbClient.Client))
```

5. **HTTP API 自动生成** ✅ 无需手动编写 REST handler

6. **前端调用** - `frontend/src/services/my_service.ts`

```typescript
import api from "@/services/api";

export const createItem = (data: CreateItemRequest) => {
  return api.post("/items", data);
};
```

**API 响应示例**（注意字段使用 snake_case）：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "conversation_id": "52df2b40-3ff3-4ce4-b869-432c33e7fc56",
    "agent_id": "abc123",
    "user_id": "user-001",
    "messages": [
      {
        "id": "b9ae5fc5-4581-4f7d-b967-72c75be859fd",
        "role": "user",
        "content": "hi",
        "timestamp": "2025-11-16T14:11:12.066929466Z"
      }
    ],
    "created_at": "2025-11-16T14:11:12Z",
    "updated_at": "2025-11-16T14:11:16Z"
  }
}
```

### 添加数据库模型

使用 ent ORM：

1. 定义 Schema - `backend/internal/model/schema/my_model.go`

```go
type MyModel struct {
    ent.Schema
}

func (MyModel) Fields() []ent.Field {
    return []ent.Field{
        field.String("name"),
        field.Time("created_at").Default(time.Now),
    }
}
```

2. 生成代码

```bash
cd backend
go generate ./...
```

3. 使用 Repository 模式访问数据

- 在 `backend/internal/repository/` 创建 repository
- 在 gRPC 服务中调用 repository

---

## 环境配置

### 后端环境变量

文件：`backend/.env`（参考 `.env.example`）

**核心配置**

```bash
# 服务端口
GRPC_PORT=9000              # gRPC 端口
HTTP_PORT=8000              # HTTP REST API 端口
SERVER_MODE=debug           # debug 或 release

# 数据库
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=123456
POSTGRES_DATABASE=agent_platform

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# AI 配置
OPENAI_API_KEY=sk-xxx                           # OpenAI API Key (可选)
SILICONFLOW_API_KEY=sk-xxx                      # SiliconFlow API Key (可选)
SILICONFLOW_API_BASE=https://api.siliconflow.cn/v1
SILICONFLOW_MODEL=deepseek-ai/DeepSeek-V3       # DeepSeek 模型

# Embedding 配置
EMBEDDING_MODEL=text-embedding-ada-002
EMBEDDING_DIMENSION=1536

# 安全
JWT_SECRET=your-secret-key
CORS_ORIGINS=http://localhost:5173
```

### 前端代理配置

`frontend/vite.config.ts` 已配置代理：

```typescript
proxy: {
  '/api': {
    target: 'http://localhost:8000',  # 后端 HTTP 端口
    changeOrigin: true,
  },
}
```

前端请求 `/api/v1/xxx` → 自动转发到 `http://localhost:8000/api/v1/xxx`

### AI 模型配置说明

系统支持多个 AI 服务提供商，可以同时配置或仅配置其中一个：

**1. OpenAI（默认）**
```bash
OPENAI_API_KEY=sk-xxx
OPENAI_API_BASE=https://api.openai.com/v1
```

**2. SiliconFlow DeepSeek（推荐 - 性价比高）**
```bash
SILICONFLOW_API_KEY=sk-xxx
SILICONFLOW_API_BASE=https://api.siliconflow.cn/v1
SILICONFLOW_MODEL=deepseek-ai/DeepSeek-V3
```

**模型路由规则**：
- 模型名包含 `deepseek` → 使用 SiliconFlow DeepSeek
- 其他模型（`gpt-4`、`gpt-3.5-turbo` 等）→ 使用 OpenAI
- 如果只配置了一个服务，自动使用该服务作为备选

**获取 SiliconFlow API Key**：
1. 访问 https://siliconflow.cn
2. 注册账号并登录
3. 前往控制台 → API Keys
4. 创建新的 API Key
5. 复制到 `.env` 文件的 `SILICONFLOW_API_KEY`

---

## 项目结构

```
agent-opus/
├── backend/
│   ├── cmd/server/
│   │   ├── main.go           # 主程序入口（启动 gRPC + HTTP Gateway）
│   │   ├── gateway.go        # HTTP Gateway 配置
│   │   └── marshaler.go      # 统一响应格式
│   ├── internal/
│   │   ├── grpc/             # gRPC 服务实现（核心业务逻辑）
│   │   ├── repository/       # 数据访问层
│   │   ├── model/schema/     # ent 数据模型定义
│   │   ├── config/           # 配置加载
│   │   ├── middleware/       # 中间件（认证、日志、CORS）
│   │   ├── ai/               # AI 模型调用
│   │   ├── knowledge/        # 知识库 + 向量检索
│   │   └── auth/             # JWT 认证
│   ├── gen/go/               # 生成的 proto 代码
│   └── .env                  # 环境变量配置
│
├── frontend/
│   ├── src/
│   │   ├── pages/            # 页面组件
│   │   ├── components/       # 公共组件
│   │   ├── services/
│   │   │   └── api.ts        # 统一的 API 客户端配置 ⭐
│   │   ├── store/            # Redux 状态管理
│   │   └── types/            # TypeScript 类型
│   └── vite.config.ts        # Vite 配置（含 proxy）
│
├── proto/                    # Proto 定义文件（服务契约）
│   ├── agent.proto
│   ├── conversation.proto
│   ├── tool.proto
│   ├── knowledge_base.proto
│   └── user.proto
│
├── Makefile                  # 常用命令快捷方式
├── docker-compose.yml        # Docker 编排
└── CLAUDE.md                 # 本文件
```

---

## 关键依赖

### 后端

- `google.golang.org/grpc` - gRPC 框架
- `grpc-ecosystem/grpc-gateway` - HTTP Gateway (REST API 自动生成)
- `entgo.io/ent` - ORM 框架
- `pgvector/pgvector-go` - PostgreSQL 向量扩展
- `sashabaranov/go-openai` - OpenAI/兼容 API 客户端 (支持 SiliconFlow)
- `go.uber.org/zap` - 结构化日志

### 前端

- `react` 18 + `react-router-dom` v6
- `antd` 5 - UI 组件库
- `@ant-design/pro-chat` - 专业级聊天 UI 组件
- `@reduxjs/toolkit` - 状态管理
- `axios` - HTTP 客户端
- `reactflow` - 工作流可视化

---

## 当前状态

**已完成 ✅**

- gRPC + HTTP Gateway 双协议架构
- 完整的 CRUD API（Agent、Conversation、Tool、KnowledgeBase）
- 前端 UI 界面（所有主要页面）
- JWT 认证
- PostgreSQL + pgvector 向量检索
- AI 对话集成（OpenAI + SiliconFlow DeepSeek）
- 多 AI 模型支持（智能路由）
- Docker 容器化部署

**待完善 🚧**

- 多 Agent 协作编排
- 工作流可视化编辑器
- 更多 AI 模型支持（Anthropic Claude、Google Gemini）
- WebSocket 实时通信
- 权限细粒度控制

---

## 开发规范

### 代码风格

- 后端：遵循 Go 官方风格，使用 `go fmt`
- 前端：遵循 ESLint 配置，使用 `npm run lint`

### API 设计原则

1. **单一数据源**：所有 API 定义在 proto 文件中
2. **统一响应格式**：使用 `{code, message, data, timestamp, request_id}`
3. **RESTful 规范**：使用 HTTP 方法语义（GET/POST/PUT/DELETE）
4. **版本化**：API 路径包含版本号 `/api/v1/`
5. **命名规范**：
   - Proto 字段使用 **snake_case**（`agent_id`, `created_at`）
   - HTTP API 响应保持 **snake_case**（配置：`UseProtoNames: true`）
   - 前端 TypeScript 接口使用 **snake_case** 与后端保持一致

### 前端开发规范

1. **使用统一 API 客户端**：导入 `@/services/api` 而非直接使用 axios
2. **类型安全**：所有 API 调用定义 TypeScript 接口
3. **状态管理**：复杂状态使用 Redux，局部状态用 `useState`
4. **组件拆分**：页面组件 < 200 行，提取可复用组件

---

## 故障排查

### 常见问题

**1. 后端启动失败 - 数据库连接错误**

- 检查 PostgreSQL 是否运行：`psql -U postgres -h localhost`
- 检查 `backend/.env` 中的数据库配置
- 确认数据库已创建：`CREATE DATABASE agent_platform;`

**2. Proto 代码生成失败**

```bash
# 重新安装工具
make install-proto-tools

# 验证 protoc 安装
protoc --version
```

**3. 前端 API 调用 404**

- 检查后端是否运行在 8000 端口
- 检查 `vite.config.ts` proxy 配置
- 确认 API 路径是否正确（`/api/v1/xxx`）

**4. CORS 错误**

- 检查 `backend/.env` 中的 `CORS_ORIGINS` 配置
- 确保包含前端地址（如 `http://localhost:5173`）

**5. Docker 端口冲突**

```bash
# 修改 docker-compose.yml 中的端口映射
ports:
  - "3307:3306"  # 如本地 3306 被占用
```

---

## 测试

### 后端测试

```bash
cd backend
go test ./...                    # 运行所有测试
go test ./internal/grpc/...      # 测试特定包
go test -v -run TestXxx          # 运行特定测试
```

### API 测试

```bash
# 使用脚本测试
./test-api.sh

# 手动测试
curl http://localhost:8000/api/v1/agents
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}'
```

### 前端测试

```bash
cd frontend
npm run test
```

---

## 资源链接

- [gRPC Gateway 文档](https://grpc-ecosystem.github.io/grpc-gateway/)
- [Ent ORM 文档](https://entgo.io/docs/getting-started)
- [pgvector 文档](https://github.com/pgvector/pgvector)
- [Ant Design 组件](https://ant.design/components/overview)
- [React Router v6](https://reactrouter.com/)

---

**最后更新**：2025-11-16
**维护者**：Agent Platform Team
