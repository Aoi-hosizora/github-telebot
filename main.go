package main

import (
	"flag"
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
	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatalln("Failed to load config file:", err)
	}
	bot := newBot(config)
	defer func() {
		bot.Stop()
	}()

	go bot.Start()
	polling(config, bot)
}
