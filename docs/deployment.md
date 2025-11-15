# Agent 平台部署文档

## 部署架构

```
┌─────────────────────────────────────────────────────────┐
│                      Load Balancer                       │
└─────────────────────────────────────────────────────────┘
                            │
        ┌───────────────────┴───────────────────┐
        │                                       │
┌───────▼────────┐                    ┌────────▼───────┐
│   Frontend     │                    │    Backend     │
│   (Nginx)      │◄──────────────────►│     (Go)       │
└────────────────┘                    └────────────────┘
                                              │
                    ┌─────────────────────────┼─────────────────────┐
                    │                         │                     │
            ┌───────▼────────┐      ┌────────▼────────┐  ┌─────────▼────────┐
            │     MySQL      │      │    MongoDB      │  │      Redis       │
            └────────────────┘      └─────────────────┘  └──────────────────┘
```

## 环境要求

### 最低配置

- CPU: 2 核
- 内存: 4GB
- 磁盘: 20GB
- 操作系统: Linux (Ubuntu 20.04+ / CentOS 7+)

### 推荐配置

- CPU: 4 核+
- 内存: 8GB+
- 磁盘: 50GB+
- 操作系统: Ubuntu 22.04 LTS

## 部署方式

### 1. Docker Compose 部署 (推荐)

#### 1.1 安装 Docker 和 Docker Compose

```bash
# Ubuntu
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 安装 Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

#### 1.2 配置环境

```bash
# 克隆项目
git clone https://github.com/yourname/agent-platform.git
cd agent-platform

# 配置后端环境变量
cp backend/.env.example backend/.env
vim backend/.env
```

必须修改的配置项：

- `JWT_SECRET`: 生产环境的 JWT 密钥
- `MYSQL_PASSWORD`: MySQL 密码
- `OPENAI_API_KEY`: OpenAI API Key (如使用)
- `CORS_ORIGINS`: 允许的跨域源

#### 1.3 启动服务

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

#### 1.4 验证部署

```bash
# 检查后端健康状态
curl http://localhost:8000/health

# 访问前端
http://localhost:3000
```

### 2. 手动部署

#### 2.1 部署数据库

##### MySQL

```bash
# 安装 MySQL
sudo apt update
sudo apt install mysql-server

# 创建数据库和用户
mysql -u root -p
CREATE DATABASE agent_platform;
CREATE USER 'agent_user'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON agent_platform.* TO 'agent_user'@'localhost';
FLUSH PRIVILEGES;
```

##### MongoDB

```bash
# 安装 MongoDB
wget -qO - https://www.mongodb.org/static/pgp/server-6.0.asc | sudo apt-key add -
echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/6.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-6.0.list
sudo apt update
sudo apt install -y mongodb-org
sudo systemctl start mongod
sudo systemctl enable mongod
```

##### Redis

```bash
# 安装 Redis
sudo apt install redis-server
sudo systemctl start redis
sudo systemctl enable redis
```

#### 2.2 部署后端

```bash
# 安装 Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 构建后端
cd backend
cp .env.example .env
vim .env  # 配置环境变量
go mod download
go build -o agent-platform ./cmd/server

# 使用 systemd 管理服务
sudo vim /etc/systemd/system/agent-platform.service
```

systemd 配置文件内容：

```ini
[Unit]
Description=Agent Platform Backend
After=network.target mysql.service mongodb.service redis.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/agent-platform/backend
ExecStart=/opt/agent-platform/backend/agent-platform
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl start agent-platform
sudo systemctl enable agent-platform
sudo systemctl status agent-platform
```

#### 2.3 部署前端

```bash
# 安装 Node.js
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs

# 构建前端
cd frontend
npm install
npm run build

# 安装 Nginx
sudo apt install nginx

# 配置 Nginx
sudo vim /etc/nginx/sites-available/agent-platform
```

Nginx 配置：

```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /opt/agent-platform/frontend/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml;
}
```

启用配置：

```bash
sudo ln -s /etc/nginx/sites-available/agent-platform /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### 3. 生产环境优化

