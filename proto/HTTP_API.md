# HTTP REST API 文档

本项目同时支持 **gRPC** 和 **HTTP REST API**，通过 gRPC-Gateway 自动生成 RESTful 接口。

## 架构说明

```
┌─────────────┐
│   Frontend  │
│   (Client)  │
└──────┬──────┘
       │
       ├──────────────────┐
       │                  │
       ▼                  ▼
  HTTP REST           gRPC-Web
 (JSON格式)         (Protobuf)
       │                  │
       ▼                  ▼
┌─────────────────────────────┐
│    gRPC-Gateway/Envoy       │
│  (HTTP -> gRPC 转换)        │
└──────────┬──────────────────┘
           │
           ▼
      ┌─────────┐
      │ Backend │
      │  (gRPC) │
      └─────────┘
```

## API 路径规范

所有 HTTP REST API 都遵循统一的路径规范：`/api/v1/{resource}`

### Agent Service

| 方法   | 路径                | 描述       | gRPC 方法   |
| ------ | ------------------- | ---------- | ----------- |
| POST   | /api/v1/agents      | 创建 Agent | CreateAgent |
| GET    | /api/v1/agents      | 获取列表   | ListAgents  |
| GET    | /api/v1/agents/{id} | 获取详情   | GetAgent    |
| PUT    | /api/v1/agents/{id} | 更新 Agent | UpdateAgent |
| DELETE | /api/v1/agents/{id} | 删除 Agent | DeleteAgent |

### Conversation Service

| 方法 | 路径                                             | 描述     | gRPC 方法          |
| ---- | ------------------------------------------------ | -------- | ------------------ |
| POST | /api/v1/conversations                            | 创建对话 | CreateConversation |
| GET  | /api/v1/conversations                            | 获取列表 | ListConversations  |
| GET  | /api/v1/conversations/{id}                       | 获取详情 | GetConversation    |
| POST | /api/v1/conversations/{conversation_id}/messages | 发送消息 | SendMessage        |

### Tool Service

| 方法   | 路径               | 描述     | gRPC 方法  |
| ------ | ------------------ | -------- | ---------- |
| POST   | /api/v1/tools      | 创建工具 | CreateTool |
| GET    | /api/v1/tools      | 获取列表 | ListTools  |
| GET    | /api/v1/tools/{id} | 获取详情 | GetTool    |
| DELETE | /api/v1/tools/{id} | 删除工具 | DeleteTool |

### KnowledgeBase Service

| 方法   | 路径                                                  | 描述       | gRPC 方法           |
| ------ | ----------------------------------------------------- | ---------- | ------------------- |
| POST   | /api/v1/knowledge-bases                               | 创建知识库 | CreateKnowledgeBase |
| GET    | /api/v1/knowledge-bases                               | 获取列表   | ListKnowledgeBases  |
| GET    | /api/v1/knowledge-bases/{id}                          | 获取详情   | GetKnowledgeBase    |
| POST   | /api/v1/knowledge-bases/{knowledge_base_id}/documents | 上传文档   | UploadDocument      |
| DELETE | /api/v1/knowledge-bases/{id}                          | 删除知识库 | DeleteKnowledgeBase |

## 使用示例

### 创建 Agent

**请求：**

```bash
curl -X POST http://localhost:8000/api/v1/agents \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Customer Service Agent",
    "description": "AI agent for customer service",
    "type": "single",
    "model_config": {
      "model": "gpt-4",
      "temperature": 0.7
    }
  }'
```

**响应：**

```json
{
  "id": "agent-123",
  "name": "Customer Service Agent",
  "description": "AI agent for customer service",
  "type": "single",
  "model_config": {
    "model": "gpt-4",
    "temperature": 0.7
  },
  "status": "draft",
  "version": "1.0.0",
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

### 获取 Agent 列表

**请求：**

```bash
curl -X GET "http://localhost:8000/api/v1/agents?page=1&page_size=10"
```

**响应：**

```json
{
  "items": [
    {
      "id": "agent-123",
      "name": "Customer Service Agent",
      "description": "AI agent for customer service",
      "type": "single",
      "status": "published"
    }
  ],
  "page": 1,
  "page_size": 10,
  "total": 1
}
```

### 获取 Agent 详情

**请求：**

```bash
curl -X GET http://localhost:8000/api/v1/agents/agent-123
```

### 更新 Agent

**请求：**

```bash
curl -X PUT http://localhost:8000/api/v1/agents/agent-123 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Agent Name",
    "description": "Updated description",
    "status": "published"
  }'
```

### 删除 Agent

**请求：**

```bash
curl -X DELETE http://localhost:8000/api/v1/agents/agent-123
```

### 创建对话

**请求：**

```bash
curl -X POST http://localhost:8000/api/v1/conversations \
  -H "Content-Type: application/json" \
  -d '{
    "agent_id": "agent-123",
    "title": "Customer Support Chat"
  }'
