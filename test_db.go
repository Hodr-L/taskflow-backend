package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 测试不同的连接字符串格式
	dsnFormats := []string{
		// 格式1: 标准格式
		"host=localhost port=5432 user=taskflow password=taskflow123 dbname=taskflow sslmode=disable",
		// 格式2: 带TimeZone
		"host=localhost port=5432 user=taskflow password=taskflow123 dbname=taskflow sslmode=disable TimeZone=Asia/Shanghai",
		// 格式3: URL格式
		"postgres://taskflow:taskflow123@localhost:5432/taskflow?sslmode=disable",
		// 格式4: 带连接参数
		"host=localhost port=5432 user=taskflow password=taskflow123 dbname=taskflow sslmode=disable connect_timeout=10",
	}

	for i, dsn := range dsnFormats {
		fmt.Printf("\n=== 测试格式 %d ===\n", i+1)
		fmt.Printf("DSN: %s\n", dsn)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Printf("❌ 连接失败: %v\n", err)
			continue
		}

		sqlDB, err := db.DB()
		if err != nil {
			fmt.Printf("❌ 获取数据库实例失败: %v\n", err)
			continue
		}

		if err := sqlDB.Ping(); err != nil {
			fmt.Printf("❌ Ping失败: %v\n", err)
		} else {
			fmt.Printf("✅ 连接成功!\n")
		}

		sqlDB.Close()
	}
}