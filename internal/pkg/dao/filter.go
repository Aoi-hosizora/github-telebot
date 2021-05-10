package dao

import (
	"github.com/Aoi-hosizora/ahlib-db/xgorm"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/Aoi-hosizora/github-telebot/internal/model"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/database"
)

func QueryFilters(chatID int64) []*model.Filter {
	filters := make([]*model.Filter, 0)
	database.DB().Model(&model.Filter{}).Where("chat_id = ?", chatID).Find(&filters)
	return filters
}

func CheckFilter(chatID int64, repoName, actor, eventType string) bool {
	filters := make([]*model.Filter, 0)
	database.DB().Model(&model.Filter{}).Where("chat_id = ? AND repo_name = ?", chatID, repoName).Find(&filters)
	for _, filter := range filters {
		if filter.Actor == "*" || filter.Actor == actor {
			return true
		}
		if filter.EventType == "*" || filter.EventType == eventType {
			return true
		}
	}
	return false
}

func CreateFilter(chatID int64, repoName, actor, eventType string) xstatus.DbStatus {
	rdb := database.DB().Model(&model.Filter{}).Create(&model.Filter{
		ChatID:    chatID,
		Actor:     actor,
		RepoName:  repoName,
		EventType: eventType,
	})
	sts, _ := xgorm.CreateErr(rdb)
	return sts
}

func DeleteFilter(chatID int64, repoName, actor, eventType string) xstatus.DbStatus {
	rdb := database.DB().Model(&model.Filter{}).Where("chat_id = ? AND actor = ? AND repo_name = ? AND event_type = ?",
		chatID, actor, repoName, eventType).Delete(&model.Filter{})
	sts, _ := xgorm.DeleteErr(rdb)
	return sts
}
