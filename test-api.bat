@echo off
echo 🔍 测试 TaskFlow API...

echo.
echo 1. 测试健康检查...
curl -s http://localhost:8080/api/v1/health | python -m json.tool

echo.
echo 2. 测试用户注册...
curl -s -X POST http://localhost:8080/api/v1/auth/register ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"testuser\",\"email\":\"test@example.com\",\"password\":\"password123\"}" ^
  | python -m json.tool

echo.
echo 3. 测试用户登录...
curl -s -X POST http://localhost:8080/api/v1/auth/login ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"test@example.com\",\"password\":\"password123\"}" ^
  | python -m json.tool

echo.
echo ✅ API 测试完成！
pause