# 变更日志

## [未发布] - 2024-11-15

### 重大变更 🔥

#### 移除 MongoDB 依赖

**原因：**
- 简化数据库架构，减少技术栈复杂度
- MySQL 已能满足所有结构化数据存储需求
- 降低部署和运维成本

**影响的文件：**
- `internal/config/config.go` - 移除 `MongoDBConfig` 结构体和相关配置加载代码
- `.env` - 移除 `MONGODB_URI` 和 `MONGODB_DATABASE` 配置项
- `.env.example` - 移除 MongoDB 配置示例
- `REFACTORING.md` - 更新重构说明，添加数据库架构说明
- `DATABASE.md` - 新增数据库配置和使用文档

**迁移指南：**

如果你之前使用了 MongoDB，需要：
1. 将 MongoDB 中的数据迁移到 MySQL
2. 更新 `.env` 文件，移除 MongoDB 相关配置
3. 重新编译项目：`make build`

详细迁移步骤请参考 [DATABASE.md](./DATABASE.md)

**当前数据库架构：**
- MySQL - 主数据库（结构化数据）
- Redis - 缓存数据库
- Milvus - 向量数据库

### 新增 ✨

#### Proto 定义和代码生成

**新增文件：**
- `api/proto/common.proto` - 通用响应结构和分页定义
- `api/proto/agent.proto` - Agent 相关的 proto 定义
- `api/proto/conversation.proto` - 对话相关的 proto 定义
- `api/proto/tool.proto` - 工具相关的 proto 定义
- `api/proto/knowledge_base.proto` - 知识库相关的 proto 定义
- `Makefile` - 构建脚本，支持 proto 编译

**使用方式：**
```bash
# 生成 proto 代码
make proto

# 清理生成的代码
make proto-clean
```

#### 统一响应格式

**新增文件：**
- `pkg/response/response.go` - 统一的 HTTP 响应处理工具包

**响应格式：**
```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

**可用函数：**
- `Success(c, data)` - 返回成功响应
- `Created(c, data)` - 返回创建成功响应（201）
- `BadRequest(c, message)` - 返回 400 错误
- `Unauthorized(c, message)` - 返回 401 错误
- `NotFound(c, message)` - 返回 404 错误
- `InternalServerError(c, message)` - 返回 500 错误

#### 文档

**新增文件：**
- `REFACTORING.md` - 重构说明文档
- `DATABASE.md` - 数据库配置和使用指南
- `CHANGELOG.md` - 变更日志（本文件）

### 改进 🚀

#### Handler 重构

所有 Handler 已重构为使用 proto 定义的消息类型：

**修改的文件：**
- `internal/handler/agent_handler.go` - 使用 `pb.Agent` 和相关 proto 消息
- `internal/handler/conversation_handler.go` - 使用 `pb.Conversation` 和 `pb.Message`
- `internal/handler/tool_handler.go` - 使用 `pb.Tool` 相关消息
- `internal/handler/knowledge_base_handler.go` - 使用 `pb.KnowledgeBase` 和 `pb.Document`

**改进点：**
- 类型安全：使用强类型的 proto 消息，减少运行时错误
- 统一格式：所有接口使用统一的响应格式
- 代码简洁：使用 response 包的工具函数，减少重复代码
- 易于维护：proto 定义作为接口契约，便于前后端协作

**示例对比：**

Before:
```go
func (h *AgentHandler) CreateAgent(c *gin.Context) {
    var req struct {
        Name string `json:"name" binding:"required"`
        // ...
    }
    // ...
    c.JSON(http.StatusCreated, gin.H{
        "code": 0,
        "message": "success",
        "data": agent,
    })
}
```

After:
```go
func (h *AgentHandler) CreateAgent(c *gin.Context) {
    var req pb.CreateAgentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, err.Error())
        return
    }
    // ...
    response.Created(c, agentMap)
}
```

### 优化 ⚡

#### 构建流程

**新增 Makefile 命令：**
```bash
make proto          # 生成 proto 代码
make proto-clean    # 清理生成的 proto 代码
make build          # 构建项目
make run            # 运行项目
make dev            # proto + build + run
make install-tools  # 安装 protoc 插件
make clean          # 清理所有构建产物
make help           # 显示帮助信息
```

#### 代码组织

- 统一导入路径为 `agent-platform/*`
- 清晰的目录结构：`api/proto/`, `pkg/response/`, `internal/handler/`
- 分离关注点：proto 定义、响应处理、业务逻辑

### 修复 🐛

- 修复了 proto 生成代码的导入路径问题
- 修复了 response 包中未使用的导入
- 清理了 go.mod 中的未使用依赖

### 技术债务 📝

**已解决：**
- ✅ 统一响应格式
- ✅ 使用 proto 定义接口
- ✅ 简化数据库架构

**待解决：**
- ⏳ 实现 Service 层，将业务逻辑从 Handler 分离
- ⏳ 实现 Repository 层，完成真实的数据库访问
- ⏳ 添加单元测试和集成测试
- ⏳ 实现完整的错误码体系
- ⏳ 添加 API 文档生成（基于 proto）
- ⏳ 性能优化：减少 proto 到 map 的转换

## 升级指南

### 从旧版本升级

1. **备份数据**
   ```bash
   # 如果使用了 MongoDB，先备份数据
   mongodump --db agent_platform --out backup/
   ```

2. **更新代码**
   ```bash
   git pull origin main
   ```

3. **更新配置文件**
   ```bash
   # 复制新的配置示例
   cp .env.example .env

   # 编辑 .env，移除 MongoDB 配置
   # 确保 MySQL、Redis、Milvus 配置正确
   ```

4. **安装 protoc 工具**
   ```bash
   make install-tools
   ```

5. **生成 proto 代码**
   ```bash
   make proto
   ```

6. **更新依赖**
   ```bash
   go mod tidy
   ```

7. **编译项目**
   ```bash
   make build
   ```

8. **数据迁移**（如果之前使用了 MongoDB）
   ```bash
   # 参考 DATABASE.md 中的迁移指南
   ```

9. **测试运行**
   ```bash
   make run
   ```

## 兼容性说明

### 破坏性变更

- ❌ 移除了 MongoDB 支持
- ❌ 所有 API 响应格式变更为 `{code, message, data}`
- ❌ Handler 的导入路径和包结构发生变化

### 向后兼容

- ✅ MySQL、Redis、Milvus 配置保持兼容
- ✅ 环境变量格式不变（除 MongoDB 相关）
- ✅ 核心业务逻辑不受影响

## 贡献者

感谢所有参与此次重构的贡献者！

## 相关链接

- [重构说明文档](./REFACTORING.md)
- [数据库配置指南](./DATABASE.md)
- [Protocol Buffers 文档](https://protobuf.dev/)
- [Gin Web Framework](https://gin-gonic.com/)

---

**注意：** 本次更新包含破坏性变更，升级前请仔细阅读升级指南和相关文档。
