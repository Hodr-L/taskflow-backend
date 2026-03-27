@echo off
echo 馃殌 鍚姩 TaskFlow 鍚庣鏈嶅姟...

REM 妫€鏌ユ槸鍚﹀湪姝ｇ‘鐨勭洰褰?if not exist "go.mod" (
    echo 鉂?閿欒锛氳鍦?backend 鐩綍涓嬭繍琛屾鑴氭湰
    pause
    exit /b 1
)

REM 妫€鏌?Docker 鏈嶅姟鏄惁杩愯
docker info >nul 2>&1
if errorlevel 1 (
    echo 鉂?Docker 鏈繍琛岋紝璇峰惎鍔?Docker Desktop
    pause
    exit /b 1
)

REM 鍚姩鍩虹璁炬柦
echo 馃摝 鍚姩 PostgreSQL 鍜?Redis...
cd ..\deployments
docker-compose up -d postgres redis

REM 绛夊緟鏁版嵁搴撳氨缁?echo 鈴?绛夊緟鏁版嵁搴撳氨缁?..
timeout /t 5 /nobreak >nul

REM 杩斿洖鍚庣鐩綍
cd ..\backend

REM 涓嬭浇渚濊禆
echo 馃摜 涓嬭浇 Go 渚濊禆...
go mod tidy

REM 鍚姩鏈嶅姟
echo 馃寪 鍚姩鍚庣鏈嶅姟...
echo 馃摗 API 鍦板潃: http://localhost:8080
echo 馃搳 鍋ュ悍妫€鏌? http://localhost:8080/api/v1/health
echo.
echo 鎸?Ctrl+C 鍋滄鏈嶅姟
echo.

REM 浣跨敤 air 鐑噸杞?where air >nul 2>&1
if errorlevel 1 (
    echo 鈿狅笍  air 鏈畨瑁咃紝浣跨敤 go run
    echo    瀹夎 air: go install github.com/cosmtrek/air@latest
    go run cmd\main.go
) else (
    air
)