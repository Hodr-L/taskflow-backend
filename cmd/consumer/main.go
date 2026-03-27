package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"taskflow-backend/internal/config"
	"taskflow-backend/pkg/logger"
)

type ConsumerType string

const (
	ConsumerNotification ConsumerType = "notification"
	ConsumerAudit        ConsumerType = "audit"
)

func main() {
	// 解析命令行参数
	consumerType := flag.String("type", "notification", "消费者类型: notification, audit")
	flag.Parse()

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ 加载配置失败: %v", err)
	}

	// 初始化日志
	logger.Init(cfg.Log.Level, cfg.Log.Format)

	logger.Info("🚀 启动 Kafka 消费者服务",
		zap.String("type", *consumerType),
		zap.String("env", cfg.App.Env),
	)

	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port,
		cfg.Database.Name, cfg.Database.Charset, cfg.Database.ParseTime, cfg.Database.Loc)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("❌ 数据库连接失败", logger.ErrorField(err))
	}

	// 根据消费者类型设置主题
	var topics []string
	var groupID string

	switch ConsumerType(*consumerType) {
	case ConsumerNotification:
		topics = []string{"user-registered", "task-created", "task-updated", "comment-added"}
		groupID = "notification-consumer-group"
		logger.Info("📨 启动通知消费者", zap.Strings("topics", topics))

	case ConsumerAudit:
		topics = []string{"user-login", "user-action", "task-operation", "file-upload"}
		groupID = "audit-consumer-group"
		logger.Info("📊 启动审计消费者", zap.Strings("topics", topics))

	default:
		logger.Fatal("❌ 未知的消费者类型", zap.String("type", *consumerType))
	}

	// 创建Kafka消费者
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Kafka.Brokers,
		GroupID:  groupID,
		Topic:    topics[0], // 暂时只监听一个主题，可以扩展为多个
		MinBytes: 10e3,      // 10KB
		MaxBytes: 10e6,      // 10MB
		MaxWait:  1 * time.Second,
	})

	defer reader.Close()

	logger.Info("✅ Kafka消费者已连接",
		zap.Strings("brokers", cfg.Kafka.Brokers),
		zap.String("topic", topics[0]),
		zap.String("group_id", groupID),
	)

	// 处理中断信号
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 消费消息
	for {
		select {
		case <-ctx.Done():
			logger.Info("🛑 收到停止信号，关闭消费者")
			return
		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				if err == context.Canceled {
					return
				}
				logger.Error("❌ 读取消息失败", logger.ErrorField(err))
				continue
			}

			// 处理消息
			processMessage(ConsumerType(*consumerType), msg, db)
		}
	}
}

func processMessage(consumerType ConsumerType, msg kafka.Message, db *gorm.DB) {
	logger.Info("📩 收到消息",
		zap.String("topic", msg.Topic),
		zap.Int("partition", msg.Partition),
		zap.Int64("offset", msg.Offset),
		zap.String("key", string(msg.Key)),
		zap.String("value", string(msg.Value)),
		zap.Time("time", msg.Time),
	)

	// 根据消费者类型处理消息
	switch consumerType {
	case ConsumerNotification:
		handleNotification(msg, db)
	case ConsumerAudit:
		handleAudit(msg, db)
	}
}

func handleNotification(msg kafka.Message, db *gorm.DB) {
	// TODO: 实现通知处理逻辑
	logger.Info("🔔 处理通知消息", zap.String("topic", msg.Topic))

	// 示例：将通知存入数据库
	// notification := models.Notification{
	//     UserID:  parseUserIdFromMessage(msg),
	//     Type:    msg.Topic,
	//     Title:   "新通知",
	//     Content: string(msg.Value),
	// }
	// db.Create(&notification)
}

func handleAudit(msg kafka.Message, db *gorm.DB) {
	// TODO: 实现审计日志处理逻辑
	logger.Info("📝 处理审计消息", zap.String("topic", msg.Topic))

	// 示例：将审计日志存入数据库
	// auditLog := models.AuditLog{
	//     Action:    msg.Topic,
	//     UserID:    parseUserIdFromMessage(msg),
	//     Details:   string(msg.Value),
	//     IPAddress: parseIpFromHeaders(msg.Headers),
	// }
	// db.Create(&auditLog)
}
