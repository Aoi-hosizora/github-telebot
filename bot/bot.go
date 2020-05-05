package bot

import (
	"github.com/Aoi-hosizora/ah-tgbot/bot/fsm"
	"github.com/Aoi-hosizora/ah-tgbot/config"
	"github.com/Aoi-hosizora/ah-tgbot/logger"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"time"
)

var (
	Bot        *telebot.Bot
	UserStates map[int64]fsm.UserStatus
)

func Load() error {
	cfg := config.Configs.TelegramConfig
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  cfg.BotToken,
		Poller: &telebot.LongPoller{Timeout: time.Second * time.Duration(cfg.PollerTimeout)},
	})
	if err != nil {
		return err
	}
	log.Println("[telebot] Success to connect telegram bot:", bot.Me.Username)

	Bot = bot
	UserStates = make(map[int64]fsm.UserStatus)
	makeHandle()
	return nil
}

func Start() {
	Bot.Start()
}

func Stop() {
	Bot.Stop()
}

func handleWithLogger(endpoint string, handler func(*telebot.Message)) {
	Bot.Handle(endpoint, func(m *telebot.Message) {
		logger.RcvLogger(m, endpoint)
		handler(m)
	})
}

func makeHandle() {
	handleWithLogger("/start", startCtrl)
	handleWithLogger("/bind", startBindCtrl)
	handleWithLogger(telebot.OnText, onTextCtrl)
}
