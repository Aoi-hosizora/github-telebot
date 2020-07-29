package controller

import (
	"github.com/Aoi-hosizora/github-telebot/src/bot"
	"github.com/Aoi-hosizora/github-telebot/src/bot/fsm"
	"gopkg.in/tucnak/telebot.v2"
)

// noinspection GoSnakeCaseUsage
const (
	START = "This is AoiHosizora's github telebot. Use /help to show help message."
	HELP  = `**Commands**
/start - show start message
/help - show help message
/cancel - cancel the last action
/bind - bind with a new github account
/unbind - unbind an old github account
/me - show the bound user's information
/allowIssue - allow bot to send issue periodically
/disallowIssue - allow bot to send issue periodically
/activity - show the first page of activity events
/activityN - show the nth page of activity events
/issue - show the first page of issue events
/issueN - show the nth page of issue events`

	NO_ACTION       = "There is no action now."
	ACTION_CANCELED = "Current action has been canceled."
	NUM_REQUIRED    = "Excepted integer, but got a string. Please resend an integer."

	BIND_START         = "Please send github's username, and token (split with whitespace) if you want to watch private events also. /cancel to cancel."
	BIND_ALREADY       = "You have already bound with a github account."
	BIND_NOT_YET       = "You have not bound a github account yet."
	BIND_EMPTY         = "Please resend a non-empty username again."
	BIND_FAILED        = "Failed to bind github account, please retry later."
	BIND_SUCCESS       = "Binding user %s without token success. /activity to get activity events, /issue to get issue events.\n" + BIND_SUCCESS_TIP
	BIND_TOKEN_SUCCESS = "Binding user %s with token success. /send to get events, /issue to get issue events.\n" + BIND_SUCCESS_TIP
	BIND_SUCCESS_TIP   = "(Tips: new activity events will be sent periodically, but issue events will not be sent. Use /allowIssue to allow)"

	UNBIND_START   = "Sure to unbind the current github account %s?"
	UNBIND_FAILED  = "Failed to unbind github account, please retry later."
	UNBIND_SUCCESS = "Unbind user success."

	ISSUE_ONLY_FOR_TOKEN   = "This only can be allowed if you have bound with a token."
	ISSUE_ALLOW_SUCCESS    = "Success to allow bot to send issue events periodically."
	ISSUE_DISALLOW_SUCCESS = "Success to disallow bot to send issue events periodically."
	ISSUE_ALLOW_FAILED     = "Failed to allow bot to send issue events periodically."
	ISSUE_DISALLOW_FAILED  = "Failed to disallow bot to send issue events periodically."

	GITHUB_ME        = "You have bound with user: %s without token."
	GITHUB_ME_TOKEN  = "You have bound with user: %s with token."
	GITHUB_FAILED    = "Failed to get github information, please retry later."
	GITHUB_NOT_FOUND = "Github user not found."
	GITHUB_EMPTY     = "Empty events: \\[]"
	GITHUB_SEND_PAGE = "Please send the page you want to get, number required."
)

// /start
func StartCtrl(m *telebot.Message) {
	_ = bot.Bot.Reply(m, START)
}

// /help
func HelpCtrl(m *telebot.Message) {
	_ = bot.Bot.Reply(m, HELP, telebot.ModeMarkdown)
}

// /cancel
func CancelCtrl(m *telebot.Message) {
	if bot.Bot.UserStates[m.Chat.ID] == fsm.None {
		_ = bot.Bot.Reply(m, NO_ACTION)
	} else {
		bot.Bot.UserStates[m.Chat.ID] = fsm.None
		_ = bot.Bot.Reply(m, ACTION_CANCELED, &telebot.ReplyMarkup{
			ReplyKeyboardRemove: true,
		})
	}
}

// onText
func OnTextCtrl(m *telebot.Message) {
	switch bot.Bot.UserStates[m.Chat.ID] {
	case fsm.Binding:
		fromBindingCtrl(m)
	case fsm.ActivityN:
		fromActivityNCtrl(m)
	case fsm.IssueN:
		fromIssueNCtrl(m)
	default:
		_ = bot.Bot.Reply(m, "Unknown command: "+m.Text)
	}
}
