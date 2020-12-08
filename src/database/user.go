package database

import (
	"github.com/Aoi-hosizora/ahlib-db/xgorm"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/Aoi-hosizora/github-telebot/src/model"
)

func GetUsers() []*model.User {
	users := make([]*model.User, 0)
	DB.Model(&model.User{}).Find(&users)
	return users
}

func GetUser(chatId int64) *model.User {
	user := &model.User{}
	rdb := DB.Model(&model.User{}).Where(&model.User{ChatID: chatId}).First(user)
	if rdb.RecordNotFound() {
		return nil
	}
	return user
}

func AddUser(user *model.User) xstatus.DbStatus {
	rdb := DB.Model(&model.User{}).Create(user)
	s, _ := xgorm.CreateErr(rdb)
	return s
}

func UpdateUser(user *model.User) xstatus.DbStatus {
	rdb := DB.Model(&model.User{}).Where(&model.User{ChatID: user.ChatID}).Updates(map[string]interface{}{
		"username":     user.Username,
		"token":        user.Token,
		"allow_issue":  user.AllowIssue,
		"silent":       user.Silent,
		"silent_start": user.SilentStart,
		"silent_end":   user.SilentEnd,
		"time_zone":    user.TimeZone,
	})
	s, _ := xgorm.UpdateErr(rdb)
	return s
}

func DeleteUser(chatId int64) xstatus.DbStatus {
	rdb := DB.Model(&model.User{}).Where(&model.User{ChatID: chatId}).Delete(&model.User{})
	s, _ := xgorm.DeleteErr(rdb)
	return s
}
