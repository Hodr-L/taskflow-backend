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
	// 瑙ｆ瀽鍛戒护琛屽弬鏁?	consumerType := flag.String("type", "notification", "娑堣垂鑰呯被鍨? notification, audit")
	flag.Parse()

	// 鍔犺浇閰嶇疆
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("鉂?鍔犺浇閰嶇疆澶辫触: %v", err)
	}

	// 鍒濆鍖栨棩蹇?	logger.Init(cfg.Log.Level, cfg.Log.Format)

	logger.Info("馃殌 鍚姩 Kafka 娑堣垂鑰呮湇鍔?,
		zap.String("type", *consumerType),
		zap.String("env", cfg.App.Env),
	)

	// 杩炴帴鏁版嵁搴?	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port,
		cfg.Database.Name, cfg.Database.Charset, cfg.Database.ParseTime, cfg.Database.Loc)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("鉂?鏁版嵁搴撹繛鎺ュけ璐?, logger.ErrorField(err))
	}

	// 鏍规嵁娑堣垂鑰呯被鍨嬭缃富棰?	var topics []string
	var groupID string

	switch ConsumerType(*consumerType) {
	case ConsumerNotification:
		topics = []string{"user-registered", "task-created", "task-updated", "comment-added"}
		groupID = "notification-consumer-group"
		logger.Info("馃摠 鍚姩閫氱煡娑堣垂鑰?, zap.Strings("topics", topics))

	case ConsumerAudit:
		topics = []string{"user-login", "user-action", "task-operation", "file-upload"}
		groupID = "audit-consumer-group"
		logger.Info("馃搳 鍚姩瀹¤娑堣垂鑰?, zap.Strings("topics", topics))

	default:
		logger.Fatal("鉂?鏈煡鐨勬秷璐硅€呯被鍨?, zap.String("type", *consumerType))
	}

	// 鍒涘缓Kafka娑堣垂鑰?	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Kafka.Brokers,
		GroupID:  groupID,
		Topic:    topics[0], // 鏆傛椂鍙洃鍚竴涓富棰橈紝鍙互鎵╁睍涓哄涓?		MinBytes: 10e3,      // 10KB
		MaxBytes: 10e6,      // 10MB
		MaxWait:  1 * time.Second,
	})

	defer reader.Close()

	logger.Info("鉁?Kafka娑堣垂鑰呭凡杩炴帴",
		zap.Strings("brokers", cfg.Kafka.Brokers),
		zap.String("topic", topics[0]),
		zap.String("group_id", groupID),
	)

	// 澶勭悊涓柇淇″彿
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 娑堣垂娑堟伅
	for {
		select {
		case <-ctx.Done():
			logger.Info("馃洃 鏀跺埌鍋滄淇″彿锛屽叧闂秷璐硅€?)
			return
		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				if err == context.Canceled {
					return
				}
				logger.Error("鉂?璇诲彇娑堟伅澶辫触", logger.ErrorField(err))
				continue
			}

			// 澶勭悊娑堟伅
			processMessage(ConsumerType(*consumerType), msg, db)
		}
	}
}

func processMessage(consumerType ConsumerType, msg kafka.Message, db *gorm.DB) {
	logger.Info("馃摡 鏀跺埌娑堟伅",
		zap.String("topic", msg.Topic),
		zap.Int("partition", msg.Partition),
		zap.Int64("offset", msg.Offset),
		zap.String("key", string(msg.Key)),
		zap.String("value", string(msg.Value)),
		zap.Time("time", msg.Time),
	)

	// 鏍规嵁娑堣垂鑰呯被鍨嬪鐞嗘秷鎭?	switch consumerType {
	case ConsumerNotification:
		handleNotification(msg, db)
	case ConsumerAudit:
		handleAudit(msg, db)
	}
}

func handleNotification(msg kafka.Message, db *gorm.DB) {
	// TODO: 瀹炵幇閫氱煡澶勭悊閫昏緫
	logger.Info("馃敂 澶勭悊閫氱煡娑堟伅", zap.String("topic", msg.Topic))

	// 绀轰緥锛氬皢閫氱煡瀛樺叆鏁版嵁搴?	// notification := models.Notification{
	//     UserID:  parseUserIdFromMessage(msg),
	//     Type:    msg.Topic,
	//     Title:   "鏂伴€氱煡",
	//     Content: string(msg.Value),
	// }
	// db.Create(&notification)
}

func handleAudit(msg kafka.Message, db *gorm.DB) {
	// TODO: 瀹炵幇瀹¤鏃ュ織澶勭悊閫昏緫
	logger.Info("馃摑 澶勭悊瀹¤娑堟伅", zap.String("topic", msg.Topic))

	// 绀轰緥锛氬皢瀹¤鏃ュ織瀛樺叆鏁版嵁搴?	// auditLog := models.AuditLog{
	//     Action:    msg.Topic,
	//     UserID:    parseUserIdFromMessage(msg),
	//     Details:   string(msg.Value),
	//     IPAddress: parseIpFromHeaders(msg.Headers),
	// }
	// db.Create(&auditLog)
}
