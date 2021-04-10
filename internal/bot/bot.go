package bot

import (
	"github.com/Aoi-hosizora/ahlib-web/xrecovery"
	"github.com/Aoi-hosizora/ahlib/xruntime"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/button"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/controller"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/server"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/config"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/logger"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"time"
)

func Setup() error {
	b, err := telebot.NewBot(telebot.Settings{
		Token:   config.Configs().Bot.Token,
		Verbose: false,
		Poller: &telebot.LongPoller{
			Timeout: time.Second * time.Duration(config.Configs().Bot.PollerTimeout),
		},
		Reporter: func(err error) {
			xrecovery.LogToLogrus(logger.Logger(), err, xruntime.RuntimeTraceStack(4))
		},
	})
	if err != nil {
		return err
	}

	log.Println("Success to connect telegram bot:", b.Me.Username)
	server.SetupBot(server.NewBotServer(b))
	setupHandler(server.Bot())

	return nil
}

func setupHandler(b *server.BotServer) {
	// start
	b.HandleMessage("/start", controller.StartCtrl)
	b.HandleMessage("/help", controller.HelpCtrl)
	b.HandleMessage("/cancel", controller.CancelCtrl)
	b.HandleMessage(telebot.OnText, controller.OnTextCtrl)

	// user
	b.HandleMessage("/bind", controller.BindCtrl)
	b.HandleMessage("/unbind", controller.UnbindCtrl)
	b.HandleMessage("/me", controller.MeCtrl)
	b.HandleInline(button.InlineBtnUnbind, controller.InlineBtnUnbindCtrl)
	b.HandleInline(button.InlineBtnCancel, controller.InlineBtnCancelCtrl)
	b.HandleMessage("/enablesilent", controller.EnableSilentCtrl)
	b.HandleMessage("/disablesilent", controller.DisableSilentCtrl)

	// event
	b.HandleMessage("/allowissue", controller.AllowIssueCtrl)
	b.HandleInline(button.InlineBtnFilter, controller.InlineBtnFilterCtrl)
	b.HandleInline(button.InlineBtnNotFilter, controller.InlineBtnNotFilterCtrl)
	b.HandleMessage("/disallowissue", controller.DisallowIssueCtrl)
	b.HandleMessage("/activity", controller.ActivityCtrl)
	b.HandleMessage("/activitypage", controller.ActivityPageCtrl)
	b.HandleMessage("/issue", controller.IssueCtrl)
	b.HandleMessage("/issuepage", controller.IssuePageCtrl)
}
