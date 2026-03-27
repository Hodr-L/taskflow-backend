# TaskFlow 鍚庣鏈嶅姟

## 馃搵 姒傝堪
TaskFlow 鐨勫悗绔?API 鏈嶅姟锛屼娇鐢?Go 璇█寮€鍙戯紝鍩轰簬 Gin 妗嗘灦锛屾彁渚涗换鍔＄鐞嗐€佸洟闃熷崗浣溿€佺敤鎴疯璇佺瓑鍔熻兘銆?

## 馃殌 蹇€熷紑濮?

### 鐜瑕佹眰
- Go 1.21+
- MySQL 8.0+
- Redis 7+
- (鍙€? Kafka 2.8+

### 1. 鍏嬮殕椤圭洰
```bash
git clone <repository-url>
cd taskflow/backend
```

### 2. 閰嶇疆鐜
```bash
# 澶嶅埗閰嶇疆鏂囦欢妯℃澘
cp config.example.yaml config.yaml
cp .env.example .env

# 缂栬緫 .env 鏂囦欢锛岃缃綘鐨勫瘑鐮佸拰瀵嗛挜
# 閲嶈锛氱敓浜х幆澧冨繀椤讳慨鏀归粯璁ゅ瘑鐮佸拰瀵嗛挜锛?
```

### 3. 瀹夎渚濊禆
```bash
go mod download
```

### 4. 鍚姩鍩虹璁炬柦
```bash
# 浣跨敤 Docker Compose 鍚姩鏁版嵁搴撳拰缂撳瓨
cd ../deployments
docker-compose up -d mysql redis
```

### 5. 杩愯鏈嶅姟
```bash
# 杩斿洖鍚庣鐩綍
cd ../backend

# 寮€鍙戞ā寮忥紙鐑噸杞斤級
air

# 鎴栫洿鎺ヨ繍琛?
go run cmd/main.go
```

### 6. 楠岃瘉鏈嶅姟
```bash
# 鍋ュ悍妫€鏌?
curl http://localhost:30088/api/v1/health

# 鎴栬闂?Swagger UI锛堝鏋滃凡鍚敤锛?
# http://localhost:30088/swagger/index.html
```

## 鈿欙笍 閰嶇疆绠＄悊

### 閰嶇疆鏂囦欢缁撴瀯
```
backend/
鈹溾攢鈹€ config.example.yaml      # 閰嶇疆妯℃澘锛堟彁浜ゅ埌浠撳簱锛?
鈹溾攢鈹€ .env.example            # 鐜鍙橀噺妯℃澘锛堟彁浜ゅ埌浠撳簱锛?
鈹溾攢鈹€ config.yaml             # 鏈湴閰嶇疆鏂囦欢锛?gitignore蹇界暐锛?
鈹溾攢鈹€ .env                    # 鏈湴鐜鍙橀噺锛?gitignore蹇界暐锛?
鈹斺攢鈹€ internal/config/config.go # 閰嶇疆鍔犺浇閫昏緫
```

### 閰嶇疆浼樺厛绾?
1. **鐜鍙橀噺**锛堟渶楂樹紭鍏堢骇锛?
2. **`.env` 鏂囦欢**锛堝紑鍙戠幆澧冩柟渚夸娇鐢級
3. **`config.yaml` 鏂囦欢**
4. **浠ｇ爜榛樿鍊?*锛堟渶浣庝紭鍏堢骇锛?

### 鐜鍙橀噺鍛藉悕
- **鏁版嵁搴?*: `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`
- **Redis**: `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`
- **JWT**: `JWT_SECRET`, `JWT_ACCESS_EXPIRE_HOURS`
- **鏈嶅姟鍣?*: `SERVER_PORT`, `SERVER_HOST`, `CORS_ALLOWED_ORIGINS`
- **搴旂敤**: `APP_ENV`, `APP_NAME`, `DEBUG`

### 寮€鍙戠幆澧冮厤缃ず渚?
```bash
# .env 鏂囦欢鍐呭绀轰緥
APP_ENV=development
DB_HOST=localhost
DB_PORT=53306
DB_NAME=taskflow
DB_USER=taskflow
DB_PASSWORD=taskflow123
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=redis123
JWT_SECRET=your-dev-jwt-secret-change-in-production
```

### 鐢熶骇鐜閰嶇疆瑕佹眰
1. **蹇呴』璁剧疆寮哄瘑鐮?*锛氭暟鎹簱鍜?Redis 瀵嗙爜闀垮害鑷冲皯 12 浣?
2. **蹇呴』璁剧疆寮?JWT Secret**锛氶暱搴﹁嚦灏?32 浣嶏紝浣跨敤闅忔満瀛楃涓?
3. **绂佺敤璋冭瘯妯″紡**锛氳缃?`DEBUG=false`
4. **璁剧疆姝ｇ‘鐨勭幆澧?*锛氳缃?`APP_ENV=production`

## 馃惓 Docker 閮ㄧ讲

### 鏋勫缓闀滃儚
```bash
docker build -t taskflow-backend:latest .
```

### 浣跨敤 Docker Compose
```bash
# 瀹屾暣鍫嗘爤锛堝悗绔?+ 鏁版嵁搴?+ Redis + Kafka锛?
cd ../deployments
docker-compose up -d

# 鎴栦粎鍚姩蹇呴渶鏈嶅姟
docker-compose up -d mysql redis backend
```

### 鐜鍙橀噺娉ㄥ叆
```bash
# 杩愯瀹瑰櫒鏃舵敞鍏ョ幆澧冨彉閲?
docker run -d \
  -p 30088:30088 \
  -e DB_PASSWORD=your_secure_password \
  -e JWT_SECRET=your_secure_jwt_secret \
  -e APP_ENV=production \
  taskflow-backend:latest
```

## 馃敡 寮€鍙戞寚鍗?

### 鏁版嵁搴撹縼绉?
```bash
# 鑷姩杩佺Щ锛堜粎寮€鍙戠幆澧冿級
go run scripts/migrate/main.go up

# 鎴栦娇鐢?GORM 鑷姩杩佺Щ
go run cmd/migrate/main.go
```

### API 娴嬭瘯
```bash
# 浣跨敤 curl 娴嬭瘯 API
curl -X POST http://localhost:30088/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 鎴栦娇鐢?Postman 瀵煎叆 API 鏂囨。
```

### 鏃ュ織鏌ョ湅
```bash
# 鏌ョ湅搴旂敤鏃ュ織
tail -f logs/app.log

# 鎴栨牴鎹厤缃緭鍑哄埌 stdout
```

## 馃攼 瀹夊叏閰嶇疆

### 鐢熶骇鐜瀹夊叏妫€鏌?
- [ ] 楠岃瘉 `JWT_SECRET` 宸茶缃笖闀垮害 鈮?32 浣?
- [ ] 楠岃瘉鏁版嵁搴撳瘑鐮佷笉鏄粯璁ゅ€?
- [ ] 楠岃瘉 Redis 瀵嗙爜涓嶆槸榛樿鍊?
- [ ] 纭 `DEBUG=false`
- [ ] 纭 CORS 閰嶇疆浠呭厑璁镐俊浠荤殑鍩熷悕
- [ ] 鍚敤 HTTPS锛堥€氳繃鍙嶅悜浠ｇ悊锛?

### 瀵嗛挜绠＄悊寤鸿
1. **浣跨敤瀵嗛挜绠＄悊鏈嶅姟**锛氬 HashiCorp Vault銆丄WS Secrets Manager
2. **瀹氭湡杞崲瀵嗛挜**锛氭瘡 3-6 涓湀杞崲涓€娆?JWT Secret
3. **鏈€灏忔潈闄愬師鍒?*锛氭暟鎹簱鐢ㄦ埛浠呮巿浜堝繀瑕佹潈闄?
4. **瀹¤鏃ュ織**锛氬惎鐢ㄦ墍鏈夋晱鎰熸搷浣滅殑瀹¤鏃ュ織

## 馃搳 鐩戞帶涓庡仴搴锋鏌?

### 鍋ュ悍妫€鏌ョ鐐?
- `GET /api/v1/health` - 搴旂敤鍋ュ悍鐘舵€?
- `GET /api/v1/health/db` - 鏁版嵁搴撹繛鎺ョ姸鎬?
- `GET /api/v1/health/redis` - Redis 杩炴帴鐘舵€?
- `GET /api/v1/health/kafka` - Kafka 杩炴帴鐘舵€侊紙濡傛灉鍚敤锛?

### 鐩戞帶鎸囨爣
- **Prometheus**锛氬鏋滃惎鐢紝鎸囨爣浣嶄簬 `/metrics`
- **鑷畾涔夋寚鏍?*锛氳姹傝鏁般€佸搷搴旀椂闂淬€侀敊璇巼绛?

## 馃毃 鏁呴殰鎺掗櫎

### 甯歌闂

#### 1. 鏁版嵁搴撹繛鎺ュけ璐?
```bash
# 妫€鏌?MySQL 鏈嶅姟鐘舵€?
docker ps | grep mysql

# 妫€鏌ョ鍙ｆ槸鍚﹁鍗犵敤
netstat -an | grep 53306

# 楠岃瘉鏁版嵁搴撻厤缃?
echo "DB_HOST=$DB_HOST, DB_PORT=$DB_PORT"
```

#### 2. Redis 杩炴帴澶辫触
```bash
# 妫€鏌?Redis 鏈嶅姟鐘舵€?
docker ps | grep redis

# 娴嬭瘯 Redis 杩炴帴
redis-cli -h localhost -p 6379 -a your_password ping
```

#### 3. JWT 楠岃瘉澶辫触
- 妫€鏌?`JWT_SECRET` 鐜鍙橀噺鏄惁璁剧疆
- 楠岃瘉浠ょ墝杩囨湡鏃堕棿閰嶇疆
- 纭绯荤粺鏃堕棿鍚屾

#### 4. 閰嶇疆鏂囦欢鍔犺浇闂
```bash
# 鍚敤璇︾粏鏃ュ織鏌ョ湅閰嶇疆鍔犺浇杩囩▼
DEBUG_CONFIG=true go run cmd/main.go
```

## 馃 璐＄尞鎸囧崡

1. Fork 椤圭洰浠撳簱
2. 鍒涘缓鍔熻兘鍒嗘敮 (`git checkout -b feature/amazing-feature`)
3. 鎻愪氦鏇存敼 (`git commit -m 'Add amazing feature'`)
4. 鎺ㄩ€佸埌鍒嗘敮 (`git push origin feature/amazing-feature`)
5. 鍒涘缓 Pull Request

## 馃搫 璁稿彲璇?

鏈」鐩噰鐢?MIT 璁稿彲璇?- 鏌ョ湅 [LICENSE](LICENSE) 鏂囦欢浜嗚В璇︽儏銆?

## 馃摓 鏀寔

- 闂鎶ュ憡锛歔GitHub Issues](https://github.com/your-org/taskflow/issues)
- 鏂囨。锛歔椤圭洰 Wiki](https://github.com/your-org/taskflow/wiki)
- 璁ㄨ锛歔GitHub Discussions](https://github.com/your-org/taskflow/discussions)

---

**鎻愮ず**锛氬缁堝湪鐢熶骇鐜鍓嶅湪娴嬭瘯鐜楠岃瘉閰嶇疆鏇存敼锛