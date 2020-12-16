package bot

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/github-telebot/src/bot/controller"
	"github.com/Aoi-hosizora/github-telebot/src/bot/server"
	"github.com/Aoi-hosizora/github-telebot/src/config"
	"gopkg.in/tucnak/telebot.v2"
	"time"
)

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
	fmt.Println("[Telebot] Success to connect telegram bot:", b.Me.Username)
	fmt.Println()

	server.Bot = server.NewBotServer(b)
	initHandler(server.Bot)

	return nil
}

func initHandler(b *server.BotServer) {
	// start
	b.InlineButtons["btn_cancel"] = &telebot.InlineButton{Unique: "btn_cancel", Text: "Cancel"}
	b.HandleMessage("/start", controller.StartCtrl)
	b.HandleMessage("/help", controller.HelpCtrl)
	b.HandleMessage("/cancel", controller.CancelCtrl)
	b.HandleInline(b.InlineButtons["btn_cancel"], controller.InlBtnCancelCtrl)
	b.HandleMessage(telebot.OnText, controller.OnTextCtrl)

	// user
	b.InlineButtons["btn_unbind"] = &telebot.InlineButton{Unique: "btn_unbind", Text: "Unbind"}
	b.HandleMessage("/bind", controller.BindCtrl)
	b.HandleMessage("/unbind", controller.UnbindCtrl)
	b.HandleInline(b.InlineButtons["btn_unbind"], controller.InlBtnUnbindCtrl)
	b.HandleMessage("/me", controller.MeCtrl)
	b.HandleMessage("/enablesilent", controller.EnableSilentCtrl)
	b.HandleMessage("/disablesilent", controller.DisableSilentCtrl)

	// event
	b.InlineButtons["btn_filter"] = &telebot.InlineButton{Unique: "btn_filter", Text: "Filter"}
	b.InlineButtons["btn_not_filter"] = &telebot.InlineButton{Unique: "btn_not_filter", Text: "Not Filter"}
	b.HandleMessage("/allowissue", controller.AllowIssueCtrl)
	b.HandleMessage("/disallowissue", controller.DisallowIssueCtrl)
	b.HandleInline(b.InlineButtons["btn_filter"], controller.InlBtnFilterCtrl)
	b.HandleInline(b.InlineButtons["btn_not_filter"], controller.InlBtnNotFilterCtrl)
	b.HandleMessage("/activity", controller.ActivityCtrl)
	b.HandleMessage("/activityn", controller.ActivityNCtrl)
	b.HandleMessage("/issue", controller.IssueCtrl)
	b.HandleMessage("/issuen", controller.IssueNCtrl)
}

func SendToChat(chatId int64, what interface{}, options ...interface{}) error {
	chat, err := server.Bot.ChatByID(xnumber.I64toa(chatId))
	if err != nil {
		return err
	}

	return server.Bot.Send(chat, what, options...)
}
