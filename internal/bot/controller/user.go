package controller

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/Aoi-hosizora/ahlib/xstring"
	"github.com/Aoi-hosizora/ahlib/xtime"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/button"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/fsm"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/server"
	"github.com/Aoi-hosizora/github-telebot/internal/service"
	"github.com/Aoi-hosizora/github-telebot/internal/service/dao"
	"gopkg.in/tucnak/telebot.v2"
	"strings"
)

const (
	BIND_ALREADY        = "You have already bind with a github account: %s."
	BIND_USERNAME_Q     = "Please send the github's username which you want to bind. Send /cancel to cancel."
	BIND_TOKEN_Q        = "Do you want to watch private events? Send your token if you want, otherwise send 'no'. Send /cancel to cancel."
	BIND_EMPTY_USERNAME = "Please send a non-empty username (without whitespace). Send /cancel to cancel."
	BIND_EMPTY_TOKEN    = "Please send a non-empty token (without whitespace), or send 'no' to ignore. Send /cancel to cancel."
	BIND_FAILED         = "Failed to bind github account, please retry later."
	BIND_TOKEN_SUCCESS  = "Binding user '%s' with token success. " + BIND_SUCCESS_TIP
	BIND_NOTOK_SUCCESS  = "Binding user '%s' without token success. " + BIND_SUCCESS_TIP
	BIND_SUCCESS_TIP    = "Send /activity to get activity events, send /issue to get issue events.\n" +
		"(Tips: new activity events will be sent periodically, but issue events will not be sent in default, you can use /allowissue to allow sending periodically)"

	BIND_NOT_YET   = "You have not bind with a github account yet."
	UNBIND_Q       = "Sure to unbind the current github account '%s'?"
	UNBIND_FAILED  = "Failed to unbind github account, please retry later."
	UNBIND_SUCCESS = "Unbind user success."

	GITHUB_FAILED          = "Failed to query information from github, please retry later."
	GITHUB_USER_NOT_FOUND  = "Github user is not found, or the token is invalid."
	GITHUB_ME_NOTOK        = "You have bound with user '%s' without token."
	GITHUB_ME_TOKEN        = "You have bound with user '%s' with token '%s'."
	GITHUB_ACTOR_NOT_FOUND = "Github actor is not found."
	GITHUB_REPO_NOT_FOUND  = "Github repo is not found."

	SILENT_Q = "Please send two different numbers to represent hours that you want to start and finish silent send (numbers are in range of [0, 23]), " +
		"and with a timezone (such as +8:00 or -06:30), split by whitespace. Examples: '23 6 +8:00' or '0 8 -6'. Send /cancel to cancel."
	SILENT_UNEXPECTED_FORMAT   = "Unexpected input, please send two different numbers in range of [0, 23] and with a timezone string. Send /cancel to cancel."
	SILENT_UNEXPECTED_HOUR     = "Unexpected hour value, please send two different integers in range of [0, 23]. Send /cancel to cancel."
	SILENT_UNEXPECTED_TIMEZONE = "Unexpected timezone string, please send a valid one, such as +8:00 or -06:30. Send /cancel to cancel."
	SILENT_FAILED              = "Failed to enable silence send, please retry later."
	SILENT_SUCCESS             = "Done. Now silence send will be enabled when %s."
	SILENT_NOT_YET             = "You have already disable silence send, use /enablesilent to enable."
	DISABLE_SILENT_FAILED      = "Failed to disable silent, please retry later."
	DISABLE_SILENT_SUCCESS     = "Disable silent success. Any message will be sent directly now."
)

// /bind
func BindCtrl(m *telebot.Message) {
	user := dao.QueryChat(m.Chat.ID)
	if user != nil {
		_ = server.Bot().Reply(m, fmt.Sprintf(BIND_ALREADY, user.Username))
	} else {
		server.Bot().SetStatus(m.Chat.ID, fsm.BindingUsername)
		_ = server.Bot().Reply(m, BIND_USERNAME_Q)
	}
}

// fsm.BindingUsername
func FromBindingUsernameCtrl(m *telebot.Message) {
	username := strings.TrimSpace(m.Text)
	if username == "" || strings.Contains(username, " ") {
		_ = server.Bot().Reply(m, BIND_EMPTY_USERNAME)
		return
	}

	server.Bot().SetCache(m.Chat.ID, "username", username)
	server.Bot().SetStatus(m.Chat.ID, fsm.BindingToken)
	_ = server.Bot().Reply(m, BIND_TOKEN_Q)
}

// fsm.BindingToken
func FromBindingTokenCtrl(m *telebot.Message) {
	token := strings.TrimSpace(m.Text)
	if token == "" || strings.Contains(token, " ") {
		_ = server.Bot().Reply(m, BIND_EMPTY_TOKEN)
		return
	}

	username, _ := server.Bot().GetCache(m.Chat.ID, "username")
	if strings.ToLower(token) == "no" {
		token = ""
	}

	flag := ""
	ok, err := service.CheckUserExistence(username.(string), token)
	if err != nil {
		flag = GITHUB_FAILED
	} else if !ok {
		flag = GITHUB_USER_NOT_FOUND
	} else {
		status := dao.CreateChat(m.Chat.ID, username.(string), token)
		if status == xstatus.DbExisted {
			if existed := dao.QueryChat(m.Chat.ID); existed != nil {
				flag = fmt.Sprintf(BIND_ALREADY, existed.Username)
			} else {
				flag = fmt.Sprintf(BIND_ALREADY, "?")
			}
		} else if status == xstatus.DbFailed {
			flag = BIND_FAILED
		} else if token != "" {
			flag = fmt.Sprintf(BIND_TOKEN_SUCCESS, username)
		} else {
			flag = fmt.Sprintf(BIND_NOTOK_SUCCESS, username)
		}
	}

	server.Bot().RemoveCache(m.Chat.ID, "username")
	server.Bot().SetStatus(m.Chat.ID, fsm.None)
	_ = server.Bot().Reply(m, flag)
}

