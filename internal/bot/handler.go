package bot

import (
	"context"
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
	"github.com/Aoi-hosizora/ahlib/xgeneric/xsugar"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/controller"
	"gopkg.in/tucnak/telebot.v2"
	"strings"
)

func setupHandlers(bw *xtelebot.BotWrapper) {
	// start
	bw.HandleCommand("/start", controller.Start)
	bw.HandleCommand("/help", controller.Help(help))
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

	// set commands
	_ = bw.Bot().SetCommands(commands)
}

const help = `*Start*
/start - show start message
/help - show this help message
/cancel - cancel the last action

*Subscribe*
/subscribe - subscribe with a new GitHub account
/unsubscribe - unsubscribe the current GitHub account
/me - show the subscribed user's information

*Option*
/allowissue - allow bot to notify new issue events
/disallowissue - disallow bot to notify new issue events
/enablesilent - send message with no notification
/disablesilent - send message with notification
/enablepreview - enable preview for link
/disablepreview - disable preview for link

*Event*
/activity - show the first page of activity events
/activity N - show the N-th page of activity events
/issue - show the first page of issue events
/issue N - show the N-th page of issue events

*Bug report*
https://github.com/Aoi-hosizora/github-telebot/issues`

var commands = []telebot.Command{
	{"/start", "show start message"},
	{"/help", "show this help message"},
	{"/cancel", "cancel the last action"},
	{"/subscribe", "subscribe with a new GitHub account"},
	{"/unsubscribe", "unsubscribe the current GitHub account"},
	{"/me", "show the subscribed user's information"},
	{"/allowissue", "allow bot to notify new issue events"},
	{"/disallowissue", "disallow bot to notify new issue events"},
	{"/enablesilent", "send message with no notification"},
	{"/disablesilent", "send message with notification"},
	{"/enablepreview", "enable preview for link"},
	{"/disablepreview", "disable preview for link"},
	{"/activity", "show the first page of activity events"},
	{"/issue", "show the first page of issue events"},
}

// =============
// process error
// =============

const ctxRetryCountKey = "retry_count"

func processError(bw *xtelebot.BotWrapper, typ xtelebot.RespondEventType, ev *xtelebot.RespondEvent) {
	if typ != xtelebot.RespondSendEvent && typ != xtelebot.RespondReplyEvent {
		return // ignore error
	}
	err := ev.ReturnedError
	what, ok := xsugar.IfThenElse(typ == xtelebot.RespondSendEvent, ev.SendWhat, ev.ReplyWhat).(string)
	options := xsugar.IfThenElse(typ == xtelebot.RespondSendEvent, ev.SendOptions, ev.ReplyOptions)
	if !ok || len(options) == 0 || !xtelebot.IsEntityParseError(err) {
		return // ignore error
	}

	tryCount, ok := ev.RespondContext.Value(ctxRetryCountKey).(int)
	if !ok {
		tryCount = 1 // first try (=1)
	} else if tryCount > 3 {
		return // max retry count (>3)
	} else {
		tryCount++ // after the first try (++)
	}
	ctx := context.WithValue(ev.RespondContext, ctxRetryCountKey, tryCount)

	newWhat, newOptions := generateRetriedMsg(tryCount, err, what, options)
	if typ == xtelebot.RespondSendEvent {
		_, _ = bw.RespondSendCtx(ctx, ev.SendSource, newWhat, newOptions...)
	} else {
		_, _ = bw.RespondReplyCtx(ctx, ev.ReplySource, ev.ReplyExplicit, newWhat, newOptions...)
	}
}

func generateRetriedMsg(tryCount int, err error, what string, options []interface{}) (newWhat string, newOptions []interface{}) {
	firstTry := tryCount == 1
	if firstTry {
		newWhat = strings.ReplaceAll(what, "\\", "")
		newWhat += "\n\nPlease contact to @AoiHosizora with following message: \"" + err.Error() + "\""
	} else {
		newWhat = what
	}

	newOptions = make([]interface{}, 0, len(options))
	if firstTry {
		newOptions = append(newOptions, telebot.ModeMarkdown)
	}
	for _, opt := range options {
		if (opt != telebot.ModeMarkdownV2) && (firstTry || opt != telebot.ModeMarkdown) {
			newOptions = append(newOptions, opt)
		}
	}
	return what, newOptions
}
