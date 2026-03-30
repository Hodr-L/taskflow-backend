package models

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 自定义 UUID 类型，支持 JSON 序列化和数据库存储
type UUID uuid.UUID

// 生成新的 UUID
func NewUUID() UUID {
	return UUID(uuid.New())
}

// 从字符串解析 UUID
func ParseUUID(s string) (UUID, error) {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return UUID{}, fmt.Errorf("invalid UUID format: %w", err)
	}
	return UUID(parsed), nil
}

// 转换为字符串
func (u UUID) String() string {
	return uuid.UUID(u).String()
}

// 是否为零值
func (u UUID) IsZero() bool {
	return uuid.UUID(u) == uuid.Nil
}

// 实现 json.Unmarshaler 接口
func (u *UUID) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "null" || s == "" {
		*u = UUID(uuid.Nil)
		return nil
	}

	parsed, err := uuid.Parse(s)
	if err != nil {
		return fmt.Errorf("invalid UUID format: %w", err)
	}
	*u = UUID(parsed)
	return nil
}

// 实现 json.Marshaler 接口
func (u UUID) MarshalJSON() ([]byte, error) {
	if u.IsZero() {
		return []byte(`null`), nil
	}
	return []byte(`"` + u.String() + `"`), nil
}

// 实现 binding 接口
func (u *UUID) UnmarshalText(text []byte) error {
	return u.UnmarshalJSON(text)
}

// 实现 driver.Valuer 接口（用于存储到数据库）
func (u UUID) Value() (driver.Value, error) {
	if u.IsZero() {
		return nil, nil
	}
	return u.String(), nil
}

// 实现 sql.Scanner 接口（用于从数据库读取）
func (u *UUID) Scan(value interface{}) error {
	if value == nil {
		*u = UUID(uuid.Nil)
		return nil
	}

	var str string
	switch v := value.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return fmt.Errorf("unsupported type for UUID: %T", value)
	}

	if str == "" {
		*u = UUID(uuid.Nil)
		return nil
	}

	parsed, err := uuid.Parse(str)
	if err != nil {
		return fmt.Errorf("invalid UUID format in database: %w", err)
	}
	*u = UUID(parsed)
	return nil
}

// 实现 Gorm 数据类型接口
func (UUID) GormDataType() string {
	return "char(36)"
}

// 实现 Gorm 表字段接口
func (UUID) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql", "mariadb":
		return "CHAR(36)"
	case "postgres":
		return "UUID"
	case "sqlite":
		return "TEXT"
	default:
		return "VARCHAR(36)"
	}
}

type RequestWithUUID struct {
	ID UUID `json:"id" binding:"required"`
}
