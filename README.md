# Agent 平台

一个功能完整的智能代理管理和编排系统，支持 Agent 创建、管理、对话、工具集成和知识库管理。

## 项目概述

Agent 平台是一个面向企业和开发者的智能代理管理系统，提供：

- 🤖 **Agent 管理** - 创建、配置和管理各种 AI Agent
- 💬 **对话管理** - 完整的对话会话管理和历史记录
- 🛠️ **工具系统** - 内置和自定义工具库
- 📚 **知识库** - 文档向量化和检索增强生成
- 🔄 **工作流编排** - Multi-Agent 协作和任务编排
- 📊 **监控分析** - 完善的监控和性能分析

## 技术栈

### 后端

- **语言**: Go 1.21+
- **框架**: Gin
- **ORM**: ent
- **数据库**: MySQL 8.0 + MongoDB 6 + Redis 7

### 前端

- **框架**: React 18
- **构建工具**: Vite
- **UI 库**: Ant Design 5
- **状态管理**: Redux Toolkit
- **路由**: React Router v6

### 基础设施

- **容器化**: Docker + Docker Compose
- **反向代理**: Nginx

## 项目结构

```
agent-opus/
├── backend/                    # 后端服务
│   ├── cmd/
│   │   └── server/            # 主程序入口
│   ├── internal/
│   │   ├── config/            # 配置管理
│   │   ├── handler/           # HTTP 处理器
│   │   ├── middleware/        # 中间件
│   │   ├── model/             # 数据模型
│   │   ├── repository/        # 数据访问层
│   │   └── service/           # 业务逻辑层
│   ├── pkg/                   # 公共包
│   ├── go.mod
│   └── Dockerfile
├── frontend/                   # 前端应用
│   ├── src/
│   │   ├── components/        # React 组件
│   │   ├── pages/             # 页面组件
│   │   ├── services/          # API 服务
│   │   ├── store/             # Redux store
│   │   └── types/             # TypeScript 类型
│   ├── package.json
│   └── Dockerfile
├── docs/                       # 文档
├── scripts/                    # 脚本
├── docker-compose.yml          # Docker Compose 配置
├── agent-platform-requirements.md  # 需求文档
└── README.md
```

## 快速开始

### 前置要求

- Docker 20.10+
- Docker Compose 2.0+
- (可选) Node.js 20+ (本地开发)
- (可选) Go 1.21+ (本地开发)

### 使用 Docker Compose 启动

1. **克隆项目**

```bash
cd /Users/xuzelu/workspace/agent-opus
```

2. **配置环境变量**

```bash
cp backend/.env.example backend/.env
# 编辑 backend/.env 配置必要的环境变量
```

3. **启动所有服务**

```bash
docker-compose up -d
```

4. **查看服务状态**

```bash
docker-compose ps
```

5. **访问应用**

- 前端: http://localhost:3000
- 后端 API: http://localhost:8000
- API 文档: http://localhost:8000/api/v1/ping

### 本地开发

#### 后端开发

```bash
cd backend

# 安装依赖
go mod download

# 复制配置文件
cp .env.example .env

# 启动数据库 (使用 Docker)
docker-compose up -d mysql mongodb redis

# 运行服务
go run cmd/server/main.go
```

#### 前端开发

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

访问 http://localhost:5173

## API 文档

### 核心 API 端点

#### Agent 管理

```bash
# 创建 Agent
POST /api/v1/agents

# 获取 Agent 列表
GET /api/v1/agents

# 获取 Agent 详情
GET /api/v1/agents/{id}

# 更新 Agent
PUT /api/v1/agents/{id}

# 删除 Agent
DELETE /api/v1/agents/{id}
```

#### 对话管理

```bash
# 创建对话
POST /api/v1/conversations

# 发送消息
POST /api/v1/conversations/{id}/messages

# 获取对话详情
GET /api/v1/conversations/{id}
```

#### 工具管理

```bash
# 获取工具列表
GET /api/v1/tools

# 创建工具
POST /api/v1/tools
```

#### 知识库管理

```bash
# 创建知识库
POST /api/v1/knowledge-bases

# 上传文档
POST /api/v1/knowledge-bases/{id}/documents
```

### API 响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "timestamp": 1634567890,
  "request_id": "uuid"
}
```

## 配置说明

### 后端环境变量

| 变量名         | 说明                     | 默认值                    |
| -------------- | ------------------------ | ------------------------- |
| SERVER_PORT    | 服务端口                 | 8000                      |
| SERVER_MODE    | 运行模式 (debug/release) | debug                     |
| MYSQL_HOST     | MySQL 主机               | localhost                 |
| MYSQL_PORT     | MySQL 端口               | 3306                      |
| MYSQL_DATABASE | 数据库名                 | agent_platform            |
| MONGODB_URI    | MongoDB 连接字符串       | mongodb://localhost:27017 |
| REDIS_HOST     | Redis 主机               | localhost                 |
| JWT_SECRET     | JWT 密钥                 | your-secret-key           |
| OPENAI_API_KEY | OpenAI API Key           | -                         |

完整配置见 `backend/.env.example`

## 部署

### Docker 部署

使用提供的 `docker-compose.yml` 文件可以一键部署所有服务：

```bash
# 生产环境部署
docker-compose -f docker-compose.yml up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### Kubernetes 部署

参考 `docs/kubernetes-deployment.md` (待补充)

## 开发指南

### 后端开发

1. **添加新的 API 端点**

   - 在 `internal/handler` 添加处理器
   - 在 `cmd/server/router.go` 注册路由
   - 在 `internal/service` 添加业务逻辑

2. **添加新的数据模型**
   - 在 `internal/model/schema` 定义 ent schema
   - 运行 `go generate ./...` 生成代码

### 前端开发

1. **添加新页面**

   - 在 `src/pages` 创建页面组件
   - 在 `App.tsx` 添加路由配置

2. **添加 API 服务**
   - 在 `src/services` 添加服务文件
   - 定义 API 调用方法

## 测试

### 后端测试

```bash
cd backend
go test ./...
```

### 前端测试

```bash
cd frontend
npm run test
```

## 常见问题

### Q: 如何重置数据库？

```bash
docker-compose down -v
docker-compose up -d
```

### Q: 如何查看日志？

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f backend
docker-compose logs -f frontend
```

### Q: 前端代理配置

前端开发时，API 请求会通过 Vite 代理转发到后端。配置在 `frontend/vite.config.ts`:

```typescript
proxy: {
  '/api': {
    target: 'http://localhost:8000',
    changeOrigin: true,
  },
}
```

## 性能优化

- 使用 Redis 缓存热点数据
- 数据库查询优化和索引
- 前端代码分割和懒加载
- CDN 加速静态资源

## 安全性

- JWT Token 认证
- CORS 配置
- SQL 注入防护
- XSS 攻击防护
- 敏感数据加密

## 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 许可证

MIT License

## 联系方式

- 项目地址: https://github.com/yourname/agent-platform
- 问题反馈: https://github.com/yourname/agent-platform/issues

## 致谢

感谢以下开源项目：

- [Gin](https://github.com/gin-gonic/gin)
- [ent](https://entgo.io/)
- [React](https://react.dev/)
- [Ant Design](https://ant.design/)
- [Vite](https://vitejs.dev/)
