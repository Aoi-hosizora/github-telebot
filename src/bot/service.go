package bot

import (
	"github.com/Aoi-hosizora/github-telebot/src/logger"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"strconv"
)

func SendToChat(chatId int64, render string) {
	chat, err := Bot.ChatByID(strconv.Itoa(int(chatId)))
	if err != nil {
		log.Println(err)
		return
	}

	msg, err := Bot.Send(chat, render, telebot.ModeMarkdown)
	logger.SendLogger(chat, msg, err)
}
