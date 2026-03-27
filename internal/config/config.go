package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Env     string `mapstructure:"env"`
	Debug   bool   `mapstructure:"debug"`
}

type ServerConfig struct {
	Port               int      `mapstructure:"port"`
	Host               string   `mapstructure:"host"`
	ReadTimeout        int      `mapstructure:"read_timeout"`
	WriteTimeout       int      `mapstructure:"write_timeout"`
	IdleTimeout        int      `mapstructure:"idle_timeout"`
	CORSAllowedOrigins []string `mapstructure:"cors_allowed_origins"`
}

type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Name            string `mapstructure:"name"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	Charset         string `mapstructure:"charset"`
	ParseTime       bool   `mapstructure:"parse_time"`
	Loc             string `mapstructure:"loc"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

type JWTConfig struct {
	Secret                 string `mapstructure:"secret"`
	AccessTokenExpireHours int    `mapstructure:"access_token_expire_hours"`
	RefreshTokenExpireDays int    `mapstructure:"refresh_token_expire_days"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

type UploadConfig struct {
	MaxSize      int64    `mapstructure:"max_size"`
	AllowedTypes []string `mapstructure:"allowed_types"`
	StoragePath  string   `mapstructure:"storage_path"`
}

type KafkaConfig struct {
	Enabled             bool                `mapstructure:"enabled"`
	Brokers             []string            `mapstructure:"brokers"`
	ClientID            string              `mapstructure:"client_id"`
	NotificationGroupID string              `mapstructure:"notification_group_id"`
	AuditGroupID        string              `mapstructure:"audit_group_id"`
	Topics              map[string][]string `mapstructure:"topics"`
}

var cfg *Config

func Load() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	// 加载.env文件（如果存在）
	if err := gotenv.Load(); err != nil {
		// .env文件不存在是正常的，不返回错误
		fmt.Println("ℹ️  .env文件未找到，将使用环境变量和默认配置")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")

	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("⚠️  配置文件未找到，使用默认值")
		} else {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	// 环境变量覆盖
	viper.AutomaticEnv()
	bindEnvVars()
	
	// 处理环境变量兼容性（DATABASE_ 前缀 -> DB_ 前缀）
	handleEnvCompatibility()

	// 解析配置到结构体
	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 验证配置
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	// 创建上传目录
	if err := os.MkdirAll(cfg.Upload.StoragePath, 0755); err != nil {
		return nil, fmt.Errorf("创建上传目录失败: %w", err)
	}

	return cfg, nil
}

func setDefaults() {
	viper.SetDefault("app.name", "TaskFlow")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.env", "development")
	viper.SetDefault("app.debug", true)

	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.read_timeout", 10)
	viper.SetDefault("server.write_timeout", 10)
	viper.SetDefault("server.idle_timeout", 60)
	viper.SetDefault("server.cors_allowed_origins", []string{"http://localhost:3000"})

	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.name", "taskflow")
	viper.SetDefault("database.user", "taskflow")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", 300)

	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)
	viper.SetDefault("redis.min_idle_conns", 2)

	viper.SetDefault("jwt.secret", "")
	viper.SetDefault("jwt.access_token_expire_hours", 24)
	viper.SetDefault("jwt.refresh_token_expire_days", 7)

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output", "stdout")

	viper.SetDefault("upload.max_size", 10485760) // 10MB
	viper.SetDefault("upload.allowed_types", []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"application/pdf",
	})
	viper.SetDefault("upload.storage_path", "./uploads")

	viper.SetDefault("kafka.enabled", false)
	viper.SetDefault("kafka.brokers", []string{"localhost:9092"})
	viper.SetDefault("kafka.client_id", "taskflow-backend")
	viper.SetDefault("kafka.group_id", "taskflow-consumer-group")
}

func bindEnvVars() {
	// 应用环境变量
	viper.BindEnv("app.name", "APP_NAME")
	viper.BindEnv("app.version", "APP_VERSION")
	viper.BindEnv("app.env", "APP_ENV")
	viper.BindEnv("app.debug", "APP_DEBUG")
	
	// 服务器环境变量
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("server.host", "SERVER_HOST")
	viper.BindEnv("server.read_timeout", "SERVER_READ_TIMEOUT")
	viper.BindEnv("server.write_timeout", "SERVER_WRITE_TIMEOUT")
	viper.BindEnv("server.idle_timeout", "SERVER_IDLE_TIMEOUT")
	viper.BindEnv("server.cors_allowed_origins", "CORS_ALLOWED_ORIGINS")

	// 数据库环境变量 - 使用 DB_ 前缀（与 Docker Compose 一致）
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.name", "DB_NAME")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.charset", "DB_CHARSET")
	viper.BindEnv("database.parse_time", "DB_PARSE_TIME")
	viper.BindEnv("database.loc", "DB_LOC")
	viper.BindEnv("database.max_open_conns", "DB_MAX_OPEN_CONNS")
	viper.BindEnv("database.max_idle_conns", "DB_MAX_IDLE_CONNS")
	viper.BindEnv("database.conn_max_lifetime", "DB_CONN_MAX_LIFETIME")
	
	// 数据库兼容性环境变量（旧 DATABASE_ 前缀）
	viper.BindEnv("database.host", "DATABASE_HOST")
	viper.BindEnv("database.port", "DATABASE_PORT")
	viper.BindEnv("database.name", "DATABASE_NAME")
	viper.BindEnv("database.user", "DATABASE_USER")
	viper.BindEnv("database.password", "DATABASE_PASSWORD")
	viper.BindEnv("database.charset", "DATABASE_CHARSET")
	viper.BindEnv("database.parse_time", "DATABASE_PARSE_TIME")
	viper.BindEnv("database.loc", "DATABASE_LOC")
	viper.BindEnv("database.max_open_conns", "DATABASE_MAX_OPEN_CONNS")
	viper.BindEnv("database.max_idle_conns", "DATABASE_MAX_IDLE_CONNS")
	viper.BindEnv("database.conn_max_lifetime", "DATABASE_CONN_MAX_LIFETIME")

	// Redis环境变量
	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("redis.db", "REDIS_DB")
	viper.BindEnv("redis.pool_size", "REDIS_POOL_SIZE")
	viper.BindEnv("redis.min_idle_conns", "REDIS_MIN_IDLE_CONNS")

	// JWT环境变量
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("jwt.access_token_expire_hours", "JWT_ACCESS_EXPIRE_HOURS")
	viper.BindEnv("jwt.refresh_token_expire_days", "JWT_REFRESH_EXPIRE_DAYS")

	// 日志环境变量
	viper.BindEnv("log.level", "LOG_LEVEL")
	viper.BindEnv("log.format", "LOG_FORMAT")
	viper.BindEnv("log.output", "LOG_OUTPUT")

	// 上传环境变量
	viper.BindEnv("upload.max_size", "UPLOAD_MAX_SIZE")
	viper.BindEnv("upload.allowed_types", "UPLOAD_ALLOWED_TYPES")
	viper.BindEnv("upload.storage_path", "UPLOAD_STORAGE_PATH")

	// Kafka环境变量
	viper.BindEnv("kafka.enabled", "KAFKA_ENABLED")
	viper.BindEnv("kafka.brokers", "KAFKA_BROKERS")
	viper.BindEnv("kafka.client_id", "KAFKA_CLIENT_ID")
	viper.BindEnv("kafka.notification_group_id", "KAFKA_NOTIFICATION_GROUP_ID")
	viper.BindEnv("kafka.audit_group_id", "KAFKA_AUDIT_GROUP_ID")
	// 注意：topics字段较复杂，通常通过配置文件设置
	
	// 邮件配置环境变量（可选）
	viper.BindEnv("smtp.host", "SMTP_HOST")
	viper.BindEnv("smtp.port", "SMTP_PORT")
	viper.BindEnv("smtp.user", "SMTP_USER")
	viper.BindEnv("smtp.password", "SMTP_PASSWORD")
	viper.BindEnv("smtp.from", "SMTP_FROM")
	viper.BindEnv("smtp.tls", "SMTP_TLS")
}

// Get 返回全局配置实例
func Get() *Config {
	if cfg == nil {
		panic("配置未初始化，请先调用 Load()")
	}
	return cfg
}

// GetDSN 返回数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.Charset, c.ParseTime, c.Loc)
}

// GetRedisAddr 返回Redis连接地址
func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetUploadPath 返回完整的上传路径
func (c *UploadConfig) GetUploadPath(filename string) string {
	return filepath.Join(c.StoragePath, filename)
}

// validateConfig 验证配置，确保关键配置项有效
func validateConfig(cfg *Config) error {
	// 设置开发环境的默认 JWT Secret（如果未设置）
	if cfg.JWT.Secret == "" {
		if cfg.App.Env == "production" {
			return fmt.Errorf("生产环境必须设置 JWT_SECRET 环境变量")
		}
		// 开发环境使用一个默认密钥（仅用于开发，生产环境必须更换）
		cfg.JWT.Secret = "dev-only-insecure-jwt-secret-change-in-production"
		fmt.Println("⚠️  使用开发环境默认 JWT Secret，生产环境必须设置 JWT_SECRET 环境变量")
	}

	// 环境特定的配置处理
	if cfg.App.Env == "production" {
		// 生产环境安全检查
		if cfg.Database.Password == "" || cfg.Database.Password == "taskflow123" {
			return fmt.Errorf("生产环境数据库密码不能为空或使用默认值")
		}
		if cfg.Redis.Password == "" || cfg.Redis.Password == "redis123" {
			return fmt.Errorf("生产环境Redis密码不能为空或使用默认值")
		}
		if cfg.JWT.Secret == "" || len(cfg.JWT.Secret) < 32 {
			return fmt.Errorf("生产环境JWT Secret长度至少32位")
		}
	} else {
		// 开发/测试环境：如果密码为空，设置默认值（仅用于开发）
		if cfg.Database.Password == "" {
			cfg.Database.Password = "taskflow123"
			fmt.Println("ℹ️  使用开发环境默认数据库密码，生产环境必须设置 DB_PASSWORD")
		}
		if cfg.Redis.Password == "" {
			cfg.Redis.Password = "redis123"
			fmt.Println("ℹ️  使用开发环境默认 Redis 密码，生产环境必须设置 REDIS_PASSWORD")
		}
	}

	// 验证端口范围
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("服务器端口%d无效，必须在1-65535之间", cfg.Server.Port)
	}

	// 验证数据库配置
	if cfg.Database.Host == "" {
		return fmt.Errorf("数据库主机地址不能为空")
	}
	if cfg.Database.Name == "" {
		return fmt.Errorf("数据库名称不能为空")
	}

	return nil
}

// handleEnvCompatibility 处理环境变量兼容性
func handleEnvCompatibility() {
	// 数据库配置兼容性：如果 DB_ 前缀未设置，但 DATABASE_ 前缀已设置，则使用 DATABASE_ 的值
	compatMap := map[string]string{
		"DB_HOST":                 "DATABASE_HOST",
		"DB_PORT":                 "DATABASE_PORT",
		"DB_NAME":                 "DATABASE_NAME",
		"DB_USER":                 "DATABASE_USER",
		"DB_PASSWORD":             "DATABASE_PASSWORD",
		"DB_CHARSET":              "DATABASE_CHARSET",
		"DB_PARSE_TIME":           "DATABASE_PARSE_TIME",
		"DB_LOC":                  "DATABASE_LOC",
		"DB_MAX_OPEN_CONNS":       "DATABASE_MAX_OPEN_CONNS",
		"DB_MAX_IDLE_CONNS":       "DATABASE_MAX_IDLE_CONNS",
		"DB_CONN_MAX_LIFETIME":    "DATABASE_CONN_MAX_LIFETIME",
	}

	for newKey, oldKey := range compatMap {
		newVal := os.Getenv(newKey)
		oldVal := os.Getenv(oldKey)
		
		// 如果新键未设置但旧键已设置，则设置到 viper 中
		if newVal == "" && oldVal != "" {
			// 将环境变量名转换为 viper 的键名
			viperKey := strings.ToLower(strings.ReplaceAll(newKey, "_", "."))
			// 根据数据类型转换值
			switch viperKey {
			case "db.port", "db.max.open.conns", "db.max.idle.conns", "db.conn.max.lifetime":
				if intVal, err := strconv.Atoi(oldVal); err == nil {
					viper.Set(viperKey, intVal)
				}
			case "db.parse.time":
				if boolVal, err := strconv.ParseBool(oldVal); err == nil {
					viper.Set(viperKey, boolVal)
				}
			default:
				viper.Set(viperKey, oldVal)
			}
		}
	}
	
	// 应用配置兼容性
	if os.Getenv("APP_ENV") != "" && viper.GetString("app.env") == "" {
		viper.Set("app.env", os.Getenv("APP_ENV"))
	}
}
