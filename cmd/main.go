package main

import (
	"flag"
	"github.com/Aoi-hosizora/github-telebot/internal/bot"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/server"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/config"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/database"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/logger"
	"github.com/Aoi-hosizora/github-telebot/internal/task"
	"log"
)

var (
	fConfig = flag.String("config", "./config.yaml", "config file path")
	fHelp   = flag.Bool("h", false, "show help")
)

func main() {
	flag.Parse()
	if *fHelp {
		flag.Usage()
		return
	}

	err := config.Load(*fConfig)
	if err != nil {
		log.Fatalln("Failed to load config:", err)
	}
	err = logger.Setup()
	if err != nil {
		log.Fatalln("Failed to setup logger:", err)
	}
	err = database.SetupGorm()
	if err != nil {
		log.Fatalln("Failed to setup gorm:", err)
	}
	err = database.SetupRedis()
	if err != nil {
		log.Fatalln("Failed to setup redis:", err)
	}

	err = bot.Setup()
	if err != nil {
		log.Fatalln("Failed to setup telebot:", err)
	}
	err = task.Setup()
	if err != nil {
		log.Fatalln("Failed to setup cron:", err)
	}

	defer task.Cron().Stop()
	task.Cron().Start()

	defer server.Bot().Stop()
	server.Bot().Start()
}
