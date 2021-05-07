package dao

import (
	"github.com/Aoi-hosizora/ahlib-db/xgorm"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/Aoi-hosizora/github-telebot/internal/model"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/database"
	"strings"
)

func QueryFilters(chatID int64) []*model.Filter {
	filters := make([]*model.Filter, 0)
	database.DB().Model(&model.Filter{}).Where("chat_id = ?", chatID).Where(&filters)
	return filters
}

func CheckFilter(chatID int64, username, repoName, eventType string) bool {
	if strings.Contains(repoName, "/") {
		repoName = strings.Split(repoName, "/")[1]
	}
	filters := make([]*model.Filter, 0)
	database.DB().Model(&model.Filter{}).Where("chat_id = ? AND username = ?", chatID, username).Find(&filters)

	for _, filter := range filters {
		if filter.RepoName == "*" {
			return true
		}
		if filter.RepoName == repoName {
			return false
		}
		if filter.EventType == "*" {
			return true
		}
		if filter.EventType == eventType {
			return false
		}
	}
	return false
}

func CreateFilter(chatID int64, username, repoName, eventType string) xstatus.DbStatus {
	rdb := database.DB().Model(&model.Filter{}).Create(&model.Filter{
		ChatID:    chatID,
		Username:  username,  // Actor.Login
		RepoName:  repoName,  // Repo.Name (trim username part)
		EventType: eventType, // xxxEvent or xxx_yyy
	})
	sts, _ := xgorm.CreateErr(rdb)
	return sts
}

func DeleteFilter(chatID int64, username, repoName, eventType string) xstatus.DbStatus {
	rdb := database.DB().Model(&model.Filter{}).Where("chat_id = ? AND username = ? AND repo_name = ? AND event_type = ?",
		chatID, username, repoName, eventType).Delete(&model.Filter{})
	sts, _ := xgorm.DeleteErr(rdb)
	return sts
}
