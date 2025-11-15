#!/bin/bash

echo "========================================="
echo "  Agent Platform API 测试"
echo "========================================="
echo ""

BASE_URL="http://localhost:8000"

# 测试健康检查
echo "1. 测试健康检查..."
curl -s $BASE_URL/health | python3 -m json.tool
echo ""

# 测试 Ping
echo "2. 测试 Ping..."
curl -s $BASE_URL/api/v1/ping | python3 -m json.tool
echo ""

# 获取 Agent 列表
echo "3. 获取 Agent 列表..."
curl -s $BASE_URL/api/v1/agents | python3 -m json.tool
echo ""

# 创建 Agent
echo "4. 创建新 Agent..."
curl -s -X POST $BASE_URL/api/v1/agents \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "测试 Agent",
    "description": "这是一个测试 Agent",
    "type": "single",
    "prompt_template": "You are a helpful assistant."
  }' | python3 -m json.tool
echo ""

# 获取工具列表
echo "5. 获取工具列表..."
curl -s $BASE_URL/api/v1/tools | python3 -m json.tool
echo ""

# 创建对话
echo "6. 创建对话..."
curl -s -X POST $BASE_URL/api/v1/conversations \
  -H 'Content-Type: application/json' \
  -d '{
    "agent_id": "agent-001",
    "title": "测试对话"
  }' | python3 -m json.tool
echo ""

echo "========================================="
echo "  所有测试完成"
echo "========================================="
