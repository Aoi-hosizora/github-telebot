package bot

import (
	"fmt"
	"github.com/Aoi-hosizora/github-telebot/src/bot/fsm"
	"github.com/Aoi-hosizora/github-telebot/src/config"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"time"
)

var (
	Bot *bot
)

type bot struct {
	Bot           *telebot.Bot
	UserStates    map[int64]fsm.UserStatus
	InlineButtons map[string]*telebot.InlineButton
	ReplyButtons  map[string]*telebot.ReplyButton
}

func Setup() error {
	b, err := telebot.NewBot(telebot.Settings{
		Token:   config.Configs.Bot.Token,
		Verbose: false,
		Poller: &telebot.LongPoller{
			Timeout: time.Second * time.Duration(config.Configs.Bot.PollerTimeout),
		},
	})
	if err != nil {
		return err
	}
	fmt.Println()
	log.Println("[Telebot] Success to connect telegram bot:", b.Me.Username)

	Bot = &bot{
		Bot:           b,
		UserStates:    make(map[int64]fsm.UserStatus),
		InlineButtons: make(map[string]*telebot.InlineButton),
		ReplyButtons:  make(map[string]*telebot.ReplyButton),
	}
	Bot.initHandler()

	return nil
}

func (b *bot) initHandler() {
	b.InlineButtons["btn_unbind"] = &telebot.InlineButton{Unique: "btn_unbind", Text: "Unbind"}
	b.InlineButtons["btn_cancel"] = &telebot.InlineButton{Unique: "btn_cancel", Text: "Cancel"}

	b.handleMessage("/start", startCtrl)
	b.handleMessage("/help", helpCtrl)
	b.handleMessage("/bind", startBindCtrl)
	b.handleMessage("/me", meCtrl)
	b.handleMessage("/unbind", startUnbindCtrl)
	b.handleMessage("/cancel", cancelCtrl)
	b.handleMessage("/send", sendCtrl)
	b.handleMessage("/sendn", startSendnCtrl)
	b.handleMessage("/issue", sendIssueCtrl)
	b.handleMessage("/issuen", startSendIssuenCtrl)

	b.handleInline(b.InlineButtons["btn_unbind"], unbindBtnCtrl)
	b.handleInline(b.InlineButtons["btn_cancel"], cancelBtnCtrl)

	b.handleMessage(telebot.OnText, onTextCtrl)
}

func (b *bot) Start() {
	b.Bot.Start()
}

func (b *bot) Stop() {
	b.Bot.Stop()
}
