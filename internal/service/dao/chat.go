package dao

import (
	"github.com/Aoi-hosizora/ahlib-db/xgorm"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/Aoi-hosizora/github-telebot/internal/model"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/database"
)

func QueryChats() ([]*model.Chat, error) {
	chats := make([]*model.Chat, 0)
	rdb := database.GormDB().Model(&model.Chat{}).Find(&chats)
	if rdb.Error != nil {
		return nil, rdb.Error
	}
	return chats, nil
}

func QueryChat(chatID int64) (*model.Chat, error) {
	chat := &model.Chat{}
	rdb := database.GormDB().Model(&model.Chat{}).Where("chat_id = ?", chatID).First(chat)
	if rdb.RecordNotFound() {
		return nil, nil
	} else if rdb.Error != nil {
		return nil, rdb.Error
	}
	return chat, nil
}

func CreateChat(chatID int64, username, token string) (xstatus.DbStatus, error) {
	chat := &model.Chat{ChatID: chatID, Username: username, Token: token}
	rdb := database.GormDB().Model(&model.Chat{}).Create(chat)
	return xgorm.CreateErr(rdb)
}

func UpdateChatIssue(chatID int64, allow, filterMe bool) (xstatus.DbStatus, error) {
	rdb := database.GormDB().Model(&model.Chat{}).Where("chat_id = ?", chatID).Update(map[string]interface{}{
		"issue":     allow,
		"filter_me": filterMe,
	})
	return xgorm.UpdateErr(rdb)
}

func UpdateChatSilent(chatID int64, silent bool) (xstatus.DbStatus, error) {
	rdb := database.GormDB().Model(&model.Chat{}).Where("chat_id = ?", chatID).Update(map[string]interface{}{
		"silent": silent,
	})
	return xgorm.UpdateErr(rdb)
}

func DeleteChat(chatID int64) (xstatus.DbStatus, error) {
	rdb := database.GormDB().Model(&model.Chat{}).Where("chat_id = ?", chatID).Delete(&model.Chat{})
	return xgorm.DeleteErr(rdb)
}
