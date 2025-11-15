# 数据库配置说明

## 数据库架构

项目采用多数据库架构，每个数据库服务于不同的用途：

### 1. MySQL - 主数据库

**用途：** 存储结构化业务数据

**配置项（.env）：**
```bash
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=your_password
MYSQL_DATABASE=agent_platform
```

**存储内容：**
- Agent 配置和元数据
- 用户信息和权限
- 对话历史记录
- 工具定义
- 知识库元数据

**优势：**
- 支持复杂查询和事务
- 成熟稳定，生态丰富
- 适合结构化数据

### 2. Redis - 缓存数据库

**用途：** 缓存和会话存储

**配置项（.env）：**
```bash
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

**使用场景：**
- 用户会话管理
- JWT Token 黑名单
- 热点数据缓存
- 分布式锁
- 消息队列

**优势：**
- 高性能读写
- 丰富的数据结构
- 支持过期时间

### 3. Milvus - 向量数据库

**用途：** 存储和检索向量数据

**配置项（.env）：**
```bash
MILVUS_HOST=localhost
MILVUS_PORT=19530
```

**使用场景：**
- 文档向量存储
- 语义搜索
- 知识库检索
- 相似度匹配

**优势：**
- 专为向量检索优化
- 支持大规模向量数据
- 高效的相似度搜索

## 为什么移除 MongoDB？

### 架构简化

1. **减少技术栈复杂度**
   - 减少一个数据库服务的部署和维护
   - 降低学习和运维成本
   - 简化数据一致性管理

2. **功能重叠**
   - MySQL 已能满足结构化数据存储需求
   - Redis 提供了高性能缓存能力
   - Milvus 专注于向量检索

3. **性能和成本**
   - 减少跨数据库查询的开销
   - 降低服务器资源消耗
   - 简化数据备份和恢复流程

### 数据迁移建议

如果之前使用了 MongoDB 存储数据，建议按以下步骤迁移：

1. **识别 MongoDB 中的数据类型**
   - 结构化数据 → 迁移到 MySQL
   - 缓存数据 → 迁移到 Redis
   - 向量数据 → 迁移到 Milvus

2. **创建 MySQL 表结构**
   ```sql
   -- 使用 Ent 自动生成表结构
   -- 或手动创建对应的表
   ```

3. **编写数据迁移脚本**
   - 从 MongoDB 导出数据
   - 转换数据格式
   - 导入到对应的数据库

4. **验证数据完整性**
   - 对比数据条数
   - 验证关键字段
   - 测试查询功能

## 数据库连接示例

### MySQL 连接

```go
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

// 使用配置文件中的 DSN
db, err := sql.Open("mysql", cfg.MySQL.DSN())
if err != nil {
    log.Fatal(err)
}
```

### Redis 连接

```go
import (
    "github.com/go-redis/redis/v8"
)

client := redis.NewClient(&redis.Options{
    Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
    Password: cfg.Redis.Password,
    DB:       cfg.Redis.DB,
})
```

### Milvus 连接

```go
import (
    "github.com/milvus-io/milvus-sdk-go/v2/client"
)

milvusClient, err := client.NewGrpcClient(
    context.Background(),
    cfg.Milvus.Host + ":" + cfg.Milvus.Port,
)
```

## 环境配置

### 开发环境

使用 Docker Compose 快速启动所有数据库服务：

```yaml
version: '3.8'
services:
  mysql:
    image: mysql:8.0
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: root_password
      MYSQL_DATABASE: agent_platform
      MYSQL_USER: agent_user
      MYSQL_PASSWORD: agent_password
    volumes:
      - mysql_data:/var/lib/mysql

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  milvus:
    image: milvusdb/milvus:latest
    ports:
      - "19530:19530"
    volumes:
      - milvus_data:/var/lib/milvus

volumes:
  mysql_data:
  redis_data:
  milvus_data:
```

### 生产环境

建议使用云服务商提供的托管数据库服务：

- **MySQL**: AWS RDS, Google Cloud SQL, Azure Database for MySQL
- **Redis**: AWS ElastiCache, Google Cloud Memorystore, Azure Cache for Redis
- **Milvus**: Zilliz Cloud (托管 Milvus 服务)

## 性能优化建议

### MySQL

1. **索引优化**
   - 为常用查询字段创建索引
   - 避免在大字段上创建索引
   - 定期分析和优化索引

2. **查询优化**
   - 使用 EXPLAIN 分析查询计划
   - 避免 SELECT *
   - 使用分页查询大数据集

3. **连接池配置**
   ```go
   db.SetMaxOpenConns(100)
   db.SetMaxIdleConns(10)
   db.SetConnMaxLifetime(time.Hour)
   ```

### Redis

1. **合理设置过期时间**
   - 避免缓存雪崩
   - 使用随机过期时间

2. **选择合适的数据结构**
   - String: 简单键值
   - Hash: 对象存储
   - Set: 去重集合
   - ZSet: 排行榜

3. **监控内存使用**
   - 设置 maxmemory
   - 配置淘汰策略

### Milvus

1. **合理设置向量维度**
   - 根据模型选择维度（如 text-embedding-ada-002 使用 1536）

2. **选择合适的索引类型**
   - IVF_FLAT: 平衡性能和精度
   - HNSW: 高性能，高内存消耗
   - IVF_PQ: 压缩存储

3. **批量操作**
   - 批量插入向量
   - 批量搜索请求

## 备份策略

### MySQL 备份

```bash
# 定期全量备份
mysqldump -u root -p agent_platform > backup_$(date +%Y%m%d).sql

# 增量备份（启用 binlog）
mysqlbinlog --start-datetime="2024-01-01 00:00:00" mysql-bin.000001 > incremental.sql
```

### Redis 备份

```bash
# RDB 快照备份（自动）
# 配置 redis.conf
save 900 1
save 300 10
save 60 10000

# AOF 备份（手动）
redis-cli BGSAVE
```

### Milvus 备份

```bash
# 使用 Milvus 备份工具
# 参考: https://milvus.io/docs/backup_and_restore.md
```

## 监控和告警

### 监控指标

1. **MySQL**
   - 连接数
   - 查询延迟
   - 慢查询数
   - 表锁等待

2. **Redis**
   - 内存使用率
   - 命中率
   - 键过期数
   - 连接数

3. **Milvus**
   - 向量检索延迟
   - 索引构建状态
   - 内存使用

### 推荐工具

- **Prometheus + Grafana**: 指标采集和可视化
- **mysql_exporter**: MySQL 指标导出
- **redis_exporter**: Redis 指标导出
- **Milvus 内置监控**: 通过 Prometheus 接口

## 故障排查

### 常见问题

1. **MySQL 连接失败**
   ```
   Error: can't connect to MySQL server
   解决：检查 MySQL 服务状态、端口、防火墙
   ```

2. **Redis 连接超时**
   ```
   Error: i/o timeout
   解决：检查 Redis 服务、网络连接、超时配置
   ```

3. **Milvus 检索慢**
   ```
   解决：检查索引类型、向量维度、批量大小
   ```

## 总结

当前的数据库架构专注于三个核心数据库：

- **MySQL** - 可靠的结构化数据存储
- **Redis** - 高性能缓存层
- **Milvus** - 专业的向量检索引擎

这种架构既满足了业务需求，又保持了系统的简洁性和可维护性。移除 MongoDB 后，系统更加精简，运维成本降低，同时不影响核心功能的实现。
