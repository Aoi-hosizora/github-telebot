package bot

import (
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/controller"
	"gopkg.in/tucnak/telebot.v2"
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
/subscribe - subscribe with a new github account
/unsubscribe - unsubscribe the current github account
/me - show the subscribed user's information

*Option*
/allowissue - allow bot to send issue events
/disallowissue - disallow bot to send issue events
/enablesilent - send message with no notification
/disablesilent - send message with notification
/enablepreview - enable preview for links
/disablepreview - disable preview for links

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
	{"/subscribe", "subscribe with a new github account"},
	{"/unsubscribe", "unsubscribe the current github account"},
	{"/me", "show the subscribed user's information"},
	{"/allowissue", "allow bot to send issue events"},
	{"/disallowissue", "disallow bot to send issue events"},
	{"/enablesilent", "send message with no notification"},
	{"/disablesilent", "send message with notification"},
	{"/enablepreview", "enable preview for links"},
	{"/disablepreview", "disable preview for links"},
	{"/activity", "show the first page of activity events"},
	{"/issue", "show the first page of issue events"},
}
