# TaskFlow 后端服务

## 📋 概述
TaskFlow 的后端 API 服务，使用 Go 语言开发，基于 Gin 框架，提供任务管理、团队协作、用户认证等功能。

## 🚀 快速开始

### 环境要求
- Go 1.21+
- MySQL 8.0+
- Redis 7+
- (可选) Kafka 2.8+

### 1. 克隆项目
```bash
git clone <repository-url>
cd taskflow/backend
```

### 2. 配置环境
```bash
# 复制配置文件模板
cp config.example.yaml config.yaml
cp .env.example .env

# 编辑 .env 文件，设置你的密码和密钥
# 重要：生产环境必须修改默认密码和密钥！
```

### 3. 安装依赖
```bash
go mod download
```

### 4. 启动基础设施
```bash
# 使用 Docker Compose 启动数据库和缓存
cd ../deployments
docker-compose up -d mysql redis
```

### 5. 运行服务
```bash
# 返回后端目录
cd ../backend

# 开发模式（热重载）
air

# 或直接运行
go run cmd/main.go
```

### 6. 验证服务
```bash
# 健康检查
curl http://localhost:30088/api/v1/health

# 或访问 Swagger UI（如果已启用）
# http://localhost:30088/swagger/index.html
```

## ⚙️ 配置管理

### 配置文件结构
```
backend/
├── config.example.yaml      # 配置模板（提交到仓库）
├── .env.example            # 环境变量模板（提交到仓库）
├── config.yaml             # 本地配置文件（.gitignore忽略）
├── .env                    # 本地环境变量（.gitignore忽略）
└── internal/config/config.go # 配置加载逻辑
```

### 配置优先级
1. **环境变量**（最高优先级）
2. **`.env` 文件**（开发环境方便使用）
3. **`config.yaml` 文件**
4. **代码默认值**（最低优先级）

### 环境变量命名
- **数据库**: `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`
- **Redis**: `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`
- **JWT**: `JWT_SECRET`, `JWT_ACCESS_EXPIRE_HOURS`
- **服务器**: `SERVER_PORT`, `SERVER_HOST`, `CORS_ALLOWED_ORIGINS`
- **应用**: `APP_ENV`, `APP_NAME`, `DEBUG`

### 开发环境配置示例
```bash
# .env 文件内容示例
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

### 生产环境配置要求
1. **必须设置强密码**：数据库和 Redis 密码长度至少 12 位
2. **必须设置强 JWT Secret**：长度至少 32 位，使用随机字符串
3. **禁用调试模式**：设置 `DEBUG=false`
4. **设置正确的环境**：设置 `APP_ENV=production`

## 🐳 Docker 部署

### 构建镜像
```bash
docker build -t taskflow-backend:latest .
```

### 使用 Docker Compose
```bash
# 完整堆栈（后端 + 数据库 + Redis + Kafka）
cd ../deployments
docker-compose up -d

# 或仅启动必需服务
docker-compose up -d mysql redis backend
```

### 环境变量注入
```bash
# 运行容器时注入环境变量
docker run -d \
  -p 30088:30088 \
  -e DB_PASSWORD=your_secure_password \
  -e JWT_SECRET=your_secure_jwt_secret \
  -e APP_ENV=production \
  taskflow-backend:latest
```

## 🔧 开发指南

### 数据库迁移
```bash
# 自动迁移（仅开发环境）
go run scripts/migrate/main.go up

# 或使用 GORM 自动迁移
go run cmd/migrate/main.go
```

### API 测试
```bash
# 使用 curl 测试 API
curl -X POST http://localhost:30088/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 或使用 Postman 导入 API 文档
```

### 日志查看
```bash
# 查看应用日志
tail -f logs/app.log

# 或根据配置输出到 stdout
```

## 🔐 安全配置

### 生产环境安全检查
- [ ] 验证 `JWT_SECRET` 已设置且长度 ≥ 32 位
- [ ] 验证数据库密码不是默认值
- [ ] 验证 Redis 密码不是默认值
- [ ] 确认 `DEBUG=false`
- [ ] 确认 CORS 配置仅允许信任的域名
- [ ] 启用 HTTPS（通过反向代理）

### 密钥管理建议
1. **使用密钥管理服务**：如 HashiCorp Vault、AWS Secrets Manager
2. **定期轮换密钥**：每 3-6 个月轮换一次 JWT Secret
3. **最小权限原则**：数据库用户仅授予必要权限
4. **审计日志**：启用所有敏感操作的审计日志

## 📊 监控与健康检查

### 健康检查端点
- `GET /api/v1/health` - 应用健康状态
- `GET /api/v1/health/db` - 数据库连接状态
- `GET /api/v1/health/redis` - Redis 连接状态
- `GET /api/v1/health/kafka` - Kafka 连接状态（如果启用）

### 监控指标
- **Prometheus**：如果启用，指标位于 `/metrics`
- **自定义指标**：请求计数、响应时间、错误率等

## 🚨 故障排除

### 常见问题

#### 1. 数据库连接失败
```bash
# 检查 MySQL 服务状态
docker ps | grep mysql

# 检查端口是否被占用
netstat -an | grep 53306

# 验证数据库配置
echo "DB_HOST=$DB_HOST, DB_PORT=$DB_PORT"
```

#### 2. Redis 连接失败
```bash
# 检查 Redis 服务状态
docker ps | grep redis

# 测试 Redis 连接
redis-cli -h localhost -p 6379 -a your_password ping
```

#### 3. JWT 验证失败
- 检查 `JWT_SECRET` 环境变量是否设置
- 验证令牌过期时间配置
- 确认系统时间同步

#### 4. 配置文件加载问题
```bash
# 启用详细日志查看配置加载过程
DEBUG_CONFIG=true go run cmd/main.go
```

## 🤝 贡献指南

1. Fork 项目仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 📞 支持

- 问题报告：[GitHub Issues](https://github.com/your-org/taskflow/issues)
- 文档：[项目 Wiki](https://github.com/your-org/taskflow/wiki)
- 讨论：[GitHub Discussions](https://github.com/your-org/taskflow/discussions)

---

**提示**：始终在生产环境前在测试环境验证配置更改！