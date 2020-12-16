package model

import (
	"github.com/Aoi-hosizora/ahlib-db/xgorm"
	"time"
)

type User struct {
	Id          uint32 `gorm:"primary_key;auto_increment"`
	ChatID      int64  `gorm:"not_null;unique_index:uk_chat_delete_at"`
	Username    string `gorm:"type:varchar(255);not_null"`
	Token       string `gorm:"type:varchar(255);not_null"`
	AllowIssue  bool   `gorm:"type:tinyint;not_null;default:0"`
	FilterMe    bool   `gorm:"type:tinyint;not_null;default:1"`
	Silent      bool   `gorm:"type:tinyint;non_null;default:0"`
	SilentStart int    `gorm:"non_null;default:0"`
	SilentEnd   int    `gorm:"non_null;default:0"`
	TimeZone    string `gorm:"non_null;default:'+00:00'"`

	DeletedAt *time.Time `gorm:"default:'1970-01-01 00:00:00';unique_index:uk_chat_delete_at"`
	xgorm.GormTimeWithoutDeleteAt
}
