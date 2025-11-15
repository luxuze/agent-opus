#!/bin/bash

# 启动后端服务
cd "$(dirname "$0")/backend"
echo "Starting backend server..."
go run cmd/server/main.go
