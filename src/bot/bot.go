package bot

import (
	"fmt"
	"github.com/Aoi-hosizora/github-telebot/src/bot/fsm"
	"github.com/Aoi-hosizora/github-telebot/src/bot/server"
	"github.com/Aoi-hosizora/github-telebot/src/config"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"strconv"
	"time"
)

var Bot *server.BotServer

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

	log.Println("[Telebot] Success to connect telegram bot:", b.Me.Username)
	fmt.Println()

	Bot = &server.BotServer{
		Bot:           b,
		UserStates:    make(map[int64]fsm.UserStatus),
		InlineButtons: make(map[string]*telebot.InlineButton),
		ReplyButtons:  make(map[string]*telebot.ReplyButton),
	}
	initHandler(Bot)

	return nil
}

func initHandler(b *server.BotServer) {
	b.InlineButtons["btn_unbind"] = &telebot.InlineButton{Unique: "btn_unbind", Text: "Unbind"}
	b.InlineButtons["btn_cancel"] = &telebot.InlineButton{Unique: "btn_cancel", Text: "Cancel"}

	b.HandleMessage("/start", StartCtrl)
	b.HandleMessage("/help", HelpCtrl)
	b.HandleMessage("/cancel", CancelCtrl)
	b.HandleMessage("/bind", BindCtrl)
	b.HandleMessage("/unbind", UnbindCtrl)
	b.HandleMessage("/me", MeCtrl)

	b.HandleMessage("/allowissue", AllowIssueCtrl)
	b.HandleMessage("/disallowissue", DisallowIssueCtrl)
	b.HandleMessage("/activity", ActivityCtrl)
	b.HandleMessage("/activityn", ActivityNCtrl)
	b.HandleMessage("/issue", IssueCtrl)
	b.HandleMessage("/issuen", IssueNCtrl)

	b.HandleInline(b.InlineButtons["btn_unbind"], InlBtnUnbindCtrl)
	b.HandleInline(b.InlineButtons["btn_cancel"], InlBtnCancelCtrl)

	b.HandleMessage(telebot.OnText, OnTextCtrl)
}

func SendToChat(chatId int64, what interface{}, options ...interface{}) error {
	chat, err := Bot.Bot.ChatByID(strconv.FormatInt(chatId, 10))
	if err != nil {
		return err
	} else {
		return Bot.Send(chat, what, options...)
	}
}
