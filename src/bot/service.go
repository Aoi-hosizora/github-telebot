package bot

import (
	"strconv"
)

func SendToChat(chatId int64, what interface{}, options ...interface{}) error {
	chat, err := Bot.Bot.ChatByID(strconv.FormatInt(chatId, 64))
	if err != nil {
		return err
	} else {
		return Bot.Send(chat, what, options...)
	}
}
