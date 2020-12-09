package controller

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/Aoi-hosizora/github-telebot/src/bot/fsm"
	"github.com/Aoi-hosizora/github-telebot/src/bot/server"
	"github.com/Aoi-hosizora/github-telebot/src/database"
	"github.com/Aoi-hosizora/github-telebot/src/model"
	"github.com/Aoi-hosizora/github-telebot/src/service"
	"gopkg.in/tucnak/telebot.v2"
	"strings"
)

// /allowissue
func AllowIssueCtrl(m *telebot.Message) {
	user := database.GetUser(m.Chat.ID)
	if user == nil {
		_ = server.Bot.Reply(m, BIND_NOT_YET)
		return
	} else if user.Token == "" {
		_ = server.Bot.Reply(m, ISSUE_ONLY_FOR_TOKEN)
		return
	}

	flag := ""
	user.AllowIssue = true
	status := database.UpdateUser(user)
	if status == xstatus.DbNotFound {
		flag = BIND_NOT_YET
	} else if status == xstatus.DbFailed {
		flag = ISSUE_ALLOW_FAILED
	} else {
		flag = ISSUE_ALLOW_SUCCESS
	}

	_ = server.Bot.Reply(m, flag)
}

// /disallowissue
func DisallowIssueCtrl(m *telebot.Message) {
	user := database.GetUser(m.Chat.ID)
	if user == nil {
		_ = server.Bot.Reply(m, BIND_NOT_YET)
		return
	} else if user.Token == "" {
		_ = server.Bot.Reply(m, ISSUE_ONLY_FOR_TOKEN)
		return
	}

	flag := ""
	user.AllowIssue = false
	status := database.UpdateUser(user)
	if status == xstatus.DbNotFound {
		flag = BIND_NOT_YET
	} else if status == xstatus.DbFailed {
		flag = ISSUE_DISALLOW_FAILED
	} else {
		flag = ISSUE_DISALLOW_SUCCESS
	}

	_ = server.Bot.Reply(m, flag)
}

// /activity
func ActivityCtrl(m *telebot.Message) {
	m.Text = "1"
	fromActivityNCtrl(m)
}

// /activityn
func ActivityNCtrl(m *telebot.Message) {
	server.Bot.UsersData.SetStatus(m.Chat.ID, fsm.ActivityPage)
	_ = server.Bot.Reply(m, GITHUB_SEND_PAGE_Q)
}

// /activityn -> x
func fromActivityNCtrl(m *telebot.Message) {
	page, err := xnumber.Atoi(m.Text)
	if err != nil {
		_ = server.Bot.Reply(m, NUM_REQUIRED)
		return
	} else if page <= 0 {
		page = 1
	}

	flag := ""
	v2 := false
	user := database.GetUser(m.Chat.ID)
	if user == nil {
		flag = BIND_NOT_YET
	} else {
		if resp, err := service.GetActivityEvents(user.Username, user.Token, page); err != nil {
			flag = GITHUB_FAILED
		} else if events, err := model.UnmarshalActivityEvents(resp); err != nil {
			flag = GITHUB_FAILED
		} else if render := service.RenderActivities(events); render == "" {
			flag = GITHUB_EMPTY
		} else {
			flag = service.RenderResult(render, user.Username) + fmt.Sprintf(" \\(page %d\\)", page) // <<<<<<
			v2 = true
		}
	}

	server.Bot.UsersData.SetStatus(m.Chat.ID, fsm.None)
	if !v2 {
		_ = server.Bot.Reply(m, flag, telebot.ModeMarkdown)
	} else {
		err := server.Bot.Reply(m, flag, telebot.ModeMarkdownV2)
		if err != nil && strings.Contains(err.Error(), "must be escaped") {
			flag = strings.ReplaceAll(flag, "\\", "")
			flag += "\n\nPlease contact with the developer with the message:\n" + err.Error()
			_ = server.Bot.Reply(m, flag, telebot.ModeMarkdown)
		}
	}
}

// /issue
func IssueCtrl(m *telebot.Message) {
	m.Text = "1"
	fromIssueNCtrl(m)
}

// /issuen
func IssueNCtrl(m *telebot.Message) {
	server.Bot.UsersData.SetStatus(m.Chat.ID, fsm.IssuePage)
	_ = server.Bot.Reply(m, GITHUB_SEND_PAGE_Q)
}

// /issuen -> x
func fromIssueNCtrl(m *telebot.Message) {
	page, err := xnumber.Atoi(m.Text)
	if err != nil {
		_ = server.Bot.Reply(m, NUM_REQUIRED)
		return
	} else if page <= 0 {
		page = 1
	}

	v2 := false
	flag := ""
	user := database.GetUser(m.Chat.ID)
	if user == nil {
		flag = BIND_NOT_YET
	} else if user.Token == "" {
		flag = ISSUE_ONLY_FOR_TOKEN
	} else {
		if resp, err := service.GetIssueEvents(user.Username, user.Token, page); err != nil {
			flag = GITHUB_FAILED
		} else if events, err := model.UnmarshalIssueEvents(resp); err != nil {
			flag = GITHUB_FAILED
		} else if render := service.RenderIssues(events); render == "" {
			flag = GITHUB_EMPTY
		} else {
			v2 = true
			flag = service.RenderResult(render, user.Username) + fmt.Sprintf(" \\(page %d\\)", page) // <<<<<<
		}
	}

	server.Bot.UsersData.SetStatus(m.Chat.ID, fsm.None)
	if !v2 {
		_ = server.Bot.Reply(m, flag, telebot.ModeMarkdown)
	} else {
		err := server.Bot.Reply(m, flag, telebot.ModeMarkdownV2)
		if err != nil && strings.Contains(err.Error(), "must be escaped") {
			flag = strings.ReplaceAll(flag, "\\", "")
			flag += "\n\nPlease contact with the developer with the message:\n" + err.Error()
			_ = server.Bot.Reply(m, flag, telebot.ModeMarkdown)
		}
	}
}
