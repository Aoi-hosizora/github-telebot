package main

import (
	"flag"
	"github.com/Aoi-hosizora/ah-tgbot/bot"
	"github.com/Aoi-hosizora/ah-tgbot/config"
	"github.com/Aoi-hosizora/ah-tgbot/model"
	"github.com/Aoi-hosizora/ah-tgbot/task"
	"log"
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
	err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalln("Failed to load config:", err)
	}
	err = model.SetupGorm()
	if err != nil {
		log.Fatalln("Failed to database:", err)
	}
	err = bot.Load()
	if err != nil {
		log.Fatalln("Failed to connect telegram bot:", err)
	}
	defer bot.Stop()

	go task.Polling()
	bot.Start()
}
