package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"
)

const (
	GithubReceivedEventApi string = "https://api.github.com/users/%s/received_events"
)

func newBot(config *Config) *telebot.Bot {
	poller := &telebot.LongPoller{
		Timeout: 2 * time.Second,
	}
	mwPoller := telebot.NewMiddlewarePoller(poller, func(update *telebot.Update) bool {
		// update.Message.Text
		botLog(update.Message, nil, nil)
		return true
	})

	bot, err := telebot.NewBot(telebot.Settings{
		Token:  config.TelegramConfig.BotToken,
		Poller: mwPoller,
	})
	if err != nil {
		log.Fatalln("Failed to connect bot server:", err)
	}
	log.Printf("Success to connect bot server \"%s\"\n", bot.Me.Username)

	route(bot)
	return bot
}

func route(bot *telebot.Bot) {
	bot.Handle("/start", func(m *telebot.Message) {
		msg, err := bot.Send(m.Sender, "This is AoiHosizora's bot")
		botLog(m, msg, err)
	})
	bot.Handle("/hello", func(m *telebot.Message) {
		msg, err := bot.Send(m.Sender, "Hello world")
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
			req, err := http.NewRequest("GET", api, nil)
			if err != nil {
				dataCh <- ""
				return
			}
			req.Header.Add("Authorization", fmt.Sprintf("Token %s", config.GithubConfig.Token))
			resp, err := (&http.Client{}).Do(req)
			if err != nil {
				dataCh <- ""
				return
			}
			defer resp.Body.Close()

			content, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				dataCh <- ""
				return
			}
			dataCh <- string(content)
		}()
		newStr := <-dataCh
		if newStr != "" && newStr != oldStr {
			newObjs := make([]*GithubEvent, 0)
			err := json.Unmarshal([]byte(newStr), &newObjs)
			if err == nil {
				diff := sliceDiff(newObjs, oldObjs)
				if len(diff) != 0 {
					msg, err := bot.Send(chat, GithubUtil.WrapGithubActions(diff))
					botLog(nil, msg, err)
				}
				oldStr = newStr
				oldObjs = newObjs
			}
		}
		time.Sleep(time.Second * time.Duration(config.ServerConfig.PollingDuration))
	}
}

// [{"type":"WatchEvent","actor":{"login":"iamcco","display_login":"iamcco","url":"https://api.github.com/users/iamcco"},"repo":{"name":"shanyuhai123/learnCSS","url":"https://api.github.com/repos/shanyuhai123/learnCSS"},"public":true,"created_at":"2020-02-02T02:22:18Z"}]

func sliceDiff(s1 []*GithubEvent, s2 []*GithubEvent) []*GithubEvent {
	result := make([]*GithubEvent, 0)
	for _, item1 := range s1 {
		exist := false
		for _, item2 := range s2 {
			if reflect.DeepEqual(item1, item2) {
				exist = true
				break
			}
		}
		if !exist {
			result = append(result, item1)
		}
	}
	return result
}
