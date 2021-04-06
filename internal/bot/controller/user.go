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
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/database"
	"github.com/Aoi-hosizora/github-telebot/internal/service"
	"gopkg.in/tucnak/telebot.v2"
	"strings"
)

// /bind
func BindCtrl(m *telebot.Message) {
	user := database.GetUser(m.Chat.ID)
	if user != nil {
		_ = server.Bot().Reply(m, BIND_ALREADY)
	} else {
		server.Bot().SetStatus(m.Chat.ID, fsm.Binding)
		_ = server.Bot().Reply(m, BIND_Q)
	}
}

// /bind -> x
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
		status := database.AddUser(user) // id username token
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
	user := database.GetUser(m.Chat.ID)
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
	user := database.GetUser(m.Chat.ID)
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

// inl:btn_unbind
func InlBtnUnbindCtrl(c *telebot.Callback) {
	m := c.Message
	_ = server.Bot().Delete(m)

	flag := ""
	status := database.DeleteUser(m.Chat.ID)
	if status == xstatus.DbNotFound {
		flag = BIND_NOT_YET
	} else if status == xstatus.DbFailed {
		flag = UNBIND_FAILED
	} else {
		flag = UNBIND_SUCCESS
	}

	_ = server.Bot().Reply(m, flag)
}

// /enablesilent
func EnableSilentCtrl(m *telebot.Message) {
	user := database.GetUser(m.Chat.ID)
	if user == nil {
		_ = server.Bot().Reply(m, BIND_NOT_YET)
		return
	}

	server.Bot().SetCache(m.Chat.ID, "user", user)
	server.Bot().SetStatus(m.Chat.ID, fsm.SilentHour)
	_ = server.Bot().Reply(m, SILENT_Q)
}

// /enablesilent -> x
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

	status := database.UpdateUserSilent(user.ChatID, true, start, end, zone)
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
	user := database.GetUser(m.Chat.ID)
	if user == nil {
		_ = server.Bot().Reply(m, BIND_NOT_YET)
		return
	}

	if !user.Silent {
		_ = server.Bot().Reply(m, SILENT_NOT_YET)
		return
	}

	status := database.UpdateUserSilent(user.ChatID, false, 0, 0, user.TimeZone)
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
