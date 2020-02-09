package server

import (
	"encoding/json"
	"github.com/Aoi-hosizora/ah-tgbot/src/model"
	"github.com/Aoi-hosizora/ah-tgbot/src/util"
	"github.com/Aoi-hosizora/ahlib/xslice"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"time"
)

func Polling(server *BotServer) {
	bot := server.Bot
	config := server.Config

	chat, err := bot.ChatByID(config.TelegramConfig.ChannelId)
	if err != nil {
		log.Fatalf("Failed to get chat %s:%v\n", config.TelegramConfig.ChannelId, err)
	}
	log.Printf("Success to find channel \"%s\"\n", config.TelegramConfig.ChannelId)

	oldStr := ""
	oldObjs := make([]*model.GithubEvent, 0)
	dataCh := make(chan string)
	for {
		go func() {
			content, err := util.GetActions(config.GithubConfig, 1)
			if err != nil {
				dataCh <- ""
				return
			}
			dataCh <- content
		}()
		newStr := <-dataCh
		if newStr != "" && newStr != oldStr {
			newObjs := make([]*model.GithubEvent, 0)
			err := json.Unmarshal([]byte(newStr), &newObjs)
			if err == nil {
				diffItf := xslice.SliceDiff(xslice.Sti(newObjs), xslice.Sti(oldObjs))
				diff := xslice.Its(diffItf, &model.GithubEvent{}).([]*model.GithubEvent)
				if len(diff) != 0 { // new
					server.SendTo(chat, util.WrapGithubActions(diff), telebot.ModeMarkdown)
				}
				oldStr = newStr
				oldObjs = newObjs
			}
		}
		time.Sleep(time.Second * time.Duration(config.ServerConfig.PollingDuration))
	}
}
