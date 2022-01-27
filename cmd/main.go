package main

import (
	"flag"
	"fmt"
	"github.com/Aoi-hosizora/github-telebot/internal/bot"
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
	err = database.SetupGormDB()
	if err != nil {
		log.Fatalln("Failed to setup gorm db:", err)
	}
	err = database.SetupRedisClient()
	if err != nil {
		log.Fatalln("Failed to setup redis client:", err)
	}

	c, err := bot.NewConsumer()
	if err != nil {
		log.Fatalln("Failed to create consumer:", err)
	}
	t, err := task.NewTask(c.BotWrapper())
	if err != nil {
		log.Fatalln("Failed to create task:", err)
	}

	fmt.Println()
	t.Start()
	defer t.Finish()
	c.Start()
}
