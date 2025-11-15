#!/bin/bash

echo "========================================="
echo "  Agent Platform 启动脚本"
echo "========================================="
echo ""

# 检查是否在正确的目录
if [ ! -d "backend" ] || [ ! -d "frontend" ]; then
    echo "错误: 请在项目根目录运行此脚本"
    exit 1
fi

# 启动后端
echo "📦 编译并启动后端服务..."
cd backend
go build -o agent-platform cmd/server/*.go
if [ $? -ne 0 ]; then
    echo "❌ 后端编译失败"
    exit 1
fi

# 后台运行后端
./agent-platform &
BACKEND_PID=$!
echo "✅ 后端已启动 (PID: $BACKEND_PID) - http://localhost:8000"

# 等待后端启动
sleep 2

# 测试后端
echo "🔍 测试后端健康检查..."
HEALTH=$(curl -s http://localhost:8000/health)
if [[ $HEALTH == *"healthy"* ]]; then
    echo "✅ 后端运行正常"
else
    echo "❌ 后端健康检查失败"
    kill $BACKEND_PID
    exit 1
fi

cd ../frontend

# 检查是否已安装依赖
if [ ! -d "node_modules" ]; then
    echo "📦 安装前端依赖..."
    npm install
fi

echo ""
echo "========================================="
echo "  服务已启动"
echo "========================================="
echo "  后端: http://localhost:8000"
echo "  后端 PID: $BACKEND_PID"
echo ""
echo "启动前端请运行:"
echo "  cd frontend && npm run dev"
echo ""
echo "测试 API:"
echo "  curl http://localhost:8000/health"
echo "  curl http://localhost:8000/api/v1/agents"
echo ""
echo "停止后端:"
echo "  kill $BACKEND_PID"
echo "========================================="
