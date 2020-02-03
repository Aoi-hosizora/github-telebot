package main

import (
	"encoding/json"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xslice"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"strconv"
	"time"
)

const (
	GithubReceivedEventApi string = "https://api.github.com/users/%s/received_events"
)

func newBot(config *Config) *telebot.Bot {
	poller := &telebot.LongPoller{
		Timeout: time.Second * time.Duration(config.ServerConfig.PollerTimeout),
	}
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  config.TelegramConfig.BotToken,
		Poller: poller,
	})
	if err != nil {
		log.Fatalln("Failed to connect bot server:", err)
	}
	log.Printf("Success to connect bot server \"%s\"\n", bot.Me.Username)

	route(config, bot)
	bot.Handle(telebot.OnText, func(m *telebot.Message) {
		botLog(m, nil, err)
	})

	return bot
}

func route(config *Config, bot *telebot.Bot) {
	bot.Handle("/start", func(m *telebot.Message) {
		msg, err := bot.Send(m.Sender, "This is AoiHosizora's bot")
		botLog(m, msg, err)
	})
	bot.Handle("/hello", func(m *telebot.Message) {
		msg, err := bot.Send(m.Sender, "Hello world")
		botLog(m, msg, err)
	})
	bot.Handle("/test", func(m *telebot.Message) {
		msg, err := bot.Send(m.Sender, fmt.Sprintf("The payload is: \"%s\"", m.Payload))
		botLog(m, msg, err)
	})
	bot.Handle("/github", func(m *telebot.Message) {
		page, err := strconv.Atoi(m.Payload)
		if err != nil {
			msg, err := bot.Send(m.Sender, fmt.Sprintf("Please input a integer payload, \"%s\" is illegal.", m.Payload))
			botLog(m, msg, err)
			return
		}
		api := fmt.Sprintf(GithubReceivedEventApi, config.GithubConfig.Username)
		content, err := GithubUtil.httpGet(api, page, config.GithubConfig.Token)
		if err != nil {
			msg, err := bot.Send(m.Sender, "Failed to access github event: "+err.Error())
			botLog(m, msg, err)
			return
		}
		objs := make([]*GithubEvent, 0)
		err = json.Unmarshal([]byte(content), &objs)
		if err != nil {
			msg, err := bot.Send(m.Sender, "Failed to unmarshal github event: "+err.Error())
			botLog(m, msg, err)
			return
		}
		msg, err := bot.Send(m.Sender, GithubUtil.WrapGithubActions(objs), telebot.ModeMarkdown)
		botLog(m, msg, err)
	})
}

func polling(config *Config, bot *telebot.Bot) {
	chat, err := bot.ChatByID(config.TelegramConfig.ChannelId)
	if err != nil {
		log.Fatalf("Failed to get chat %s:%v\n", config.TelegramConfig.ChannelId, err)
	}
	log.Printf("Success to find channel \"%s\"\n", config.TelegramConfig.ChannelId)
	api := fmt.Sprintf(GithubReceivedEventApi, config.GithubConfig.Username)

	oldStr := ""
	oldObjs := make([]*GithubEvent, 0)
	dataCh := make(chan string)
	for {
		go func() {
			content, err := GithubUtil.httpGet(api, 1, config.GithubConfig.Token)
			if err != nil {
				dataCh <- ""
				return
			}
			dataCh <- content
		}()
		newStr := <-dataCh
		if newStr != "" && newStr != oldStr {
			newObjs := make([]*GithubEvent, 0)
			err := json.Unmarshal([]byte(newStr), &newObjs)
			if err == nil {
				diff := xslice.Its(xslice.SliceDiff(xslice.Sti(newObjs), xslice.Sti(oldObjs)), &GithubEvent{}).([]*GithubEvent)
				if len(diff) != 0 {
					msg, err := bot.Send(chat, GithubUtil.WrapGithubActions(diff), telebot.ModeMarkdown)
					botLog(nil, msg, err)
				}
				oldStr = newStr
				oldObjs = newObjs
			}
		}
		time.Sleep(time.Second * time.Duration(config.ServerConfig.PollingDuration))
	}
}
