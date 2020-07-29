package bot

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xstatus"
	"github.com/Aoi-hosizora/github-telebot/src/bot/fsm"
	"github.com/Aoi-hosizora/github-telebot/src/model"
	"github.com/Aoi-hosizora/github-telebot/src/service"
	"gopkg.in/tucnak/telebot.v2"
	"strconv"
	"strings"
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
/activityn - show the nth page of activity events
/issue - show the first page of issue events
/issuen - show the nth page of issue events`

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
	BIND_SUCCESS_TIP   = "(Tips: new activity events will be sent periodically, but issue events will not be sent automatically, use /allowIssue to allow)"

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
func startCtrl(m *telebot.Message) {
	_ = Bot.Reply(m, START)
}

// /help
func helpCtrl(m *telebot.Message) {
	_ = Bot.Reply(m, HELP, telebot.ModeMarkdown)
}

// /cancel
func cancelCtrl(m *telebot.Message) {
	if Bot.UserStates[m.Chat.ID] == fsm.None {
		_ = Bot.Reply(m, NO_ACTION)
	} else {
		Bot.UserStates[m.Chat.ID] = fsm.None
		_ = Bot.Reply(m, ACTION_CANCELED, &telebot.ReplyMarkup{
			ReplyKeyboardRemove: true,
		})
	}
}

// /bind
func bindCtrl(m *telebot.Message) {
	user := model.GetUser(m.Chat.ID)
	if user != nil {
		_ = Bot.Reply(m, BIND_ALREADY)
	} else {
		Bot.UserStates[m.Chat.ID] = fsm.Binding
		_ = Bot.Reply(m, BIND_START)
	}
}

// /bind -> x
func fromBindingCtrl(m *telebot.Message) {
	text := strings.TrimSpace(m.Text)
	if text == "" {
		_ = Bot.Reply(m, BIND_EMPTY)
		return
	}

	sp := strings.Split(text, " ")
	username := sp[0]
	user := &model.User{ChatID: m.Chat.ID, Username: username}
	if len(sp) >= 2 {
		user.Token = sp[1]
	}

	flag := ""
	ok, err := service.CheckUser(user.Username, user.Token)
	if err != nil {
		flag = GITHUB_FAILED
	} else if !ok {
		flag = GITHUB_NOT_FOUND
	} else {
		status := model.AddUser(user)
		if status == xstatus.DbExisted {
			flag = BIND_ALREADY
		} else if status == xstatus.DbFailed {
			flag = BIND_FAILED
		} else if user.Token != "" {
			flag = fmt.Sprintf(BIND_TOKEN_SUCCESS, username)
		} else {
			flag = fmt.Sprintf(BIND_SUCCESS, username)
		}
	}

	Bot.UserStates[m.Chat.ID] = fsm.None
	_ = Bot.Reply(m, flag)
}

// /me
func meCtrl(m *telebot.Message) {
	user := model.GetUser(m.Chat.ID)
	flag := ""
	if user == nil {
		flag = BIND_NOT_YET
	} else {
		n := fmt.Sprintf("[%s](https://github.com/%s)", user.Username, user.Username)
		if user.Token != "" {
			flag = fmt.Sprintf(GITHUB_ME_TOKEN, n)
		} else {
			flag = fmt.Sprintf(GITHUB_ME, n)
		}
	}
	_ = Bot.Reply(m, flag)
}

// /unbind
func unbindCtrl(m *telebot.Message) {
	user := model.GetUser(m.Chat.ID)
	if user == nil {
		_ = Bot.Reply(m, BIND_NOT_YET)
		return
	}

	flag := fmt.Sprintf(UNBIND_START, user.Username)
	_ = Bot.Reply(m, flag, &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{*Bot.InlineButtons["btn_unbind"]}, {*Bot.InlineButtons["btn_cancel"]},
		},
	})
}

// inl:btn_cancel
func inlBtnCancelCtrl(c *telebot.Callback) {
	_ = Bot.Bot.Delete(c.Message)
	_ = Bot.Reply(c.Message, ACTION_CANCELED)
}

// inl:btn_unbind
func inlBtnUnbindCtrl(c *telebot.Callback) {
	_ = Bot.Bot.Delete(c.Message)
	flag := ""
	status := model.DeleteUser(c.Message.Chat.ID)
	if status == xstatus.DbNotFound {
		flag = BIND_NOT_YET
	} else if status == xstatus.DbFailed {
		flag = UNBIND_FAILED
	} else {
		flag = UNBIND_SUCCESS
	}

	_ = Bot.Reply(c.Message, flag)
}

// /allowIssue
func allowIssueCtrl(m *telebot.Message) {
	user := model.GetUser(m.Chat.ID)
	if user == nil {
		_ = Bot.Reply(m, BIND_NOT_YET)
		return
	} else if user.Token == "" {
		_ = Bot.Reply(m, ISSUE_ONLY_FOR_TOKEN)
		return
	}

	flag := ""
	user.AllowIssue = true
	status := model.UpdateUser(user)
	if status == xstatus.DbNotFound {
		flag = BIND_NOT_YET
	} else if status == xstatus.DbFailed {
		flag = ISSUE_ALLOW_FAILED
	} else {
		flag = ISSUE_ALLOW_SUCCESS
	}

	_ = Bot.Reply(m, flag)
}

// /disallowIssue
func disallowIssueCtrl(m *telebot.Message) {
	user := model.GetUser(m.Chat.ID)
	if user == nil {
		_ = Bot.Reply(m, BIND_NOT_YET)
		return
	} else if user.Token == "" {
		_ = Bot.Reply(m, ISSUE_ONLY_FOR_TOKEN)
		return
	}

	flag := ""
	user.AllowIssue = false
	status := model.UpdateUser(user)
	if status == xstatus.DbNotFound {
		flag = BIND_NOT_YET
	} else if status == xstatus.DbFailed {
		flag = ISSUE_DISALLOW_FAILED
	} else {
		flag = ISSUE_DISALLOW_SUCCESS
	}

	_ = Bot.Reply(m, flag)
}

// /activity
func activityCtrl(m *telebot.Message) {
	m.Text = "1"
	fromActivitynCtrl(m)
}

// /activityn
func activitynCtrl(m *telebot.Message) {
	Bot.UserStates[m.Chat.ID] = fsm.Activityn
	_ = Bot.Reply(m, GITHUB_SEND_PAGE)
}

// /activityn -> x
func fromActivitynCtrl(m *telebot.Message) {
	page, err := strconv.Atoi(m.Text)
	if err != nil {
		_ = Bot.Reply(m, NUM_REQUIRED)
		return
	} else if page <= 0 {
		page = 1
	}

	flag := ""
	user := model.GetUser(m.Chat.ID)
	if user == nil {
		flag = BIND_NOT_YET
	} else {
		if resp, err := service.GetActivityEvents(user.Username, user.Token, page); err != nil {
			flag = GITHUB_FAILED
		} else if events, err := model.UnmarshalActivityEvents(resp); err != nil {
			flag = GITHUB_FAILED
		} else if render := service.RenderActivities(events); render == "" {
			render = GITHUB_EMPTY
		} else {
			flag = fmt.Sprintf("%s\n---\nFrom [%s](https://github.com/%s) (page %d)", render, user.Username, user.Username, page)
		}
	}

	Bot.UserStates[m.Chat.ID] = fsm.None
	_ = Bot.Reply(m, flag, telebot.ModeMarkdown)
}

// /issue
func issueCtrl(m *telebot.Message) {
	m.Text = "1"
	fromIssuenCtrl(m)
}

// /issuen
func issuenCtrl(m *telebot.Message) {
	Bot.UserStates[m.Chat.ID] = fsm.Issuen
	_ = Bot.Reply(m, GITHUB_SEND_PAGE)
}

// /issuen -> x
func fromIssuenCtrl(m *telebot.Message) {
	page, err := strconv.Atoi(m.Text)
	if err != nil {
		_ = Bot.Reply(m, NUM_REQUIRED)
		return
	} else if page <= 0 {
		page = 1
	}

	flag := ""
	user := model.GetUser(m.Chat.ID)
	if user == nil {
		flag = BIND_NOT_YET
	} else {
		if resp, err := service.GetIssueEvents(user.Username, user.Token, page); err != nil {
			flag = GITHUB_FAILED
		} else if events, err := model.UnmarshalIssueEvents(resp); err != nil {
			flag = GITHUB_FAILED
		} else if render := service.RenderIssues(events); render == "" {
			render = GITHUB_EMPTY
		} else {
			flag = service.RenderResult(render, user.Username) + fmt.Sprintf(" (page %d)", page)
		}
	}

	Bot.UserStates[m.Chat.ID] = fsm.None
	_ = Bot.Reply(m, flag, telebot.ModeMarkdown)
}

// onText
func onTextCtrl(m *telebot.Message) {
	switch Bot.UserStates[m.Chat.ID] {
	case fsm.Binding:
		fromBindingCtrl(m)
	case fsm.Activityn:
		fromActivitynCtrl(m)
	case fsm.Issuen:
		fromIssuenCtrl(m)
	default:
		_ = Bot.Reply(m, "Unknown command: "+m.Text)
	}
}
