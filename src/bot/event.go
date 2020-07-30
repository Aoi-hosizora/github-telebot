package bot

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xstatus"
	"github.com/Aoi-hosizora/github-telebot/src/bot/fsm"
	"github.com/Aoi-hosizora/github-telebot/src/model"
	"github.com/Aoi-hosizora/github-telebot/src/service"
	"gopkg.in/tucnak/telebot.v2"
	"strconv"
)

// /allowissue
func AllowIssueCtrl(m *telebot.Message) {
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

// /disallowissue
func DisallowIssueCtrl(m *telebot.Message) {
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
func ActivityCtrl(m *telebot.Message) {
	m.Text = "1"
	fromActivityNCtrl(m)
}

// /activityn
func ActivityNCtrl(m *telebot.Message) {
	Bot.UserStates[m.Chat.ID] = fsm.ActivityN
	_ = Bot.Reply(m, GITHUB_SEND_PAGE)
}

// /activityn -> x
func fromActivityNCtrl(m *telebot.Message) {
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
			flag = service.RenderResult(render, user.Username) + fmt.Sprintf(" (page %d)", page)
		}
	}

	Bot.UserStates[m.Chat.ID] = fsm.None
	_ = Bot.Reply(m, flag, telebot.ModeMarkdown)
}

// /issue
func IssueCtrl(m *telebot.Message) {
	m.Text = "1"
	fromIssueNCtrl(m)
}

// /issuen
func IssueNCtrl(m *telebot.Message) {
	Bot.UserStates[m.Chat.ID] = fsm.IssueN
	_ = Bot.Reply(m, GITHUB_SEND_PAGE)
}

// /issuen -> x
func fromIssueNCtrl(m *telebot.Message) {
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
	} else if user.Token == "" {
		flag = ISSUE_ONLY_FOR_TOKEN
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
