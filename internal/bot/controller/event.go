package controller

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/fsm"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/server"
	"github.com/Aoi-hosizora/github-telebot/internal/model"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/dao"
	"github.com/Aoi-hosizora/github-telebot/internal/service"
	"gopkg.in/tucnak/telebot.v2"
)

const (
	SUPPORTED_EVENTS = "Supported activity and issue events:\n\n" +
		"1. Activity events: PushEvent, WatchEvent, CreateEvent, ForkEvent, DeleteEvent, PublicEvent, IssuesEvent, IssueCommentEvent, " +
		"PullRequestEvent, PullRequestReviewEvent, PullRequestReviewCommentEvent, CommitCommentEvent, MemberEvent, ReleaseEvent, GollumEvent.\n\n" +
		"2. Issue events: mentioned, opened, closed, reopened, renamed, labeled, unlabeled, locked, unlocked, milestoned, demilestoned, pinned, unpinned, " +
		"assigned, commented, merged, head_ref_deleted, head_ref_restored, added_to_project, removed_from_project, moved_columns_in_project, cross-referenced, referenced."

	GITHUB_PAGE_Q     = "Please send the page number you want to get. Send /cancel to cancel."
	UNEXPECTED_NUMBER = "Unexpected page number. Please send an integer value. Send /cancel to cancel."
	EMPTY_EVENT       = "You have empty event."
)

// /supportedevents
func SupportedEventCtrl(m *telebot.Message) {
	_ = server.Bot().Reply(m, SUPPORTED_EVENTS)
}

// /activity
func ActivityCtrl(m *telebot.Message) {
	m.Text = "1"
	FromActivityPageCtrl(m)
}

// /activitypage
func ActivityPageCtrl(m *telebot.Message) {
	server.Bot().SetStatus(m.Chat.ID, fsm.ActivityPage)
	_ = server.Bot().Reply(m, GITHUB_PAGE_Q)
}

// fsm.ActivityPage
func FromActivityPageCtrl(m *telebot.Message) {
	page, err := xnumber.Atoi(m.Text)
	if err != nil {
		_ = server.Bot().Reply(m, UNEXPECTED_NUMBER)
		return
	}
	if page <= 0 {
		page = 1
	}
	user := dao.QueryUser(m.Chat.ID)
	if user == nil {
		server.Bot().SetStatus(m.Chat.ID, fsm.None)
		_ = server.Bot().Reply(m, BIND_NOT_YET)
		return
	}

	flag := ""
	v2md := false
	if resp, err := service.GetActivityEvents(user.Username, user.Token, page); err != nil {
		flag = GITHUB_FAILED
	} else if events, err := model.UnmarshalActivityEvents(resp); err != nil {
		flag = GITHUB_FAILED
	} else {
		tempEvents := make([]*model.ActivityEvent, 0)
		for _, e := range events {
			if !dao.CheckFilter(user.ChatID, e.Repo.Name, e.Actor.Login, e.Type) {
				tempEvents = append(tempEvents, e)
			}
		}
		events = tempEvents
		if render := service.RenderActivityEvents(events); render == "" {
			flag = EMPTY_EVENT
		} else {
			flag = service.ConcatListAndUsername(render, user.Username) + fmt.Sprintf(" \\(page %d\\)", page) // <<<
			v2md = true
		}
	}

	server.Bot().SetStatus(m.Chat.ID, fsm.None)
	if !v2md {
		_ = server.Bot().Reply(m, flag, telebot.ModeMarkdown)
	} else {
		_ = server.Bot().Reply(m, flag, telebot.ModeMarkdownV2)
	}
}

// /issue
func IssueCtrl(m *telebot.Message) {
	m.Text = "1"
	FromIssuePageCtrl(m)
}

// /issuepage
func IssuePageCtrl(m *telebot.Message) {
	server.Bot().SetStatus(m.Chat.ID, fsm.IssuePage)
	_ = server.Bot().Reply(m, GITHUB_PAGE_Q)
}

// fsm.IssuePage
func FromIssuePageCtrl(m *telebot.Message) {
	page, err := xnumber.Atoi(m.Text)
	if err != nil {
		_ = server.Bot().Reply(m, UNEXPECTED_NUMBER)
		return
	}
	if page <= 0 {
		page = 1
	}
	user := dao.QueryUser(m.Chat.ID)
	if user == nil {
		server.Bot().SetStatus(m.Chat.ID, fsm.None)
		_ = server.Bot().Reply(m, BIND_NOT_YET)
		return
	}
	if user.Token == "" {
		_ = server.Bot().Reply(m, ISSUE_ONLY_FOR_TOKEN)
		return
	}

	flag := ""
	v2md := false
	if resp, err := service.GetIssueEvents(user.Username, user.Token, page); err != nil {
		flag = GITHUB_FAILED
	} else if events, err := model.UnmarshalIssueEvents(resp); err != nil {
		flag = GITHUB_FAILED
	} else if render := service.RenderIssueEvents(events); render == "" {
		flag = EMPTY_EVENT
	} else {
		flag = service.ConcatListAndUsername(render, user.Username) + fmt.Sprintf(" \\(page %d\\)", page) // <<<
		v2md = true
	}

	server.Bot().SetStatus(m.Chat.ID, fsm.None)
	if !v2md {
		_ = server.Bot().Reply(m, flag, telebot.ModeMarkdown)
	} else {
		_ = server.Bot().Reply(m, flag, telebot.ModeMarkdownV2)
	}
}
