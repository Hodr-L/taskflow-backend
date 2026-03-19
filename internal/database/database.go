package database

import (
	"fmt"
	"log"
	"time"

	"taskflow-backend/internal/config"
	"taskflow-backend/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// Connect 连接到数据库
func Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	if db != nil {
		return db, nil
	}

	// 使用简单的连接字符串
	// 使用MySQL格式的连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.Charset, cfg.ParseTime, cfg.Loc)

	// 配置Gorm日志
	gormLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// 连接数据库
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接测试失败: %w", err)
	}

	return db, nil
}

// AutoMigrate 自动迁移数据库表（仅用于开发环境）
func AutoMigrate(db *gorm.DB) error {
	// 按依赖顺序迁移表
	tables := []interface{}{
		&models.User{},
		&models.Team{},
		&models.TeamMember{},
		&models.Project{},
		&models.Task{},
		&models.TaskComment{},
		&models.Attachment{},
		&models.Notification{},
	}

	for _, table := range tables {
		if err := db.AutoMigrate(table); err != nil {
			return fmt.Errorf("迁移表失败: %w", err)
		}
	}

	return nil
}

// GetDB 返回数据库实例
func GetDB() *gorm.DB {
	if db == nil {
		panic("数据库未连接，请先调用 Connect()")
	}
	return db
}

// Close 关闭数据库连接
func Close() error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// Transaction 执行数据库事务
func Transaction(fn func(tx *gorm.DB) error) error {
	return db.Transaction(fn)
}

// CreateTables 手动创建表（备用方案）
func CreateTables(db *gorm.DB) error {
	// 创建用户表
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			avatar_url VARCHAR(255),
			role VARCHAR(20) DEFAULT 'user',
			status VARCHAR(20) DEFAULT 'active',
			email_verified BOOLEAN DEFAULT FALSE,
			last_login_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`).Error; err != nil {
		return fmt.Errorf("创建users表失败: %w", err)
	}

	// 创建团队表
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS teams (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			owner_id INTEGER NOT NULL,
			avatar_url VARCHAR(255),
			invite_code VARCHAR(50) UNIQUE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
		)
	`).Error; err != nil {
		return fmt.Errorf("创建teams表失败: %w", err)
	}

	// 创建团队成员表
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS team_members (
			id SERIAL PRIMARY KEY,
			team_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			role VARCHAR(20) DEFAULT 'member',
			joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(team_id, user_id),
			FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)
	`).Error; err != nil {
		return fmt.Errorf("创建team_members表失败: %w", err)
	}

	// 创建项目表
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			team_id INTEGER NOT NULL,
			owner_id INTEGER NOT NULL,
			status VARCHAR(20) DEFAULT 'active',
			color VARCHAR(7) DEFAULT '#3B82F6',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
			FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE SET NULL
		)
	`).Error; err != nil {
		return fmt.Errorf("创建projects表失败: %w", err)
	}

	// 创建任务表
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			title VARCHAR(200) NOT NULL,
			description TEXT,
			project_id INTEGER NOT NULL,
			assignee_id INTEGER,
			reporter_id INTEGER NOT NULL,
			priority VARCHAR(10) DEFAULT 'medium',
			status VARCHAR(20) DEFAULT 'todo',
			due_date TIMESTAMP WITH TIME ZONE,
			estimated_hours DECIMAL(5,2),
			actual_hours DECIMAL(5,2) DEFAULT 0,
			tags JSONB DEFAULT '[]',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			completed_at TIMESTAMP WITH TIME ZONE,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
			FOREIGN KEY (assignee_id) REFERENCES users(id) ON DELETE SET NULL,
			FOREIGN KEY (reporter_id) REFERENCES users(id) ON DELETE SET NULL
		)
	`).Error; err != nil {
		return fmt.Errorf("创建tasks表失败: %w", err)
	}

	// 创建索引
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_tasks_project_id ON tasks(project_id)",
		"CREATE INDEX IF NOT EXISTS idx_tasks_assignee_id ON tasks(assignee_id)",
		"CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status)",
		"CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks(due_date)",
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)",
		"CREATE INDEX IF NOT EXISTS idx_team_members_user_id ON team_members(user_id)",
	}

	for _, index := range indexes {
		if err := db.Exec(index).Error; err != nil {
			return fmt.Errorf("创建索引失败: %w", err)
		}
	}

	return nil
}
