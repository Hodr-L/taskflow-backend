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

// Init 鍒濆鍖栨棩蹇?func Init(level, format string) {
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

	// 缂栫爜鍣ㄩ厤缃?	encoderConfig := zapcore.EncoderConfig{
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

	// 璁剧疆杈撳嚭鏍煎紡
	var encoder zapcore.Encoder
	if format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 璁剧疆杈撳嚭
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zapLevel,
	)

	// 鍒涘缓logger
	log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	// 鏇挎崲鍏ㄥ眬logger
	zap.ReplaceGlobals(log)
}

// customTimeEncoder 鑷畾涔夋椂闂存牸寮?func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// Debug 璋冭瘯鏃ュ織
func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

// Info 淇℃伅鏃ュ織
func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// Warn 璀﹀憡鏃ュ織
func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

// Error 閿欒鏃ュ織
func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

// Fatal 鑷村懡閿欒鏃ュ織
func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

// Debugf 鏍煎紡鍖栬皟璇曟棩蹇?func Debugf(format string, args ...interface{}) {
	log.Sugar().Debugf(format, args...)
}

// Infof 鏍煎紡鍖栦俊鎭棩蹇?func Infof(format string, args ...interface{}) {
	log.Sugar().Infof(format, args...)
}

// Warnf 鏍煎紡鍖栬鍛婃棩蹇?func Warnf(format string, args ...interface{}) {
	log.Sugar().Warnf(format, args...)
}

// Errorf 鏍煎紡鍖栭敊璇棩蹇?func Errorf(format string, args ...interface{}) {
	log.Sugar().Errorf(format, args...)
}

// Fatalf 鏍煎紡鍖栬嚧鍛介敊璇棩蹇?func Fatalf(format string, args ...interface{}) {
	log.Sugar().Fatalf(format, args...)
}

// With 娣诲姞瀛楁
func With(fields ...zap.Field) *zap.Logger {
	return log.With(fields...)
}

// Sync 鍒锋柊鏃ュ織缂撳啿鍖?func Sync() error {
	return log.Sync()
}

// GetLogger 鑾峰彇logger瀹炰緥
func GetLogger() *zap.Logger {
	return log
}

// 蹇嵎鏂规硶
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

// HTTP璇锋眰鏃ュ織瀛楁
func HTTPRequest(method, path string, status int, duration time.Duration) zap.Field {
	return zap.Object("http", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		enc.AddString("method", method)
		enc.AddString("path", path)
		enc.AddInt("status", status)
		enc.AddDuration("duration", duration)
		return nil
	}))
}

// 鏁版嵁搴撴煡璇㈡棩蹇楀瓧娈?func DBQuery(sql string, duration time.Duration, rows int64) zap.Field {
	return zap.Object("db", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		enc.AddString("query", sql)
		enc.AddDuration("duration", duration)
		enc.AddInt64("rows", rows)
		return nil
	}))
}

// 鐢ㄦ埛鐩稿叧鏃ュ織瀛楁
func UserID(id uint) zap.Field {
	return zap.Uint("user_id", id)
}

func Username(name string) zap.Field {
	return zap.String("username", name)
}

func Email(email string) zap.Field {
	return zap.String("email", email)
}

// 閿欒鏃ュ織瀛楁
func ErrorField(err error) zap.Field {
	return zap.Error(err)
}

func ErrorString(err string) zap.Field {
	return zap.String("error", err)
}

// 涓氬姟鐩稿叧瀛楁
func TaskID(id uint) zap.Field {
	return zap.Uint("task_id", id)
}

func ProjectID(id uint) zap.Field {
	return zap.Uint("project_id", id)
}

func TeamID(id uint) zap.Field {
	return zap.Uint("team_id", id)
}
