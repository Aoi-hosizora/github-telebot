package model

import (
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
	"github.com/Aoi-hosizora/ahlib-web/xstatus"
	"time"
)

type User struct {
	Id         uint32 `gorm:"primary_key;auto_increment"`
	ChatID     int64  `gorm:"not_null;unique_index:uk_chat_delete_at"`
	Username   string `gorm:"type:varchar(255);not_null"`
	Token      string `gorm:"type:varchar(255);not_null"`
	AllowIssue bool   `gorm:"not_null;default:0"`

	DeletedAt *time.Time `gorm:"default:'2000-01-01 00:00:00';unique_index:uk_chat_delete_at"`
	xgorm.GormTimeWithoutDeletedAt
}

func GetUsers() []*User {
	users := make([]*User, 0)
	DB.Model(&User{}).Find(&users)
	return users
}

func GetUser(chatId int64) *User {
	user := &User{}
	rdb := DB.Model(&User{}).Where(&User{ChatID: chatId}).First(user)
	if rdb.RecordNotFound() {
		return nil
	}
	return user
}

func AddUser(user *User) xstatus.DbStatus {
	return xgorm.WithDB(DB).Insert(&User{}, user)
}

func UpdateUser(user *User) xstatus.DbStatus {
	return xgorm.WithDB(DB).Update(&User{}, &User{ChatID: user.ChatID}, user)
}

func DeleteUser(chatId int64) xstatus.DbStatus {
	return xgorm.WithDB(DB).Delete(&User{}, nil, &User{ChatID: chatId})
}
