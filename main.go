package main

import (
	"flag"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"time"
)

var (
	help       bool
	configPath string
)

func init() {
	flag.BoolVar(&help, "h", false, "show help")
	flag.StringVar(&configPath, "config", "./config.yaml", "change the config path")
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
	} else {
		run()
	}
}

func run() {
	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatalln("Failed to load config file:", err)
	}

	poller := &telebot.LongPoller{
		Timeout: 2 * time.Second,
	}
	mwPoller := telebot.NewMiddlewarePoller(poller, func(update *telebot.Update) bool {
		log.Printf("[receive]\t%d\t\"%s\"\n", update.Message.ID, update.Message.Text)
		return true
	})

	bot, err := telebot.NewBot(telebot.Settings{
		Token:   config.TelegramConfig.Token,
		Updates: 0,
		Poller:  mwPoller,
	})
	if err != nil {
		log.Fatalln("Failed to connect bot server:", err)
	}
	log.Println("Success to connect bot server:", bot.Me.Username)

	bot.Handle("/hello", func(m *telebot.Message) {
		msg, err := bot.Send(m.Sender, "Hello world")
		if err == nil {
			timeSpan := float64(msg.Time().Sub(m.Time()).Nanoseconds()) / 1e6
			log.Printf("[send]\t\t%d\t\"%s\"\t|%d\t%.0fms\n", m.ID, m.Text, msg.ID, timeSpan)
		}
	})

	defer func() {
		bot.Stop()
	}()

	bot.Start()
}
