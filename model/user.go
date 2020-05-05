package model

import (
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
	"time"
)

type User struct {
	Chat     string `gorm:"primary_key;auto_increment"`
	Private  bool   `gorm:"not_null;index:idx_unique_chat_delete_at"`
	Username string `gorm:"not_null"`
	Token    string `gorm:"not_null"`

	DeletedAt *time.Time `gorm:"default:'2000-01-01 00:00:00';index:idx_unique_chat_delete_at"`
	xgorm.GormTimeWithoutDeletedAt
}
