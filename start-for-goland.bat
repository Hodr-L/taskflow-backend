@echo off
echo 🚀 为 Goland 准备开发环境...

REM 设置环境变量（Goland运行配置中也需要设置）
set DB_HOST=localhost
set DB_PORT=3307
set DB_NAME=taskflow
set DB_USER=taskflow
set DB_PASSWORD=taskflow123
set REDIS_HOST=localhost
set REDIS_PORT=6379
set REDIS_PASSWORD=redis123
set APP_ENV=development

echo ✅ 环境变量已设置
echo.
echo 📋 在Goland中配置运行参数：
echo   工作目录: F:\openclaw\workspace\taskflow\backend
echo   运行包: taskflow-backend/cmd
echo   环境变量: 同上
echo.
echo 🐳 请确保Docker容器已启动：
echo   cd F:\openclaw\workspace\taskflow\deployments
echo   docker-compose -f docker-compose-simple.yml up -d
echo.
pause