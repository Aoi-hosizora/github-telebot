package database

import (
	"github.com/Aoi-hosizora/ahlib-db/xgorm"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/Aoi-hosizora/github-telebot/internal/model"
)

func GetUsers() []*model.User {
	users := make([]*model.User, 0)
	DB().Model(&model.User{}).Find(&users)
	return users
}

func GetUser(chatId int64) *model.User {
	user := &model.User{}
	rdb := DB().Model(&model.User{}).Where(&model.User{ChatID: chatId}).First(user)
	if rdb.RecordNotFound() {
		return nil
	}
	return user
}

func AddUser(user *model.User) xstatus.DbStatus {
	rdb := DB().Model(&model.User{}).Create(user)
	s, _ := xgorm.CreateErr(rdb)
	return s
}

func UpdateUserAllowIssue(chatID int64, allowIssue, filterMe bool) xstatus.DbStatus {
	rdb := DB().Model(&model.User{}).Where(&model.User{ChatID: chatID}).Updates(map[string]interface{}{
		"allow_issue": allowIssue,
		"filter_me":   filterMe,
	})
	s, _ := xgorm.UpdateErr(rdb)
	return s
}

func UpdateUserSilent(chatID int64, silent bool, silentStart, silentEnd int, timeZone string) xstatus.DbStatus {
	rdb := DB().Model(&model.User{}).Where(&model.User{ChatID: chatID}).Updates(map[string]interface{}{
		"silent":       silent,
		"silent_start": silentStart,
		"silent_end":   silentEnd,
		"time_zone":    timeZone,
	})
	s, _ := xgorm.UpdateErr(rdb)
	return s
}

func DeleteUser(chatId int64) xstatus.DbStatus {
	rdb := DB().Model(&model.User{}).Where(&model.User{ChatID: chatId}).Delete(&model.User{})
	s, _ := xgorm.DeleteErr(rdb)
	return s
}
