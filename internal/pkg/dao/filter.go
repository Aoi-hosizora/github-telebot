package dao

import (
	"github.com/Aoi-hosizora/ahlib-db/xgorm"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/Aoi-hosizora/github-telebot/internal/model"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/database"
)

func QueryFilters(chatID int64) []*model.Filter {
	filters := make([]*model.Filter, 0)
	database.DB().Model(&model.Filter{}).Where("chat_id = ?", chatID).Where(&filters)
	return filters
}

func CheckFilter(chatID int64, isActivity bool, repoName, eventType string) bool {
	count := 0
	if eventType == "*" {
		database.DB().Model(&model.Filter{}).Where("chat_id = ? AND is_activity = ? AND repo_name = ?",
			chatID, isActivity, repoName).Count(&count)
	} else {
		database.DB().Model(&model.Filter{}).Where("chat_id = ? AND is_activity = ? AND repo_name = ? AND event_type = ?",
			chatID, isActivity, repoName, eventType).Count(&count)
	}
	return count > 0
}

func CreateFilter(chatID int64, isActivity bool, repoName, eventType string) xstatus.DbStatus {
	rdb := database.DB().Model(&model.Filter{}).Create(&model.Filter{
		ChatID:     chatID,
		IsActivity: isActivity,
		RepoName:   repoName,
		EventType:  eventType,
	})
	sts, _ := xgorm.CreateErr(rdb)
	return sts
}

func DeleteFilter(chatID int64, isActivity bool, repoName, eventType string) xstatus.DbStatus {
	rdb := database.DB().Model(&model.Filter{}).Where("chat_id = ? AND is_activity = ? AND repo_name = ? AND event_type = ?",
		chatID, isActivity, repoName, eventType).Delete(&model.Filter{})
	sts, _ := xgorm.DeleteErr(rdb)
	return sts
}
