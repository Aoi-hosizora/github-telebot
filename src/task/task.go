package task

import (
	"github.com/Aoi-hosizora/github-telebot/src/bot"
	"github.com/Aoi-hosizora/github-telebot/src/config"
	"github.com/Aoi-hosizora/github-telebot/src/database"
	"github.com/Aoi-hosizora/github-telebot/src/model"
	"github.com/Aoi-hosizora/github-telebot/src/service"
	"github.com/robfig/cron/v3"
	"gopkg.in/tucnak/telebot.v2"
	"sync"
)

var Cron *cron.Cron

func Setup() error {
	Cron = cron.New(cron.WithSeconds())

	_, err := Cron.AddFunc(config.Configs.Task.Activity, activityTask)
	if err != nil {
		return err
	}

	_, err = Cron.AddFunc(config.Configs.Task.Issue, issueTask)
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

	wg := sync.WaitGroup{}
	wg.Add(len(users))
	for _, user := range users {
		// ------------------------------------------------------------------------------------------------------ //
		go func(user *model.User) {
			// get event and unmarshal
			resp, err := service.GetActivityEvents(user.Username, user.Token, 1)
			if err != nil {
				wg.Done()
				return
			}
			events, err := model.UnmarshalActivityEvents(resp)
			if err != nil {
				wg.Done()
				return
			}

			// check events and get diff
			oldEvents, ok := database.GetOldActivities(user.ChatID)
			if !ok {
				wg.Done()
				return
			}
			diff := model.ActivitySliceDiff(events, oldEvents)

			// update old events
			ok = database.SetOldActivities(user.ChatID, events)
			if !ok {
				wg.Done()
				return
			}

			// render and send
			if len(diff) != 0 {
				render := service.RenderActivities(diff)
				if render != "" {
					flag := service.RenderResult(render, user.Username)
					_ = bot.SendToChat(user.ChatID, flag, telebot.ModeMarkdown)
				}
			}

			wg.Done()
		}(user)
		// ------------------------------------------------------------------------------------------------------ //
	}
	wg.Wait()
}

func issueTask() {
	defer func() { recover() }()

	users := database.GetUsers()
	if len(users) == 0 {
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(users))
	for _, user := range users {
		// ------------------------------------------------------------------------------------------------------ //
		go func(user *model.User) {
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

			// check events and get diff
			oldEvents, ok := database.GetOldIssues(user.ChatID)
			if !ok {
				wg.Done()
				return
			}
			diff := model.IssueSliceDiff(events, oldEvents)

			// update old events
			ok = database.SetOldIssues(user.ChatID, events)
			if !ok {
				wg.Done()
				return
			}

			// render and send
			if len(diff) != 0 {
				render := service.RenderIssues(diff)
				if render != "" {
					flag := service.RenderResult(render, user.Username)
					_ = bot.SendToChat(user.ChatID, flag, telebot.ModeMarkdown)
				}
			}

			wg.Done()
		}(user)
		// ------------------------------------------------------------------------------------------------------ //
	}
	wg.Wait()
}
