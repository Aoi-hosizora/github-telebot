package task

import (
	"github.com/Aoi-hosizora/github-telebot/internal/bot/server"
	"github.com/Aoi-hosizora/github-telebot/internal/model"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/config"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/dao"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/logger"
	"github.com/Aoi-hosizora/github-telebot/internal/service"
	"github.com/robfig/cron/v3"
	"gopkg.in/tucnak/telebot.v2"
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

	_cron = cr
	return nil
}

func activityTask() {
	defer func() { recover() }()

	users := dao.QueryUsers()
	if len(users) == 0 {
		return
	}

	foreachUsers(users, func(user *model.User) {
		// get new events
		resp, err := service.GetActivityEvents(user.Username, user.Token, 1)
		if err != nil {
			return
		}
		newEvents, err := model.UnmarshalActivityEvents(resp)
		if err != nil {
			return
		}

		// filter and check empty
		tempEvents := make([]*model.ActivityEvent, 0)
		for _, e := range newEvents {
			if !dao.CheckFilter(user.ChatID, e.Actor.Login, e.Repo.Name, e.Type) {
				tempEvents = append(tempEvents, e)
			}
		}
		newEvents = tempEvents
		if len(newEvents) == 0 {
			return
		}

		// get old events and calc diff
		oldEvents, ok := dao.GetActivities(user.ChatID)
		if !ok {
			return
		}
		logger.Logger().Infof("Get old ativities: #%d | (%d %s)", len(oldEvents), user.ChatID, user.Username)
		diff := model.ActivitySliceDiff(newEvents, oldEvents)
		logger.Logger().Infof("Get diff ativities: #%d | (%d %s)", len(diff), user.ChatID, user.Username)
		if len(diff) == 0 {
			return
		}

		// update old events
		ok = dao.SetActivities(user.ChatID, newEvents)
		if !ok {
			return
		}
		logger.Logger().Infof("Set new ativities: #%d | (%d %s)", len(newEvents), user.ChatID, user.Username)

		// render
		render := service.RenderActivityEvents(diff) // <<<
		if render == "" {
			return
		}
		flag := service.ConcatListAndUsername(render, user.Username) + " \\(Activity events\\)" // <<<

		// send
		if checkSilent(user) {
			_ = server.Bot().SendToChat(user.ChatID, flag, telebot.ModeMarkdownV2, telebot.Silent)
		} else {
			_ = server.Bot().SendToChat(user.ChatID, flag, telebot.ModeMarkdownV2)
		}
	})
}

func issueTask() {
	defer func() { recover() }()

	users := dao.QueryUsers()
	if len(users) == 0 {
		return
	}

	foreachUsers(users, func(user *model.User) {
		// allow to send issue
		if user.Token == "" || !user.AllowIssue {
			return
		}

		// get new events
		resp, err := service.GetIssueEvents(user.Username, user.Token, 1)
		if err != nil {
			return
		}
		newEvents, err := model.UnmarshalIssueEvents(resp)
		if err != nil {
			return
		}

		// filter and check empty
		tempEvents := make([]*model.IssueEvent, 0)
		for _, e := range newEvents {
			if !dao.CheckFilter(user.ChatID, e.Actor.Login, e.Repo, e.Event) {
				tempEvents = append(tempEvents, e)
			}
		}
		newEvents = tempEvents
		if user.FilterMe {
			tempEvents = make([]*model.IssueEvent, 0)
			for _, e := range newEvents {
				if e.Actor.Login != user.Username {
					tempEvents = append(tempEvents, e)
				}
			}
			newEvents = tempEvents
		}
		if len(newEvents) == 0 {
			return
		}

		// get old events and calc diff
		oldEvents, ok := dao.GetIssues(user.ChatID)
		if !ok {
			return
		}
		logger.Logger().Infof("Get old issues: #%d | (%d %s)", len(oldEvents), user.ChatID, user.Username)
		diff := model.IssueSliceDiff(newEvents, oldEvents)
		logger.Logger().Infof("Get diff issues: #%d | (%d %s)", len(diff), user.ChatID, user.Username)
		if len(diff) == 0 {
			return
		}

		// update old events
		ok = dao.SetIssues(user.ChatID, newEvents)
		if !ok {
			return
		}
		logger.Logger().Infof("Set new issues: #%d | (%d %s)", len(newEvents), user.ChatID, user.Username)

		// render
		render := service.RenderIssueEvents(diff) // <<<
		if render == "" {
			return
		}
		flag := service.ConcatListAndUsername(render, user.Username) + " \\(Issue events\\)" // <<<

		// send
		if checkSilent(user) {
			_ = server.Bot().SendToChat(user.ChatID, flag, telebot.ModeMarkdownV2, telebot.Silent)
		} else {
			_ = server.Bot().SendToChat(user.ChatID, flag, telebot.ModeMarkdownV2)
		}
	})
}
