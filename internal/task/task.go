package task

import (
	"github.com/Aoi-hosizora/github-telebot/internal/bot/server"
	"github.com/Aoi-hosizora/github-telebot/internal/model"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/config"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/database"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/logger"
	"github.com/Aoi-hosizora/github-telebot/internal/service"
	"github.com/robfig/cron/v3"
	"gopkg.in/tucnak/telebot.v2"
	"strings"
)

// _cron represents the global cron.Cron.
var _cron *cron.Cron

func Cron() *cron.Cron {
	return _cron
}

func Setup() error {
	cr := cron.New(cron.WithSeconds())

	_, err := cr.AddFunc(config.Configs().Task.Activity, activityTask)
	if err != nil {
		return err
	}

	_, err = cr.AddFunc(config.Configs().Task.Issue, issueTask)
	if err != nil {
		return err
	}

	return nil
}

func activityTask() {
	defer func() { recover() }()

	users := database.GetUsers()
	if len(users) == 0 {
		return
	}

	foreachUsers(users, func(user *model.User) {
		// get event and unmarshal
		resp, err := service.GetActivityEvents(user.Username, user.Token, 1)
		if err != nil {
			return
		}
		events, err := model.UnmarshalActivityEvents(resp)
		if err != nil {
			return
		}

		// check events and get diff
		oldEvents, ok := database.GetOldActivities(user.ChatID)
		if !ok {
			return
		}
		logger.Logger().Infof("Get old ativities: #%d | (%d %s)", len(oldEvents), user.ChatID, user.Username)
		diff := model.ActivitySliceDiff(events, oldEvents)
		logger.Logger().Infof("Get diff ativities: #%d | (%d %s)", len(diff), user.ChatID, user.Username)

		// update old events
		ok = database.SetOldActivities(user.ChatID, events)
		logger.Logger().Infof("Set new ativities: #%d | (%d %s)", len(events), user.ChatID, user.Username)
		if !ok {
			return
		}

		// render and send
		if len(diff) == 0 {
			return
		}
		render := service.RenderActivities(diff) // <<<
		if render == "" {
			return
		}
		flag := service.RenderResult(render, user.Username) + " \\(Activity events\\)" // <<<
		var sendErr error
		if checkSilent(user) {
			sendErr = server.Bot().SendToChat(user.ChatID, flag, telebot.ModeMarkdownV2, telebot.Silent)
		} else {
			sendErr = server.Bot().SendToChat(user.ChatID, flag, telebot.ModeMarkdownV2)
		}
		if sendErr != nil && strings.Contains(sendErr.Error(), "must be escaped") {
			flag = strings.ReplaceAll(flag, "\\", "")
			flag += "\n\nPlease contact with the developer with the message:\n" + sendErr.Error()
			if checkSilent(user) {
				_ = server.Bot().SendToChat(user.ChatID, flag, telebot.ModeMarkdown, telebot.Silent)
			} else {
				_ = server.Bot().SendToChat(user.ChatID, flag, telebot.ModeMarkdown)
			}
		}
	})
}

func issueTask() {
	defer func() { recover() }()

	users := database.GetUsers()
	if len(users) == 0 {
		return
	}

	foreachUsers(users, func(user *model.User) {
		// allow to send issue
		if user.Token == "" || !user.AllowIssue {
			return
		}

		// get event and unmarshal
		resp, err := service.GetIssueEvents(user.Username, user.Token, 1)
		if err != nil {
			return
		}
		events, err := model.UnmarshalIssueEvents(resp)
		if err != nil {
			return
		}
		if user.FilterMe {
			tempEvents := make([]*model.IssueEvent, 0)
			for _, e := range events {
				if e.Actor.Login != user.Username {
					tempEvents = append(tempEvents, e)
				}
			}
			events = tempEvents
		}

		// check events and get diff
		oldEvents, ok := database.GetOldIssues(user.ChatID)
		logger.Logger().Infof("Get old issues: #%d | (%d %s)", len(oldEvents), user.ChatID, user.Username)
		diff := model.IssueSliceDiff(events, oldEvents)
		logger.Logger().Infof("Get diff issues: #%d | (%d %s)", len(diff), user.ChatID, user.Username)

		// update old events
		ok = database.SetOldIssues(user.ChatID, events)
		logger.Logger().Infof("Set new issues: #%d | (%d %s)", len(events), user.ChatID, user.Username)
		if !ok {
			return
		}

		// render and send
		if len(diff) == 0 {
			return
		}
		render := service.RenderIssues(diff) // <<<
		if render == "" {
			return
		}
		flag := service.RenderResult(render, user.Username) + " \\(Issue events\\)" // <<<
		var sendErr error
		if checkSilent(user) {
			sendErr = server.Bot().SendToChat(user.ChatID, flag, telebot.ModeMarkdownV2, telebot.Silent)
		} else {
			sendErr = server.Bot().SendToChat(user.ChatID, flag, telebot.ModeMarkdownV2)
		}
		if sendErr != nil && strings.Contains(sendErr.Error(), "must be escaped") {
			flag = strings.ReplaceAll(flag, "\\", "")
			flag += "\n\nPlease contact with the developer with the message:\n" + sendErr.Error()
			if checkSilent(user) {
				_ = server.Bot().SendToChat(user.ChatID, flag, telebot.ModeMarkdown, telebot.Silent)
			} else {
				_ = server.Bot().SendToChat(user.ChatID, flag, telebot.ModeMarkdown)
			}
		}
	})
}
