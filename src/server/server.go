package server

import (
	"github.com/Aoi-hosizora/ah-tgbot/src/config"
	"github.com/Aoi-hosizora/ah-tgbot/src/util"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"time"
)

type BotServer struct {
	Config *config.Config
	Bot    *telebot.Bot
}

func NewBotServer(config *config.Config) *BotServer {
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
	server := &BotServer{
		Config: config,
		Bot:    bot,
	}
	log.Printf("Success to connect bot server \"%s\"\n", bot.Me.Username)

	setupRoute(server)
	return server
}

func (b *BotServer) Send(from *telebot.Message, what interface{}, options ...interface{}) {
	msg, err := b.Bot.Send(from.Sender, what, options...)
	util.BotLog(from, msg, err)
}

func (b *BotServer) SendTo(to telebot.Recipient, what interface{}, options ...interface{}) {
	msg, err := b.Bot.Send(to, what, options...)
	util.BotLog(nil, msg, err)
}

func (b *BotServer) Serve() {
	go b.Bot.Start()
	Polling(b)
}