#### 3.1 启用 HTTPS

```bash
# 使用 Let's Encrypt
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

#### 3.2 配置防火墙

```bash
sudo ufw allow 22
sudo ufw allow 80
sudo ufw allow 443
sudo ufw enable
```

#### 3.3 数据库优化

MySQL 配置 (`/etc/mysql/mysql.conf.d/mysqld.cnf`):

```ini
[mysqld]
max_connections = 200
innodb_buffer_pool_size = 1G
innodb_log_file_size = 256M
query_cache_type = 1
query_cache_size = 64M
```

#### 3.4 性能监控

安装 Prometheus 和 Grafana：

```bash
# Prometheus
docker run -d -p 9090:9090 --name prometheus prom/prometheus

# Grafana
docker run -d -p 3001:3000 --name grafana grafana/grafana
```

## 备份策略

### 数据库备份

```bash
# MySQL 备份脚本
#!/bin/bash
BACKUP_DIR="/backup/mysql"
DATE=$(date +%Y%m%d_%H%M%S)
mysqldump -u agent_user -p agent_platform > $BACKUP_DIR/backup_$DATE.sql
find $BACKUP_DIR -name "*.sql" -mtime +7 -delete

# 添加到 crontab
0 2 * * * /path/to/backup.sh
```

### MongoDB 备份

```bash
# MongoDB 备份
mongodump --db agent_platform --out /backup/mongodb/$(date +%Y%m%d)
```

## 故障排查

### 常见问题

1. **后端无法连接数据库**

```bash
# 检查数据库服务状态
sudo systemctl status mysql
sudo systemctl status mongodb
sudo systemctl status redis

# 检查网络连接
telnet localhost 3306
telnet localhost 27017
telnet localhost 6379
```

2. **前端无法访问**

```bash
# 检查 Nginx 状态
sudo systemctl status nginx

# 检查 Nginx 错误日志
sudo tail -f /var/log/nginx/error.log
```

3. **后端日志查看**

```bash
# systemd 服务日志
sudo journalctl -u agent-platform -f

# Docker 日志
docker-compose logs -f backend
```

## 扩展部署

### 负载均衡

使用 Nginx 作为负载均衡器：

```nginx
upstream backend {
    server backend1:8000;
    server backend2:8000;
    server backend3:8000;
}

server {
    location /api {
        proxy_pass http://backend;
    }
}
```

### 数据库主从复制

参考 MySQL 官方文档配置主从复制以提高可用性。

## 监控告警

### 监控指标

- CPU 使用率
- 内存使用率
- 磁盘使用率
- API 响应时间
- 数据库连接数
- 错误率

### 告警配置

使用 Prometheus Alertmanager 配置告警规则。

## 安全加固

1. **定期更新系统**

```bash
sudo apt update
sudo apt upgrade
```

2. **配置防火墙规则**
3. **定期备份数据**
4. **监控异常访问**
5. **使用强密码策略**
6. **定期安全审计**

## 版本升级

### 后端升级

```bash
# 停止服务
docker-compose stop backend

# 拉取新代码
git pull origin main

# 重新构建
docker-compose build backend

# 启动服务
docker-compose up -d backend
```

### 前端升级

```bash
docker-compose stop frontend
git pull origin main
docker-compose build frontend
docker-compose up -d frontend
```

## 回滚策略

保留最近 3 个版本的备份，如遇问题可快速回滚：

```bash
# 恢复数据库
mysql -u agent_user -p agent_platform < /backup/mysql/backup_20240115.sql

# 回滚代码
git checkout <previous-commit>
docker-compose up -d --build
```

## 支持

如遇部署问题，请查看：

- 项目 Issues: https://github.com/yourname/agent-platform/issues
- 文档: https://docs.your-domain.com
