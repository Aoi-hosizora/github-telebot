package server

import (
	"encoding/json"
	"fmt"
	"github.com/Aoi-hosizora/ah-tgbot/src/model"
	"github.com/Aoi-hosizora/ah-tgbot/src/util"
	"gopkg.in/tucnak/telebot.v2"
	"strconv"
)

type HandleFunc = func(m *telebot.Message)

func setupRoute(server *BotServer) {
	bot := server.Bot

	bot.Handle("/start", start(server))
	bot.Handle("/github", github(server))

	bot.Handle(telebot.OnText, func(m *telebot.Message) {
		util.BotLog(m, nil, nil)
	})
}

func start(server *BotServer) HandleFunc {
	return func(m *telebot.Message) {
		server.Send(m, "This is AoiHosizora's bot")
	}
}

func github(server *BotServer) HandleFunc {
	return func(m *telebot.Message) {
		page, err := strconv.Atoi(m.Payload)
		if err != nil {
			server.Send(m, fmt.Sprintf("Please input a integer payload, \"%s\" is illegal.", m.Payload))
			return
		}
		content, err := util.GetActions(server.Config.GithubConfig, page)
		if err != nil {
			server.Send(m, "Failed to access github event: "+err.Error())
			return
		}
		objs := make([]*model.GithubEvent, 0)
		err = json.Unmarshal([]byte(content), &objs)
		if err != nil {
			server.Send(m, "Failed to unmarshal github event: "+err.Error())
			return
		}
		server.Send(m, util.WrapGithubActions(objs), telebot.ModeMarkdown)
	}
}
