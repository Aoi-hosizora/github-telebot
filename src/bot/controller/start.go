package controller

import (
	"fmt"
	"github.com/Aoi-hosizora/github-telebot/src/bot/fsm"
	"github.com/Aoi-hosizora/github-telebot/src/bot/server"
	"gopkg.in/tucnak/telebot.v2"
)

// noinspection GoSnakeCaseUsage
const (
	START = "Here is AoiHosizora's github telebot. Use /help to show help message."
	HELP  = `*Commands*
/start - show start message
/help - show help message
/cancel - cancel the last action

*Account*
/bind - bind with a new github account
/unbind - unbind an old github account
/me - show the bound user's information
/enablesilent - enable client send
/disablesilent - disable client send

*Events*
/allowissue - allow bot to send issue
/disallowissue - allow bot to send issue
/activity - show the first page of activity events
/activityn - show the nth page of activity events
/issue - show the first page of issue events
/issuen - show the nth page of issue events

*Bug report*
https://github.com/Aoi-hosizora/github-telebot/issues/new
`

	NO_ACTION       = "There is no action now."
	ACTION_CANCELED = "Current action has been canceled."
	UNKNOWN_COMMAND = "Unknown command: %s. Send /help to see help."
	NUM_REQUIRED    = "Excepted integer, but got a string. Please resend an integer."

	BIND_Q             = "Please send github's username, and token (split by whitespace) if you want to watch private events or issue events also. /cancel to cancel."
	BIND_ALREADY       = "You have already bound with a github account."
	BIND_NOT_YET       = "You have not bound a github account yet."
	BIND_EMPTY         = "Please resend a non-empty username again."
	BIND_FAILED        = "Failed to bind github account, please retry later."
	BIND_SUCCESS       = "Binding user %s without token success. /activity to get activity events, /issue to get issue events.\n" + BIND_SUCCESS_TIP
	BIND_TOKEN_SUCCESS = "Binding user %s with token success. /activity to get events, /issue to get issue events.\n" + BIND_SUCCESS_TIP
	BIND_SUCCESS_TIP   = "(Tips: new activity events will be sent periodically, but issue events will not be sent. Use /allowissue to allow)"

	SILENT_Q               = "Please send 2 different numbers as hour (in [0, 23]) you want to start and finish silent send, and with a timezone (such as +8:00 or -06:30), split by whitespace. Examples: 23 6 +8 or 0 8 -6."
	SILENT_FORMAT_REQUIRED = "Excepted input, please send 2 different numbers as hour (in [0, 23]) you want to start and finish silent send, and with a timezone."
	SILENT_HOUR_REQUIRED   = "Excepted hour, please send an integer in [0, 23]."
	SILENT_ZONE_REQUIRED   = "Excepted timezone, please send a right time zone, such as +8:00 or -06:30"
	SILENT_NOT_YET         = "You have not set silent yet, use /enablesilent to set."
	SILENT_SUCCESS         = "Success. Now it will be silent when send message in %s."
	SILENT_FAILED          = "Failed to set silent send, please retry later."
	DISABLE_SILENT_SUCCESS = "Disable silent success. Any message will be sent directly now."
	DISABLE_SILENT_FAILED  = "Failed to disable silent, please retry later."

	UNBIND_Q       = "Sure to unbind the current github account %s?"
	UNBIND_FAILED  = "Failed to unbind github account, please retry later."
	UNBIND_SUCCESS = "Unbind user success."

	ISSUE_ONLY_FOR_TOKEN   = "This only can be allowed if you have bound with a token."
	ISSUE_ALLOW_SUCCESS    = "Success to allow bot to send issue events periodically."
	ISSUE_DISALLOW_SUCCESS = "Success to disallow bot to send issue events periodically."
	ISSUE_ALLOW_FAILED     = "Failed to allow bot to send issue events periodically."
	ISSUE_DISALLOW_FAILED  = "Failed to disallow bot to send issue events periodically."

	GITHUB_ME          = "You have bound with user: %s without token."
	GITHUB_ME_TOKEN    = "You have bound with user: %s with token."
	GITHUB_FAILED      = "Failed to get github information, please retry later."
	GITHUB_NOT_FOUND   = "Github user not found."
	GITHUB_EMPTY       = "Empty events: \\[]"
	GITHUB_SEND_PAGE_Q = "Please send the page you want to get, number required."
)

// /start
func StartCtrl(m *telebot.Message) {
	_ = server.Bot.Reply(m, START)
}

// /help
func HelpCtrl(m *telebot.Message) {
	_ = server.Bot.Reply(m, HELP, telebot.ModeMarkdown)
}

// /cancel
func CancelCtrl(m *telebot.Message) {
	if server.Bot.UsersData.GetStatus(m.Chat.ID) == fsm.None {
		_ = server.Bot.Reply(m, NO_ACTION)
	} else {
		server.Bot.UsersData.SetStatus(m.Chat.ID, fsm.None)
		_ = server.Bot.Reply(m, ACTION_CANCELED, &telebot.ReplyMarkup{
			ReplyKeyboardRemove: true,
		})
	}
}

// $onText
func OnTextCtrl(m *telebot.Message) {
	switch server.Bot.UsersData.GetStatus(m.Chat.ID) {
	case fsm.Binding:
		fromBindingCtrl(m)
	case fsm.ActivityPage:
		fromActivityNCtrl(m)
	case fsm.IssuePage:
		fromIssueNCtrl(m)
	case fsm.SilentHour:
		fromSilentHourCtrl(m)
	default:
		_ = server.Bot.Reply(m, fmt.Sprintf(UNKNOWN_COMMAND, m.Text))
	}
}