// /unbind
func UnbindCtrl(m *telebot.Message) {
	user := dao.QueryChat(m.Chat.ID)
	if user == nil {
		_ = server.Bot().Reply(m, BIND_NOT_YET)
		return
	}

	flag := fmt.Sprintf(UNBIND_Q, user.Username)
	_ = server.Bot().Reply(m, flag, &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{*button.InlineBtnUnbind}, {*button.InlineBtnCancel},
		},
	})
}

// button.InlineBtnUnbind
func InlineBtnUnbindCtrl(c *telebot.Callback) {
	m := c.Message
	_, _ = server.Bot().Edit(m, fmt.Sprintf("%s (unbind)", m.Text))

	flag := ""
	status := dao.DeleteChat(m.Chat.ID)
	if status == xstatus.DbNotFound {
		flag = BIND_NOT_YET
	} else if status == xstatus.DbFailed {
		flag = UNBIND_FAILED
	} else {
		flag = UNBIND_SUCCESS
	}

	_ = server.Bot().Reply(m, flag)
}

// /me
func MeCtrl(m *telebot.Message) {
	user := dao.QueryChat(m.Chat.ID)
	if user == nil {
		_ = server.Bot().Reply(m, BIND_NOT_YET)
		return
	}

	flag := ""
	url := fmt.Sprintf("[%s](https://github.com/%s)", user.Username, user.Username)
	if user.Token == "" {
		flag = fmt.Sprintf(GITHUB_ME_NOTOK, url)
	} else {
		flag = fmt.Sprintf(GITHUB_ME_TOKEN, url, xstring.MaskTokenR(user.Token, 1, 2, 3, -1, -2, -3))
	}
	_ = server.Bot().Reply(m, flag, telebot.ModeMarkdown)
}

// /enablesilent
func EnableSilentCtrl(m *telebot.Message) {
	user := dao.QueryChat(m.Chat.ID)
	if user == nil {
		_ = server.Bot().Reply(m, BIND_NOT_YET)
		return
	}

	server.Bot().SetCache(m.Chat.ID, "chatID", user.ChatID)
	server.Bot().SetStatus(m.Chat.ID, fsm.EnablingSilent)
	_ = server.Bot().Reply(m, SILENT_Q)
}

// fsm.EnablingSilent
func FromEnablingSilentCtrl(m *telebot.Message) {
	sp := strings.Split(strings.TrimSpace(m.Text), " ")
	if len(sp) < 3 || sp[0] == sp[1] {
		_ = server.Bot().Reply(m, SILENT_UNEXPECTED_FORMAT)
		return
	}
	start, err1 := xnumber.Atoi(sp[0])
	end, err2 := xnumber.Atoi(sp[1])
	if (err1 != nil || start < 0 || start > 23) || (err2 != nil || end < 0 || end > 23) {
		_ = server.Bot().Reply(m, SILENT_UNEXPECTED_HOUR)
		return
	}
	zone := sp[2]
	loc, err := xtime.ParseTimezone(zone)
	if err != nil {
		_ = server.Bot().Reply(m, SILENT_UNEXPECTED_TIMEZONE)
		return
	}

	chatID, _ := server.Bot().GetCache(m.Chat.ID, "chatID")
	status := dao.UpdateChatSilent(chatID.(int64), true, start, end, zone)
	flag := ""
	if status == xstatus.DbNotFound {
		flag = BIND_NOT_YET
	} else if status == xstatus.DbFailed {
		flag = SILENT_FAILED
	} else {
		f := ""
		if start < end {
			f = fmt.Sprintf("%dh to %dh in %s", start, end, loc.String())
		} else {
			f = fmt.Sprintf("%dh to %dh next day in %s", start, end, loc.String())
		}
		flag = fmt.Sprintf(SILENT_SUCCESS, f)
	}

	server.Bot().RemoveCache(m.Chat.ID, "chatID")
	server.Bot().SetStatus(m.Chat.ID, fsm.None)
	_ = server.Bot().Reply(m, flag)
}

// /disablesilent
func DisableSilentCtrl(m *telebot.Message) {
	user := dao.QueryChat(m.Chat.ID)
	if user == nil {
		_ = server.Bot().Reply(m, BIND_NOT_YET)
		return
	}

	if !user.Silent {
		_ = server.Bot().Reply(m, SILENT_NOT_YET)
		return
	}

	status := dao.UpdateChatSilent(user.ChatID, false, 0, 0, user.TimeZone)
	flag := ""
	if status == xstatus.DbNotFound {
		flag = BIND_NOT_YET
	} else if status == xstatus.DbFailed {
		flag = DISABLE_SILENT_FAILED
	} else {
		flag = DISABLE_SILENT_SUCCESS
	}
	_ = server.Bot().Reply(m, flag)
}
