package bot

import (
	"context"
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
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

func processError(what string, errMsg string, options []interface{}) (string, []interface{}) {
	what = strings.ReplaceAll(what, "\\", "")
	what += "\n\nPlease contact to @AoiHosizora with the message: \"" + errMsg + "\""
	newOptions := make([]interface{}, 1, len(options))
	newOptions[0] = telebot.ModeMarkdown
	for _, opt := range options {
		if opt != telebot.ModeMarkdownV2 {
			newOptions = append(newOptions, opt)
		}
	}
	return what, newOptions
}

var ctxRetryCountKey = "retry_count"

func processSendError(bw *xtelebot.BotWrapper, err error, ev *xtelebot.RespondEvent) {
	source, what, options, ctx := ev.SendSource, ev.SendWhat, ev.SendOptions, ev.RespondContext
	if c, ok := ctx.Value(ctxRetryCountKey).(int); !ok {
		ctx = context.WithValue(ctx, ctxRetryCountKey, 1)
	} else if c <= 3 { // max retry count: 3
		ctx = context.WithValue(ctx, ctxRetryCountKey, c+1)
	} else {
		return
	}

	strWhat, ok := what.(string)
	if ok && xtelebot.IsEscapedParseError(err) {
		newWhat, newOptions := processError(strWhat, err.Error(), options)
		_, _ = bw.RespondSendCtx(ctx, source, newWhat, newOptions...)
	}
}

func processReplyError(bw *xtelebot.BotWrapper, err error, ev *xtelebot.RespondEvent) {
	source, explicit, what, options, ctx := ev.ReplySource, ev.ReplyExplicit, ev.ReplyWhat, ev.ReplyOptions, ev.RespondContext
	if c, ok := ctx.Value(ctxRetryCountKey).(int); !ok {
		ctx = context.WithValue(ctx, ctxRetryCountKey, 1)
	} else if c <= 3 { // max retry count: 3
		ctx = context.WithValue(ctx, ctxRetryCountKey, c+1)
	} else {
		return
	}

	strWhat, ok := what.(string)
	if ok && xtelebot.IsEscapedParseError(err) {
		newWhat, newOptions := processError(strWhat, err.Error(), options)
		_, _ = bw.RespondReplyCtx(ctx, source, explicit, newWhat, newOptions...)
	}
}
