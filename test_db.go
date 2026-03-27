package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 娴嬭瘯涓嶅悓鐨勮繛鎺ュ瓧绗︿覆鏍煎紡
	dsnFormats := []string{
		// 鏍煎紡1: 鏍囧噯鏍煎紡
		"host=localhost port=5432 user=taskflow password=taskflow123 dbname=taskflow sslmode=disable",
		// 鏍煎紡2: 甯imeZone
		"host=localhost port=5432 user=taskflow password=taskflow123 dbname=taskflow sslmode=disable TimeZone=Asia/Shanghai",
		// 鏍煎紡3: URL鏍煎紡
		"postgres://taskflow:taskflow123@localhost:5432/taskflow?sslmode=disable",
		// 鏍煎紡4: 甯﹁繛鎺ュ弬鏁?		"host=localhost port=5432 user=taskflow password=taskflow123 dbname=taskflow sslmode=disable connect_timeout=10",
	}

	for i, dsn := range dsnFormats {
		fmt.Printf("\n=== 娴嬭瘯鏍煎紡 %d ===\n", i+1)
		fmt.Printf("DSN: %s\n", dsn)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Printf("鉂?杩炴帴澶辫触: %v\n", err)
			continue
		}

		sqlDB, err := db.DB()
		if err != nil {
			fmt.Printf("鉂?鑾峰彇鏁版嵁搴撳疄渚嬪け璐? %v\n", err)
			continue
		}

		if err := sqlDB.Ping(); err != nil {
			fmt.Printf("鉂?Ping澶辫触: %v\n", err)
		} else {
			fmt.Printf("鉁?杩炴帴鎴愬姛!\n")
		}

		sqlDB.Close()
	}
}
