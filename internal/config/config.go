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

	// 鍔犺浇.env鏂囦欢锛堝鏋滃瓨鍦級
	if err := gotenv.Load(); err != nil {
		// .env鏂囦欢涓嶅瓨鍦ㄦ槸姝ｅ父鐨勶紝涓嶈繑鍥為敊璇?		fmt.Println("鈩癸笍  .env鏂囦欢鏈壘鍒帮紝灏嗕娇鐢ㄧ幆澧冨彉閲忓拰榛樿閰嶇疆")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")

	// 璁剧疆榛樿鍊?	setDefaults()

	// 璇诲彇閰嶇疆鏂囦欢
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("鈿狅笍  閰嶇疆鏂囦欢鏈壘鍒帮紝浣跨敤榛樿鍊?)
		} else {
			return nil, fmt.Errorf("璇诲彇閰嶇疆鏂囦欢澶辫触: %w", err)
		}
	}

	// 鐜鍙橀噺瑕嗙洊
	viper.AutomaticEnv()
	bindEnvVars()

	// 澶勭悊鐜鍙橀噺鍏煎鎬э紙DATABASE_ 鍓嶇紑 -> DB_ 鍓嶇紑锛?	handleEnvCompatibility()

	// 瑙ｆ瀽閰嶇疆鍒扮粨鏋勪綋
	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("瑙ｆ瀽閰嶇疆澶辫触: %w", err)
	}

	// 楠岃瘉閰嶇疆
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("閰嶇疆楠岃瘉澶辫触: %w", err)
	}

	// 鍒涘缓涓婁紶鐩綍
	if err := os.MkdirAll(cfg.Upload.StoragePath, 0755); err != nil {
		return nil, fmt.Errorf("鍒涘缓涓婁紶鐩綍澶辫触: %w", err)
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
	// 搴旂敤鐜鍙橀噺
	viper.BindEnv("app.name", "APP_NAME")
	viper.BindEnv("app.version", "APP_VERSION")
	viper.BindEnv("app.env", "APP_ENV")
	viper.BindEnv("app.debug", "APP_DEBUG")

	// 鏈嶅姟鍣ㄧ幆澧冨彉閲?	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("server.host", "SERVER_HOST")
	viper.BindEnv("server.read_timeout", "SERVER_READ_TIMEOUT")
	viper.BindEnv("server.write_timeout", "SERVER_WRITE_TIMEOUT")
	viper.BindEnv("server.idle_timeout", "SERVER_IDLE_TIMEOUT")
	viper.BindEnv("server.cors_allowed_origins", "CORS_ALLOWED_ORIGINS")

	// 鏁版嵁搴撶幆澧冨彉閲?- 浣跨敤 DB_ 鍓嶇紑锛堜笌 Docker Compose 涓€鑷达級
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

	// 鏁版嵁搴撳吋瀹规€х幆澧冨彉閲忥紙鏃?DATABASE_ 鍓嶇紑锛?	viper.BindEnv("database.host", "DATABASE_HOST")
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

	// Redis鐜鍙橀噺
	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("redis.db", "REDIS_DB")
	viper.BindEnv("redis.pool_size", "REDIS_POOL_SIZE")
	viper.BindEnv("redis.min_idle_conns", "REDIS_MIN_IDLE_CONNS")

	// JWT鐜鍙橀噺
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("jwt.access_token_expire_hours", "JWT_ACCESS_EXPIRE_HOURS")
	viper.BindEnv("jwt.refresh_token_expire_days", "JWT_REFRESH_EXPIRE_DAYS")

	// 鏃ュ織鐜鍙橀噺
	viper.BindEnv("log.level", "LOG_LEVEL")
	viper.BindEnv("log.format", "LOG_FORMAT")
	viper.BindEnv("log.output", "LOG_OUTPUT")

	// 涓婁紶鐜鍙橀噺
	viper.BindEnv("upload.max_size", "UPLOAD_MAX_SIZE")
	viper.BindEnv("upload.allowed_types", "UPLOAD_ALLOWED_TYPES")
	viper.BindEnv("upload.storage_path", "UPLOAD_STORAGE_PATH")

	// Kafka鐜鍙橀噺
	viper.BindEnv("kafka.enabled", "KAFKA_ENABLED")
	viper.BindEnv("kafka.brokers", "KAFKA_BROKERS")
	viper.BindEnv("kafka.client_id", "KAFKA_CLIENT_ID")
	viper.BindEnv("kafka.notification_group_id", "KAFKA_NOTIFICATION_GROUP_ID")
	viper.BindEnv("kafka.audit_group_id", "KAFKA_AUDIT_GROUP_ID")
	// 娉ㄦ剰锛歵opics瀛楁杈冨鏉傦紝閫氬父閫氳繃閰嶇疆鏂囦欢璁剧疆

	// 閭欢閰嶇疆鐜鍙橀噺锛堝彲閫夛級
	viper.BindEnv("smtp.host", "SMTP_HOST")
	viper.BindEnv("smtp.port", "SMTP_PORT")
	viper.BindEnv("smtp.user", "SMTP_USER")
	viper.BindEnv("smtp.password", "SMTP_PASSWORD")
	viper.BindEnv("smtp.from", "SMTP_FROM")
	viper.BindEnv("smtp.tls", "SMTP_TLS")
}

