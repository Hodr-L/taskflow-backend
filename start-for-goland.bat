@echo off
echo 馃殌 涓?Goland 鍑嗗寮€鍙戠幆澧?..

REM 璁剧疆鐜鍙橀噺锛圙oland杩愯閰嶇疆涓篃闇€瑕佽缃級
set DB_HOST=localhost
set DB_PORT=3307
set DB_NAME=taskflow
set DB_USER=taskflow
set DB_PASSWORD=taskflow123
set REDIS_HOST=localhost
set REDIS_PORT=6379
set REDIS_PASSWORD=redis123
set APP_ENV=development

echo 鉁?鐜鍙橀噺宸茶缃?echo.
echo 馃搵 鍦℅oland涓厤缃繍琛屽弬鏁帮細
echo   宸ヤ綔鐩綍: F:\openclaw\workspace\taskflow\backend
echo   杩愯鍖? taskflow-backend/cmd
echo   鐜鍙橀噺: 鍚屼笂
echo.
echo 馃惓 璇风‘淇滵ocker瀹瑰櫒宸插惎鍔細
echo   cd F:\openclaw\workspace\taskflow\deployments
echo   docker-compose -f docker-compose-simple.yml up -d
echo.
pause