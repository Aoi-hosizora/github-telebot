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
	"github.com/Aoi-hosizora/github-telebot/internal/model"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/dao"
	"github.com/Aoi-hosizora/github-telebot/internal/service"
	"gopkg.in/tucnak/telebot.v2"
	"strings"
)

const (
	BIND_Q             = "Please send your github's username, and if you want to watch private events or issue events, please send your token after username (split by whitespace). /cancel to cancel."
	BIND_ALREADY       = "You have already bound with a github account."
	BIND_NOT_YET       = "You have not bound a github account yet."
	BIND_EMPTY         = "Please resend a non-empty username again."
	BIND_FAILED        = "Failed to bind github account, please retry later."
	BIND_SUCCESS       = "Binding user '%s' without token success. /activity to get activity events, /issue to get issue events.\n" + BIND_SUCCESS_TIP
	BIND_TOKEN_SUCCESS = "Binding user '%s' with token success. /activity to get events, /issue to get issue events.\n" + BIND_SUCCESS_TIP
	BIND_SUCCESS_TIP   = "(Tips: new activity events will be sent periodically, but issue events will not be sent in default. Use /allowissue to allow)"

	UNBIND_Q       = "Sure to unbind the current github account '%s'?"
	UNBIND_FAILED  = "Failed to unbind github account, please retry later."
	UNBIND_SUCCESS = "Unbind user success."

	SILENT_Q               = "Please send 2 different numbers as hour (in [0, 23]) you want to start and finish silent send, and with a timezone (such as +8:00 or -06:30), split by whitespace. Examples: 23 6 +8 or 0 8 -6."
	SILENT_FORMAT_REQUIRED = "Excepted input, please send 2 different numbers as hour (in [0, 23]) you want to start and finish silent send, and with a timezone."
	SILENT_HOUR_REQUIRED   = "Excepted hour, please send an integer in [0, 23]."
	SILENT_ZONE_REQUIRED   = "Excepted timezone, please send a right time zone, such as +8:00 or -06:30"
	SILENT_NOT_YET         = "You have not set silent yet, use /enablesilent to set."
	SILENT_SUCCESS         = "Success. Now it will be silent when send message in %s."
	SILENT_FAILED          = "Failed to set silent send, please retry later."
	DISABLE_SILENT_SUCCESS = "Disable silent success. Any message will be sent directly now."
	DISABLE_SILENT_FAILED  = "Failed to disable silent, please retry later."
)

// /bind
func BindCtrl(m *telebot.Message) {
	user := dao.QueryUser(m.Chat.ID)
	if user != nil {
		_ = server.Bot().Reply(m, BIND_ALREADY)
	} else {
		server.Bot().SetStatus(m.Chat.ID, fsm.Binding)
		_ = server.Bot().Reply(m, BIND_Q)
	}
}

// fsm.Binding
func FromBindingCtrl(m *telebot.Message) {
	text := strings.TrimSpace(m.Text)
	if text == "" {
		_ = server.Bot().Reply(m, BIND_EMPTY)
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
		status := dao.CreateUser(user) // id username token
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

	server.Bot().SetStatus(m.Chat.ID, fsm.None)
	_ = server.Bot().Reply(m, flag)
}

// /me
func MeCtrl(m *telebot.Message) {
	user := dao.QueryUser(m.Chat.ID)
	flag := ""
	if user == nil {
		flag = BIND_NOT_YET
	} else {
		name := service.Markdown(user.Username)
		n := fmt.Sprintf("[%s](https://github.com/%s)", name, user.Username)
		if user.Token != "" {
			flag = fmt.Sprintf(GITHUB_ME_TOKEN, n, xstring.DefaultMaskToken(user.Token))
		} else {
			flag = fmt.Sprintf(GITHUB_ME, n)
		}
	}
	_ = server.Bot().Reply(m, flag, telebot.ModeMarkdown)
}

// /unbind
func UnbindCtrl(m *telebot.Message) {
	user := dao.QueryUser(m.Chat.ID)
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
	_ = server.Bot().Delete(m)

	flag := ""
	status := dao.DeleteUser(m.Chat.ID)
	if status == xstatus.DbNotFound {
		flag = BIND_NOT_YET
	} else if status == xstatus.DbFailed {
		flag = UNBIND_FAILED
	} else {
		flag = UNBIND_SUCCESS
	}

	_ = server.Bot().Reply(m, flag)
}

// button.InlineBtnCancel
func InlineBtnCancelCtrl(c *telebot.Callback) {
	m := c.Message
	_, _ = server.Bot().Edit(m, fmt.Sprintf("%s (canceled)", m.Text))
}

// /enablesilent
func EnableSilentCtrl(m *telebot.Message) {
	user := dao.QueryUser(m.Chat.ID)
	if user == nil {
		_ = server.Bot().Reply(m, BIND_NOT_YET)
		return
	}

	server.Bot().SetCache(m.Chat.ID, "user", user)
	server.Bot().SetStatus(m.Chat.ID, fsm.SilentHour)
	_ = server.Bot().Reply(m, SILENT_Q)
}

// fsm.SilentHour
func FromSilentHourCtrl(m *telebot.Message) {
	userItf, ok := server.Bot().GetCache(m.Chat.ID, "user")
	if !ok {
		_ = server.Bot().Reply(m, BIND_NOT_YET)
		return
	}
	user := userItf.(*model.User)
	sp := strings.Split(strings.TrimSpace(m.Text), " ")
	if len(sp) < 3 || sp[0] == sp[1] {
		_ = server.Bot().Reply(m, SILENT_FORMAT_REQUIRED)
		return
	}
	start, err1 := xnumber.ParseInt(sp[0], 10)
	end, err2 := xnumber.ParseInt(sp[1], 10)
	if (err1 != nil || start < 0 || start > 23) || (err2 != nil || end < 0 || end > 23) {
		_ = server.Bot().Reply(m, SILENT_HOUR_REQUIRED)
		return
	}
	zone := sp[2]
	loc, err := xtime.ParseTimezone(zone)
	if err != nil {
		_ = server.Bot().Reply(m, SILENT_ZONE_REQUIRED)
		return
	}

	status := dao.UpdateUserSilent(user.ChatID, true, start, end, zone)
	flag := ""
	if status == xstatus.DbNotFound {
		flag = BIND_NOT_YET
	} else if status == xstatus.DbFailed {
		flag = SILENT_FAILED
	} else {
		f := ""
		if start < end {
			f = fmt.Sprintf("%dh to %dh for %s", start, end, loc.String())
		} else {
			f = fmt.Sprintf("%dh to %dh next day for %s", start, end, loc.String())
		}
		flag = fmt.Sprintf(SILENT_SUCCESS, f)
	}

	server.Bot().RemoveCache(m.Chat.ID, "user")
	server.Bot().SetStatus(m.Chat.ID, fsm.None)
	_ = server.Bot().Reply(m, flag)
}

// /disablesilent
func DisableSilentCtrl(m *telebot.Message) {
	user := dao.QueryUser(m.Chat.ID)
	if user == nil {
		_ = server.Bot().Reply(m, BIND_NOT_YET)
		return
	}

	if !user.Silent {
		_ = server.Bot().Reply(m, SILENT_NOT_YET)
		return
	}

	status := dao.UpdateUserSilent(user.ChatID, false, 0, 0, user.TimeZone)
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