// Get 杩斿洖鍏ㄥ眬閰嶇疆瀹炰緥
func Get() *Config {
	if cfg == nil {
		panic("閰嶇疆鏈垵濮嬪寲锛岃鍏堣皟鐢?Load()")
	}
	return cfg
}

// GetDSN 杩斿洖鏁版嵁搴撹繛鎺ュ瓧绗︿覆
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.Charset, c.ParseTime, c.Loc)
}

// GetRedisAddr 杩斿洖Redis杩炴帴鍦板潃
func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetUploadPath 杩斿洖瀹屾暣鐨勪笂浼犺矾寰?func (c *UploadConfig) GetUploadPath(filename string) string {
	return filepath.Join(c.StoragePath, filename)
}

// validateConfig 楠岃瘉閰嶇疆锛岀‘淇濆叧閿厤缃」鏈夋晥
func validateConfig(cfg *Config) error {
	// 璁剧疆寮€鍙戠幆澧冪殑榛樿 JWT Secret锛堝鏋滄湭璁剧疆锛?	if cfg.JWT.Secret == "" {
		if cfg.App.Env == "production" {
			return fmt.Errorf("鐢熶骇鐜蹇呴』璁剧疆 JWT_SECRET 鐜鍙橀噺")
		}
		// 寮€鍙戠幆澧冧娇鐢ㄤ竴涓粯璁ゅ瘑閽ワ紙浠呯敤浜庡紑鍙戯紝鐢熶骇鐜蹇呴』鏇存崲锛?		cfg.JWT.Secret = "dev-only-insecure-jwt-secret-change-in-production"
		fmt.Println("鈿狅笍  浣跨敤寮€鍙戠幆澧冮粯璁?JWT Secret锛岀敓浜х幆澧冨繀椤昏缃?JWT_SECRET 鐜鍙橀噺")
	}

	// 鐜鐗瑰畾鐨勯厤缃鐞?	if cfg.App.Env == "production" {
		// 鐢熶骇鐜瀹夊叏妫€鏌?		if cfg.Database.Password == "" || cfg.Database.Password == "taskflow123" {
			return fmt.Errorf("鐢熶骇鐜鏁版嵁搴撳瘑鐮佷笉鑳戒负绌烘垨浣跨敤榛樿鍊?)
		}
		if cfg.Redis.Password == "" || cfg.Redis.Password == "redis123" {
			return fmt.Errorf("鐢熶骇鐜Redis瀵嗙爜涓嶈兘涓虹┖鎴栦娇鐢ㄩ粯璁ゅ€?)
		}
		if cfg.JWT.Secret == "" || len(cfg.JWT.Secret) < 32 {
			return fmt.Errorf("鐢熶骇鐜JWT Secret闀垮害鑷冲皯32浣?)
		}
	} else {
		// 寮€鍙?娴嬭瘯鐜锛氬鏋滃瘑鐮佷负绌猴紝璁剧疆榛樿鍊硷紙浠呯敤浜庡紑鍙戯級
		if cfg.Database.Password == "" {
			cfg.Database.Password = "taskflow123"
			fmt.Println("鈩癸笍  浣跨敤寮€鍙戠幆澧冮粯璁ゆ暟鎹簱瀵嗙爜锛岀敓浜х幆澧冨繀椤昏缃?DB_PASSWORD")
		}
		if cfg.Redis.Password == "" {
			cfg.Redis.Password = "redis123"
			fmt.Println("鈩癸笍  浣跨敤寮€鍙戠幆澧冮粯璁?Redis 瀵嗙爜锛岀敓浜х幆澧冨繀椤昏缃?REDIS_PASSWORD")
		}
	}

	// 楠岃瘉绔彛鑼冨洿
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("鏈嶅姟鍣ㄧ鍙?d鏃犳晥锛屽繀椤诲湪1-65535涔嬮棿", cfg.Server.Port)
	}

	// 楠岃瘉鏁版嵁搴撻厤缃?	if cfg.Database.Host == "" {
		return fmt.Errorf("鏁版嵁搴撲富鏈哄湴鍧€涓嶈兘涓虹┖")
	}
	if cfg.Database.Name == "" {
		return fmt.Errorf("鏁版嵁搴撳悕绉颁笉鑳戒负绌?)
	}

	return nil
}

