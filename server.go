package main

import (
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"time"
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
	for {
		msg, err := bot.Send(chat, "a")
		botLog(nil, msg, err)
		time.Sleep(time.Second * time.Duration(config.ServerConfig.PollingDuration))
	}
}
