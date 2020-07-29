package model

import (
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
	"time"
)

type User struct {
	Id       uint32 `gorm:"primary_key;auto_increment"`
	ChatID   int64  `gorm:"not_null;unique_index:uk_chat_delete_at"`
	Username string `gorm:"not_null"`
	Token    string `gorm:"not_null"`

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

func AddUser(user *User) DbStatus {
	rdb := DB.Model(&User{}).Create(user)
	if xgorm.IsMySqlDuplicateEntryError(rdb.Error) {
		return DbExisted
	} else if rdb.Error != nil || rdb.RowsAffected == 0 {
		return DbFailed
	}
	return DbSuccess
}

func UpdateUser(user *User) DbStatus {
	rdb := DB.Model(&User{}).Update(user)
	if xgorm.IsMySqlDuplicateEntryError(rdb.Error) {
		return DbExisted
	} else if rdb.Error != nil {
		return DbFailed
	} else if rdb.RowsAffected == 0 {
		return DbNotFound
	}
	return DbSuccess
}

func DeleteUser(chatId int64) DbStatus {
	rdb := DB.Model(&User{}).Delete(&User{ChatID: chatId})
	if rdb.Error != nil {
		return DbFailed
	} else if rdb.RowsAffected == 0 {
		return DbNotFound
	}
	return DbSuccess
}
