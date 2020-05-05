package bot

import (
	"github.com/Aoi-hosizora/ah-tgbot/bot/fsm"
	"github.com/Aoi-hosizora/ah-tgbot/config"
	"github.com/Aoi-hosizora/ah-tgbot/logger"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"time"
)

// ロミオとシンデレラ
var (
	Bot        *telebot.Bot
	UserStates map[int64]fsm.UserStatus
	InlineBtns map[string]*telebot.InlineButton
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
	InlineBtns = make(map[string]*telebot.InlineButton)
	makeHandle()
	return nil
}

func Start() {
	Bot.Start()
}

func Stop() {
	Bot.Stop()
}

func handleWithLogger(endpoint interface{}, handler interface{}) {
	if h, ok := handler.(func(*telebot.Message)); ok {
		Bot.Handle(endpoint, func(m *telebot.Message) {
			logger.RcvLogger(m, endpoint)
			h(m)
		})
	} else {
		Bot.Handle(endpoint, handler)
	}
}

func makeHandle() {
	InlineBtns["btn_unbind"] = &telebot.InlineButton{Unique: "btn_unbind", Text: "Unbind"}
	InlineBtns["btn_cancel"] = &telebot.InlineButton{Unique: "btn_cancel", Text: "Cancel"}

	handleWithLogger("/start", startCtrl)
	handleWithLogger("/bind", startBindCtrl)
	handleWithLogger("/unbind", startUnbindCtrl)
	handleWithLogger("/cancel", cancelCtrl)

	handleWithLogger(InlineBtns["btn_unbind"], unbindBtnCtrl)
	handleWithLogger(InlineBtns["btn_cancel"], cancelBtnCtrl)

	handleWithLogger(telebot.OnText, onTextCtrl)
}
