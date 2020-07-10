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
	Bot           *telebot.Bot
	UserStates    map[int64]fsm.UserStatus
	InlineButtons map[string]*telebot.InlineButton
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
	InlineButtons = make(map[string]*telebot.InlineButton)
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
	if msg, ok := handler.(func(*telebot.Message)); ok {
		Bot.Handle(endpoint, func(m *telebot.Message) {
			logger.RcvLogger(m, endpoint)
			msg(m)
		})
	} else if cb, ok := handler.(func(*telebot.Callback)); ok {
		Bot.Handle(endpoint, func(c *telebot.Callback) {
			logger.RcvLogger(c, endpoint)
			cb(c)
		})
	} else {
		Bot.Handle(endpoint, handler)
	}
}

func makeHandle() {
	InlineButtons["btn_unbind"] = &telebot.InlineButton{Unique: "btn_unbind", Text: "Unbind"}
	InlineButtons["btn_cancel"] = &telebot.InlineButton{Unique: "btn_cancel", Text: "Cancel"}

	handleWithLogger("/start", startCtrl)
	handleWithLogger("/help", helpCtrl)
	handleWithLogger("/bind", startBindCtrl)
	handleWithLogger("/me", meCtrl)
	handleWithLogger("/unbind", startUnbindCtrl)
	handleWithLogger("/cancel", cancelCtrl)
	handleWithLogger("/send", sendCtrl)
	handleWithLogger("/sendn", startSendnCtrl)

	handleWithLogger(InlineButtons["btn_unbind"], unbindBtnCtrl)
	handleWithLogger(InlineButtons["btn_cancel"], cancelBtnCtrl)

	handleWithLogger(telebot.OnText, onTextCtrl)
}
