package bot

import (
	"fmt"
	"github.com/Aoi-hosizora/github-telebot/src/bot/button"
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
	b.HandleMessage("/start", controller.StartCtrl)
	b.HandleMessage("/help", controller.HelpCtrl)
	b.HandleMessage("/cancel", controller.CancelCtrl)
	b.HandleInline(button.InlineBtnCancel, controller.InlBtnCancelCtrl)
	b.HandleMessage(telebot.OnText, controller.OnTextCtrl)

	// user
	b.HandleMessage("/bind", controller.BindCtrl)
	b.HandleMessage("/unbind", controller.UnbindCtrl)
	b.HandleInline(button.InlineBtnUnbind, controller.InlBtnUnbindCtrl)
	b.HandleMessage("/me", controller.MeCtrl)
	b.HandleMessage("/enablesilent", controller.EnableSilentCtrl)
	b.HandleMessage("/disablesilent", controller.DisableSilentCtrl)

	// event
	b.HandleMessage("/allowissue", controller.AllowIssueCtrl)
	b.HandleInline(button.InlineBtnFilter, controller.InlBtnFilterCtrl)
	b.HandleInline(button.InlineBtnNotFilter, controller.InlBtnNotFilterCtrl)
	b.HandleMessage("/disallowissue", controller.DisallowIssueCtrl)
	b.HandleMessage("/activity", controller.ActivityCtrl)
	b.HandleMessage("/activityn", controller.ActivityNCtrl)
	b.HandleMessage("/issue", controller.IssueCtrl)
	b.HandleMessage("/issuen", controller.IssueNCtrl)
}
