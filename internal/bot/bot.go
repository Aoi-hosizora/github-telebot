package bot

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xgin"
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
	"github.com/Aoi-hosizora/ahlib/xcolor"
	"github.com/Aoi-hosizora/ahlib/xruntime"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/controller"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/fsm"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/config"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/logger"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"time"
)

type Consumer struct {
	bw *xtelebot.BotWrapper
}

func (s *Consumer) BotWrapper() *xtelebot.BotWrapper {
	return s.bw
}

func NewConsumer() (*Consumer, error) {
	// telebot
	bot, err := telebot.NewBot(telebot.Settings{
		Token:   config.Configs().Meta.Token,
		Verbose: false,
		Poller:  &telebot.LongPoller{Timeout: time.Second * time.Duration(config.Configs().Meta.PollerTimeout)},
	})
	if err != nil {
		return nil, err
	}

	// wrapper
	bw := xtelebot.NewBotWrapper(bot)
	bw.Data().SetInitialState(fsm.None)
	bw.SetHandledCallback(func(_ interface{}, renderedEndpoint string, handlerName string) {
		if config.IsDebugMode() {
			fmt.Printf("[Bot-debug] %s --> %s\n", xcolor.Blue.Sprint(fmt.Sprintf("%-32s", renderedEndpoint)), handlerName)
		}
	})
	setupLoggers(bw)
	setupHandlers(bw)
	bw.SetHandledCallback(func(interface{}, string, string) {})

	s := &Consumer{bw: bw}
	return s, nil
}

func setupLoggers(bw *xtelebot.BotWrapper) {
	l := logger.Logger()
	bw.SetReceivedCallback(func(endpoint interface{}, received *telebot.Message) {
		xtelebot.LogReceiveToLogrus(l, endpoint, received)
	})
	bw.SetRepliedCallback(func(received *telebot.Message, replied *telebot.Message, err error) {
		xtelebot.LogReplyToLogrus(l, received, replied, err)
	})
	bw.SetSentCallback(func(chat *telebot.Chat, sent *telebot.Message, err error) {
		xtelebot.LogSendToLogrus(l, chat, sent, err)
	})
	bw.SetPanicHandler(func(endpoint interface{}, v interface{}) {
		xgin.LogRecoveryToLogrus(l, v, xruntime.RuntimeTraceStack(4))
	})
}

func (s *Consumer) Start() {
	log.Printf("[Bot] Starting consuming incoming update on bot %s", s.bw.Bot().Me.Username)
	s.bw.Bot().Start() // block to poll and consume
}

func setupHandlers(bw *xtelebot.BotWrapper) {
	// start
	bw.HandleCommand("/start", controller.Start)
	bw.HandleCommand("/help", controller.Help)
	bw.HandleCommand("/cancel", controller.Cancel)
	bw.HandleCommand(telebot.OnText, controller.OnText)

	// subscribe
	bw.HandleCommand("/subscribe", controller.Subscribe)
	bw.HandleCommand("/unsubscribe", controller.Unsubscribe)
	bw.HandleCommand("/me", controller.Me)

	// option
	bw.HandleCommand("/allowissue", controller.AllowIssue)
	bw.HandleCommand("/disallowissue", controller.DisallowIssue)
	bw.HandleCommand("/enablesilent", controller.EnableSilent)
	bw.HandleCommand("/disablesilent", controller.DisableSilent)
	bw.HandleCommand("/enablepreview", controller.EnablePreview)
	bw.HandleCommand("/disablepreview", controller.DisablePreview)

	// event
	bw.HandleCommand("/activity", controller.Activity)
	bw.HandleCommand("/issue", controller.Issue)
}
