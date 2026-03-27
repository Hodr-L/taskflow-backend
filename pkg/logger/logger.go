package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log *zap.Logger
)

// Init 初始化日志
func Init(level, format string) {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 设置输出格式
	var encoder zapcore.Encoder
	if format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 设置输出
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zapLevel,
	)

	// 创建logger
	log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	// 替换全局logger
	zap.ReplaceGlobals(log)
}

// customTimeEncoder 自定义时间格式
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// Debug 调试日志
func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

// Info 信息日志
func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// Warn 警告日志
func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

// Error 错误日志
func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

// Fatal 致命错误日志
func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

// Debugf 格式化调试日志
func Debugf(format string, args ...interface{}) {
	log.Sugar().Debugf(format, args...)
}

// Infof 格式化信息日志
func Infof(format string, args ...interface{}) {
	log.Sugar().Infof(format, args...)
}

// Warnf 格式化警告日志
func Warnf(format string, args ...interface{}) {
	log.Sugar().Warnf(format, args...)
}

// Errorf 格式化错误日志
func Errorf(format string, args ...interface{}) {
	log.Sugar().Errorf(format, args...)
}

// Fatalf 格式化致命错误日志
func Fatalf(format string, args ...interface{}) {
	log.Sugar().Fatalf(format, args...)
}

// With 添加字段
func With(fields ...zap.Field) *zap.Logger {
	return log.With(fields...)
}

// Sync 刷新日志缓冲区
func Sync() error {
	return log.Sync()
}

// GetLogger 获取logger实例
func GetLogger() *zap.Logger {
	return log
}

// 快捷方法
func Debugw(msg string, keysAndValues ...interface{}) {
	log.Sugar().Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	log.Sugar().Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	log.Sugar().Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	log.Sugar().Errorw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	log.Sugar().Fatalw(msg, keysAndValues...)
}

// HTTP请求日志字段
func HTTPRequest(method, path string, status int, duration time.Duration) zap.Field {
	return zap.Object("http", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		enc.AddString("method", method)
		enc.AddString("path", path)
		enc.AddInt("status", status)
		enc.AddDuration("duration", duration)
		return nil
	}))
}

// 数据库查询日志字段
func DBQuery(sql string, duration time.Duration, rows int64) zap.Field {
	return zap.Object("db", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		enc.AddString("query", sql)
		enc.AddDuration("duration", duration)
		enc.AddInt64("rows", rows)
		return nil
	}))
}

// 用户相关日志字段
func UserID(id uint) zap.Field {
	return zap.Uint("user_id", id)
}

func Username(name string) zap.Field {
	return zap.String("username", name)
}

func Email(email string) zap.Field {
	return zap.String("email", email)
}

// 错误日志字段
func ErrorField(err error) zap.Field {
	return zap.Error(err)
}

func ErrorString(err string) zap.Field {
	return zap.String("error", err)
}

// 业务相关字段
func TaskID(id uint) zap.Field {
	return zap.Uint("task_id", id)
}

func ProjectID(id uint) zap.Field {
	return zap.Uint("project_id", id)
}

func TeamID(id uint) zap.Field {
	return zap.Uint("team_id", id)
}
