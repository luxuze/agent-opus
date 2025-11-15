# 快速启动指南

由于 Docker 构建可能需要较长时间，以下是本地快速启动的方式：

## 方式一：本地开发（推荐用于测试）

### 1. 启动后端

```bash
cd backend

# 确保已安装 Go 1.21+
go version

# 安装依赖（已完成）
go mod tidy

# 直接运行（不需要数据库也可以测试 API）
go run cmd/server/main.go
```

后端会在 http://localhost:8080 启动

测试 API：
```bash
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/ping
```

### 2. 启动前端

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端会在 http://localhost:5173 启动

## 方式二：使用 Docker Compose

### 解决端口冲突

如果您本地已有 MySQL (3306) 或 Redis (6379) 运行，我们已将端口修改为：
- MySQL: 3307 -> 3306 (容器内)
- MongoDB: 27017 (保持不变)
- Redis: 需要修改或停止本地 Redis

### 修改 docker-compose.yml

将 Redis 和 MongoDB 端口也改为不冲突的端口：

```yaml
  redis:
    ports:
      - "6380:6379"  # 改用 6380

  mongodb:
    ports:
      - "27018:27017"  # 改用 27018
```

然后更新 backend/.env 中的配置：

```env
MYSQL_PORT=3307
REDIS_PORT=6380
MONGODB_URI=mongodb://root:rootpassword@localhost:27018
```

### 启动服务

```bash
# 清理旧容器
docker-compose down -v

# 启动所有服务
docker-compose up -d --build
```

## 方式三：只用 Docker 数据库 + 本地运行代码

这是最快的测试方式：

### 1. 修改端口避免冲突

编辑 `docker-compose.yml`，修改所有端口映射：

```yaml
services:
  mysql:
    ports:
      - "3307:3306"
  mongodb:
    ports:
      - "27018:27017"
  redis:
    ports:
      - "6380:6379"
```

### 2. 只启动数据库

```bash
docker-compose up -d mysql mongodb redis
```

### 3. 更新后端配置

编辑 `backend/.env`:

```env
MYSQL_HOST=localhost
MYSQL_PORT=3307
MONGODB_URI=mongodb://root:rootpassword@localhost:27018
REDIS_HOST=localhost
REDIS_PORT=6380
```

### 4. 运行后端

```bash
cd backend
go run cmd/server/main.go
```

### 5. 运行前端

```bash
cd frontend
npm install
npm run dev
```

## 快速测试（无需数据库）

后端 API 目前使用 Mock 数据，可以直接运行测试：

```bash
# 启动后端
cd backend
go run cmd/server/main.go

# 在另一个终端测试 API
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/agents
curl -X POST http://localhost:8080/api/v1/agents \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Agent","description":"A test agent","type":"single"}'
```

## 常见问题

### Q: 端口被占用怎么办？

查看占用端口的进程：
```bash
# MacOS/Linux
lsof -i :3306
lsof -i :6379
lsof -i :27017

# 或停止本地服务
brew services stop mysql
brew services stop redis
brew services stop mongodb-community
```

### Q: Docker 构建很慢怎么办？

使用方式一（本地开发）或方式三（Docker 数据库 + 本地代码）

### Q: 前端无法连接后端？

确保：
1. 后端在 8080 端口运行
2. 前端 vite.config.ts 中的代理配置正确
3. CORS 配置包含了前端地址

## 项目结构

```
├── backend/         # Go 后端
│   └── cmd/server/main.go  # 主入口
├── frontend/        # React 前端
│   └── src/main.tsx # 前端入口
└── docker-compose.yml
```

## 下一步

1. 查看 README.md 了解完整功能
2. 查看 agent-platform-requirements.md 了解需求
3. 查看 docs/deployment.md 了解生产部署
