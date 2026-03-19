@echo off
echo 🚀 启动 TaskFlow 后端服务...

REM 检查是否在正确的目录
if not exist "go.mod" (
    echo ❌ 错误：请在 backend 目录下运行此脚本
    pause
    exit /b 1
)

REM 检查 Docker 服务是否运行
docker info >nul 2>&1
if errorlevel 1 (
    echo ❌ Docker 未运行，请启动 Docker Desktop
    pause
    exit /b 1
)

REM 启动基础设施
echo 📦 启动 PostgreSQL 和 Redis...
cd ..\deployments
docker-compose up -d postgres redis

REM 等待数据库就绪
echo ⏳ 等待数据库就绪...
timeout /t 5 /nobreak >nul

REM 返回后端目录
cd ..\backend

REM 下载依赖
echo 📥 下载 Go 依赖...
go mod tidy

REM 启动服务
echo 🌐 启动后端服务...
echo 📡 API 地址: http://localhost:8080
echo 📊 健康检查: http://localhost:8080/api/v1/health
echo.
echo 按 Ctrl+C 停止服务
echo.

REM 使用 air 热重载
where air >nul 2>&1
if errorlevel 1 (
    echo ⚠️  air 未安装，使用 go run
    echo    安装 air: go install github.com/cosmtrek/air@latest
    go run cmd\main.go
) else (
    air
)