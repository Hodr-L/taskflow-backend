package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
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
	Enabled            bool                `mapstructure:"enabled"`
	Brokers            []string            `mapstructure:"brokers"`
	ClientID           string              `mapstructure:"client_id"`
	NotificationGroupID string             `mapstructure:"notification_group_id"`
	AuditGroupID       string              `mapstructure:"audit_group_id"`
	Topics             map[string][]string `mapstructure:"topics"`
}

var cfg *Config

func Load() (*Config, error) {
	if cfg != nil {
		return cfg, nil
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

	// 解析配置到结构体
	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
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
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.name", "taskflow")
	viper.SetDefault("database.user", "taskflow")
	viper.SetDefault("database.password", "taskflow123")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", 300)

	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "redis123")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)
	viper.SetDefault("redis.min_idle_conns", 2)

	viper.SetDefault("jwt.secret", "your-super-secret-jwt-key-change-in-production")
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
	// 数据库环境变量
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.name", "DB_NAME")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")

	// Redis环境变量
	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")

	// JWT环境变量
	viper.BindEnv("jwt.secret", "JWT_SECRET")

	// 应用环境变量
	viper.BindEnv("app.env", "APP_ENV")
	viper.BindEnv("app.debug", "APP_DEBUG")
	viper.BindEnv("server.port", "PORT")
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