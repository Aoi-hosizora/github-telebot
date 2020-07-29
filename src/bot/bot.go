package bot

import (
	"fmt"
	"github.com/Aoi-hosizora/github-telebot/src/bot/controller"
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

	b.handleMessage("/start", controller.StartCtrl)
	b.handleMessage("/help", controller.HelpCtrl)
	b.handleMessage("/cancel", controller.CancelCtrl)
	b.handleMessage("/bind", controller.BindCtrl)
	b.handleMessage("/unbind", controller.UnbindCtrl)
	b.handleMessage("/me", controller.MeCtrl)

	b.handleMessage("/allowIssue", controller.AllowIssueCtrl)
	b.handleMessage("/disallowIssue", controller.DisallowIssueCtrl)
	b.handleMessage("/activity", controller.ActivityCtrl)
	b.handleMessage("/activityN", controller.ActivityNCtrl)
	b.handleMessage("/issue", controller.IssueCtrl)
	b.handleMessage("/issueN", controller.IssueNCtrl)

	b.handleInline(b.InlineButtons["btn_unbind"], controller.InlBtnUnbindCtrl)
	b.handleInline(b.InlineButtons["btn_cancel"], controller.InlBtnCancelCtrl)

	b.handleMessage(telebot.OnText, controller.OnTextCtrl)
}

func (b *bot) Start() {
	b.Bot.Start()
}

func (b *bot) Stop() {
	b.Bot.Stop()
}
