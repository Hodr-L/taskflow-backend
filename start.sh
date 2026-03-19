#!/bin/bash

echo "🚀 启动 TaskFlow 后端服务..."

# 检查是否在正确的目录
if [ ! -f "go.mod" ]; then
    echo "❌ 错误：请在 backend 目录下运行此脚本"
    exit 1
fi

# 检查 Docker 服务是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker 未运行，请启动 Docker Desktop"
    exit 1
fi

# 启动基础设施
echo "📦 启动 PostgreSQL 和 Redis..."
cd ../deployments
docker-compose up -d postgres redis

# 等待数据库就绪
echo "⏳ 等待数据库就绪..."
sleep 5

# 返回后端目录
cd ../backend

# 下载依赖
echo "📥 下载 Go 依赖..."
go mod tidy

# 启动服务
echo "🌐 启动后端服务..."
echo "📡 API 地址: http://localhost:8080"
echo "📊 健康检查: http://localhost:8080/api/v1/health"
echo ""
echo "按 Ctrl+C 停止服务"

# 使用 air 热重载
if command -v air &> /dev/null; then
    air
else
    echo "⚠️  air 未安装，使用 go run"
    echo "   安装 air: go install github.com/cosmtrek/air@latest"
    go run cmd/main.go
fi