```

### 发送消息

**请求：**

```bash
curl -X POST http://localhost:8000/api/v1/conversations/conv-456/messages \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Hello, I need help",
    "role": "user"
  }'
```

**响应：**

```json
{
  "conversation_id": "conv-456",
  "messages": [
    {
      "id": "msg-1",
      "role": "user",
      "content": "Hello, I need help",
      "timestamp": "2024-01-15T10:05:00Z"
    },
    {
      "id": "msg-2",
      "role": "assistant",
      "content": "Hello! How can I help you today?",
      "timestamp": "2024-01-15T10:05:01Z"
    }
  ]
}
```

## 错误处理

HTTP REST API 使用标准的 HTTP 状态码：

- **200 OK** - 请求成功
- **201 Created** - 资源创建成功
- **400 Bad Request** - 请求参数错误
- **404 Not Found** - 资源不存在
- **500 Internal Server Error** - 服务器内部错误

**错误响应示例：**

```json
{
  "code": 3,
  "message": "invalid argument: name is required",
  "details": []
}
```

## 前端集成

### 使用 Fetch API

```typescript
// 创建 Agent
async function createAgent(data) {
  const response = await fetch("http://localhost:8000/api/v1/agents", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  return await response.json();
}

// 获取列表
async function listAgents(page = 1, pageSize = 10) {
  const response = await fetch(
    `http://localhost:8000/api/v1/agents?page=${page}&page_size=${pageSize}`
  );
  return await response.json();
}
```

### 使用 Axios

```typescript
import axios from "axios";

const apiClient = axios.create({
  baseURL: "http://localhost:8000/api/v1",
  headers: {
    "Content-Type": "application/json",
  },
});

// 创建 Agent
export const createAgent = (data) => apiClient.post("/agents", data);

// 获取列表
export const listAgents = (page, pageSize) =>
  apiClient.get("/agents", { params: { page, page_size: pageSize } });

// 获取详情
export const getAgent = (id) => apiClient.get(`/agents/${id}`);

// 更新
export const updateAgent = (id, data) => apiClient.put(`/agents/${id}`, data);

// 删除
export const deleteAgent = (id) => apiClient.delete(`/agents/${id}`);
```

## 认证和授权

在生产环境中，建议添加认证：

```typescript
const apiClient = axios.create({
  baseURL: "http://localhost:8000/api/v1",
  headers: {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  },
});
```

## 开发调试

### 使用 Postman

1. 创建新的 Collection
2. 添加环境变量：
   - `base_url`: `http://localhost:8000/api/v1`
3. 创建请求，使用 `{{base_url}}/agents`

### 使用 HTTPie

```bash
# 创建 Agent
http POST localhost:8000/api/v1/agents \
  name="Test Agent" \
  description="A test agent" \
  type="single"

# 获取列表
http GET localhost:8000/api/v1/agents page==1 page_size==10
```

## 最佳实践

1. **使用环境变量**：不要硬编码 API 地址

   ```typescript
   const API_BASE_URL =
     process.env.REACT_APP_API_URL || "http://localhost:8000/api/v1";
   ```

2. **错误处理**：统一处理 API 错误

   ```typescript
   apiClient.interceptors.response.use(
     (response) => response,
     (error) => {
       console.error("API Error:", error.response?.data);
       return Promise.reject(error);
     }
   );
   ```

3. **请求超时**：设置合理的超时时间

   ```typescript
   const apiClient = axios.create({
     baseURL: API_BASE_URL,
     timeout: 10000, // 10 seconds
   });
   ```

4. **数据验证**：在发送前验证数据
   ```typescript
   function validateAgent(data) {
     if (!data.name) throw new Error("Name is required");
     if (!data.type) throw new Error("Type is required");
     return true;
   }
   ```

## 对比：gRPC vs REST

| 特性       | gRPC              | REST HTTP          |
| ---------- | ----------------- | ------------------ |
| 协议       | HTTP/2            | HTTP/1.1           |
| 数据格式   | Protobuf (二进制) | JSON (文本)        |
| 性能       | 更快              | 较慢               |
| 浏览器支持 | 需要 gRPC-Web     | 原生支持           |
| 开发工具   | 需要特殊工具      | 通用工具 (curl 等) |
| 使用场景   | 内部服务、高性能  | 公共 API、简单集成 |

## 相关文档

- [gRPC 服务架构](../backend/GRPC.md)
- [Proto 文件说明](./README.md)
- [gRPC-Gateway 官方文档](https://grpc-ecosystem.github.io/grpc-gateway/)