// handleEnvCompatibility 澶勭悊鐜鍙橀噺鍏煎鎬?func handleEnvCompatibility() {
	// 鏁版嵁搴撻厤缃吋瀹规€э細濡傛灉 DB_ 鍓嶇紑鏈缃紝浣?DATABASE_ 鍓嶇紑宸茶缃紝鍒欎娇鐢?DATABASE_ 鐨勫€?	compatMap := map[string]string{
		"DB_HOST":              "DATABASE_HOST",
		"DB_PORT":              "DATABASE_PORT",
		"DB_NAME":              "DATABASE_NAME",
		"DB_USER":              "DATABASE_USER",
		"DB_PASSWORD":          "DATABASE_PASSWORD",
		"DB_CHARSET":           "DATABASE_CHARSET",
		"DB_PARSE_TIME":        "DATABASE_PARSE_TIME",
		"DB_LOC":               "DATABASE_LOC",
		"DB_MAX_OPEN_CONNS":    "DATABASE_MAX_OPEN_CONNS",
		"DB_MAX_IDLE_CONNS":    "DATABASE_MAX_IDLE_CONNS",
		"DB_CONN_MAX_LIFETIME": "DATABASE_CONN_MAX_LIFETIME",
	}

	for newKey, oldKey := range compatMap {
		newVal := os.Getenv(newKey)
		oldVal := os.Getenv(oldKey)

		// 濡傛灉鏂伴敭鏈缃絾鏃ч敭宸茶缃紝鍒欒缃埌 viper 涓?		if newVal == "" && oldVal != "" {
			// 灏嗙幆澧冨彉閲忓悕杞崲涓?viper 鐨勯敭鍚?			viperKey := strings.ToLower(strings.ReplaceAll(newKey, "_", "."))
			// 鏍规嵁鏁版嵁绫诲瀷杞崲鍊?			switch viperKey {
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

	// 搴旂敤閰嶇疆鍏煎鎬?	if os.Getenv("APP_ENV") != "" && viper.GetString("app.env") == "" {
		viper.Set("app.env", os.Getenv("APP_ENV"))
	}
}
