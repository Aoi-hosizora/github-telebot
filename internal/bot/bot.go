package bot

import (
	"github.com/Aoi-hosizora/ahlib-web/xgin"
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
	"github.com/Aoi-hosizora/ahlib/xruntime"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/config"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/logger"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Consumer struct {
	bw *xtelebot.BotWrapper
}

func (s *Consumer) BotWrapper() *xtelebot.BotWrapper {
	return s.bw
}

func NewConsumer() (*Consumer, error) {
	// telebot
	cfg := config.Configs().Meta
	bot, err := telebot.NewBot(telebot.Settings{
		Token:    cfg.Token,
		Reporter: func(err error) {}, // ignore
		Verbose:  false,
		Poller:   xtelebot.LongPoller(int(cfg.PollerTimeout)),
	})
	if err != nil {
		return nil, err
	}

	// wrapper
	bw := xtelebot.NewBotWrapper(bot)
	setupLoggers(bw)

	// handlers
	bw.SetHandledCallback(func(endpoint interface{}, renderedEndpoint string, handlerName string) {
		if config.IsDebugMode() {
			xtelebot.DefaultColorizedHandledCallback(endpoint, renderedEndpoint, handlerName)
		}
	})
	setupHandlers(bw)
	bw.SetHandledCallback(nil)

	s := &Consumer{bw: bw}
	return s, nil
}

func setupLoggers(bw *xtelebot.BotWrapper) {
	l := logger.Logger()
	bw.SetReceivedCallback(func(endpoint interface{}, received *telebot.Message) {
		xtelebot.LogReceiveToLogrus(l, endpoint, received)
	})
	bw.SetRespondedCallback(func(typ xtelebot.RespondEventType, event *xtelebot.RespondEvent) {
		xtelebot.LogRespondToLogrus(l, typ, event)
		if event.ReturnedError != nil {
			if typ == xtelebot.RespondSendEvent {
				processSendError(bw, event.ReturnedError, event)
			} else if typ == xtelebot.RespondReplyEvent {
				processReplyError(bw, event.ReturnedError, event)
			}
		}
	})
	bw.SetPanicHandler(func(_, _, v interface{}) {
		xgin.LogRecoveryToLogrus(l, v, xruntime.RuntimeTraceStack(3))
	})
}

func (s *Consumer) Start() {
	terminated := make(chan interface{})
	go func() {
		defer close(terminated)
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		sig := <-ch
		signal.Stop(ch)
		log.Printf("[Bot] Stopping due to %s received...", xruntime.SignalName(sig.(syscall.Signal)))
		s.bw.Bot().Stop()
	}()

	log.Printf("[Bot] Starting consuming incoming update on bot %s", s.bw.Bot().Me.Username)
	s.bw.Bot().Start()
	<-terminated
	log.Println("[Bot] Bot consumer is stopped successfully")
}
