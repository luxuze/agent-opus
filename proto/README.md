# Proto 定义文件（前后端共享）

本目录包含所有的 Protocol Buffers 定义文件，用于前后端接口的统一定义。

## ✨ 新特性：同时支持 gRPC 和 HTTP REST API

所有服务现在都同时支持：

- **gRPC** - 高性能二进制协议（端口 9000）
- **HTTP REST API** - 标准 RESTful JSON API（端口 8000）

两种方式使用相同的 proto 定义，确保接口一致性。

## 文件说明

- **`common.proto`** - 通用数据结构和响应格式
- **`agent.proto`** - Agent 服务定义和消息（含 HTTP 路由）
- **`conversation.proto`** - Conversation 服务定义和消息（含 HTTP 路由）
- **`tool.proto`** - Tool 服务定义和消息（含 HTTP 路由）
- **`knowledge_base.proto`** - KnowledgeBase 服务定义和消息（含 HTTP 路由）
- **`google/api/`** - Google API 注解（用于定义 HTTP 路由）
- **`HTTP_API.md`** - HTTP REST API 详细文档

## 生成代码

### 方式一：从项目根目录（推荐）

```bash
# 生成所有前后端代码
make proto

# 仅生成后端 Go 代码
make proto-backend

# 仅生成前端 TypeScript/JavaScript 代码
make proto-frontend
```

### 方式二：从后端目录

```bash
cd backend
make proto
```

### 方式三：从前端目录

```bash
cd frontend
npm run proto
```

## 生成的代码位置

- **后端 Go 代码**: `backend/gen/go/`

  - `*.pb.go` - Protobuf 消息定义
  - `*_grpc.pb.go` - gRPC 服务端代码

- **前端 TypeScript/JavaScript 代码**: `frontend/src/proto/`
  - `*_pb.js` - Protobuf 消息定义
  - `*_grpc_web_pb.js` - gRPC-Web 客户端代码

## 修改 Proto 文件

1. 编辑本目录下的 `.proto` 文件
2. 运行 `make proto` 重新生成前后端代码
3. 在后端实现新的 gRPC 方法（`backend/internal/grpc/`）
4. 在前端使用生成的客户端代码

## 注意事项

- 修改 proto 文件后必须重新生成代码
- 前后端使用同一份 proto 定义，确保接口一致性
- 生成的代码文件不应手动编辑
- 添加到 `.gitignore` 中的生成目录：
  - `backend/gen/go/`
  - `frontend/src/proto/`

## 工具安装

### 一键安装（推荐）

从项目根目录运行：

```bash
make install-proto-tools
```

这会自动安装所有必需的工具。

### 手动安装

#### 后端工具（必需）

```bash
# 安装 Go protoc 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

#### 前端工具（可选）

**macOS (使用 Homebrew):**

```bash
brew install protoc-gen-grpc-web protoc-gen-js
```

**其他系统:**

1. 下载 `protoc-gen-grpc-web`:

   - 访问 https://github.com/grpc/grpc-web/releases
   - 下载对应平台的二进制文件
   - 将其移动到 PATH 目录中（如 `/usr/local/bin`）
   - 添加执行权限：`chmod +x protoc-gen-grpc-web`

2. 下载 `protoc-gen-js`:
   - 访问 https://github.com/protocolbuffers/protobuf-javascript/releases
   - 下载对应平台的二进制文件
   - 将其移动到 PATH 目录中
   - 添加执行权限：`chmod +x protoc-gen-js`

**验证安装:**

```bash
which protoc-gen-go
which protoc-gen-go-grpc
which protoc-gen-grpc-web
which protoc-gen-js
```

## 更多信息

- [Protocol Buffers 文档](https://protobuf.dev/)
- [gRPC 官方文档](https://grpc.io/)
- [gRPC-Web 文档](https://github.com/grpc/grpc-web)
