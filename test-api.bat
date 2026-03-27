@echo off
echo 馃攳 娴嬭瘯 TaskFlow API...

echo.
echo 1. 娴嬭瘯鍋ュ悍妫€鏌?..
curl -s http://localhost:8080/api/v1/health | python -m json.tool

echo.
echo 2. 娴嬭瘯鐢ㄦ埛娉ㄥ唽...
curl -s -X POST http://localhost:8080/api/v1/auth/register ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"testuser\",\"email\":\"test@example.com\",\"password\":\"password123\"}" ^
  | python -m json.tool

echo.
echo 3. 娴嬭瘯鐢ㄦ埛鐧诲綍...
curl -s -X POST http://localhost:8080/api/v1/auth/login ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"test@example.com\",\"password\":\"password123\"}" ^
  | python -m json.tool

echo.
echo 鉁?API 娴嬭瘯瀹屾垚锛?pause