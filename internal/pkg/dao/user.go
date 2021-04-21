package dao

import (
	"github.com/Aoi-hosizora/ahlib-db/xgorm"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/Aoi-hosizora/github-telebot/internal/model"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/database"
)

func QueryUsers() []*model.User {
	users := make([]*model.User, 0)
	database.DB().Model(&model.User{}).Find(&users)
	return users
}

func QueryUser(chatID int64) *model.User {
	user := &model.User{}
	rdb := database.DB().Model(&model.User{}).Where("chat_id = ?", chatID).First(user)
	if rdb.RecordNotFound() {
		return nil
	}
	return user
}

func CreateUser(chatID int64, username, token string) xstatus.DbStatus {
	rdb := database.DB().Model(&model.User{}).Create(&model.User{
		ChatID:   chatID,
		Username: username,
		Token:    token,
	})
	sts, _ := xgorm.CreateErr(rdb)
	return sts
}

func UpdateUserAllowIssue(chatID int64, allowIssue, filterMe bool) xstatus.DbStatus {
	rdb := database.DB().Model(&model.User{}).Where("chat_id = ?", chatID).Update(map[string]interface{}{
		"allow_issue": allowIssue,
		"filter_me":   filterMe,
	})
	sts, _ := xgorm.UpdateErr(rdb)
	return sts
}

func UpdateUserSilent(chatID int64, silent bool, silentStart, silentEnd int, timeZone string) xstatus.DbStatus {
	rdb := database.DB().Model(&model.User{}).Where("chat_id = ?", chatID).Update(map[string]interface{}{
		"silent":       silent,
		"silent_start": silentStart,
		"silent_end":   silentEnd,
		"time_zone":    timeZone,
	})
	sts, _ := xgorm.UpdateErr(rdb)
	return sts
}

func DeleteUser(chatID int64) xstatus.DbStatus {
	rdb := database.DB().Model(&model.User{}).Where("chat_id = ?", chatID).Delete(&model.User{})
	sts, _ := xgorm.DeleteErr(rdb)
	return sts
}
