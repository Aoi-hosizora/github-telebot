package task

import (
	"fmt"
	"github.com/Aoi-hosizora/github-telebot/src/bot"
	"github.com/Aoi-hosizora/github-telebot/src/config"
	"github.com/Aoi-hosizora/github-telebot/src/model"
	"github.com/Aoi-hosizora/github-telebot/src/service"
	"time"
)

var (
	oldActivities = make(map[int64][]*model.ActivityEvent, 0)
	oldIssues     = make(map[int64][]*model.IssueEvent, 0)
)

func activityTask() {
	defer func() {
		if err := recover(); err != nil {
			activityTask()
		}
	}()

	for {
		users := model.GetUsers()
		for _, user := range users {
			// get event and unmarshal
			resp, err := service.GetActivityEvents(user.Username, user.Token, 1)
			if err != nil {
				continue
			}
			events, err := model.UnmarshalActivityEvents(resp)
			if err != nil {
				continue
			}

			// check map and diff
			if _, ok := oldActivities[user.ChatID]; !ok {
				oldActivities[user.ChatID] = []*model.ActivityEvent{}
			}
			diff := model.ActivitySliceDiff(events, oldActivities[user.ChatID])
			if len(diff) != 0 {
				// render and send
				render := service.RenderActivities(diff)
				flag := service.RenderResult(render, user.Username)
				_ = bot.SendToChat(user.ChatID, flag)
			}

			// update old map
			oldActivities[user.ChatID] = events
		}

		// wait to send next time
		time.Sleep(time.Duration(config.Configs.Task.ActivityDuration) * time.Second)
	}
}

func issueTask() {
	defer func() {
		if err := recover(); err != nil {
			issueTask()
		}
	}()

	for {
		users := model.GetUsers()
		for _, user := range users {
			// allow to send issue
			if user.Token == "" || !user.AllowIssue {
				continue
			}

			// get event and unmarshal
			resp, err := service.GetIssueEvents(user.Username, user.Token, 1)
			if err != nil {
				continue
			}
			events, err := model.UnmarshalIssueEvents(resp)
			if err != nil {
				continue
			}

			// check map and diff
			if _, ok := oldIssues[user.ChatID]; !ok {
				oldIssues[user.ChatID] = []*model.IssueEvent{}
			}
			diff := model.IssueSliceDiff(events, oldIssues[user.ChatID])
			if len(diff) != 0 {
				// render and send
				render := service.RenderIssues(diff)
				flag := fmt.Sprintf("%s\n---\nFrom [%s](https://github.com/%s) updated.", render, user.Username, user.Username)
				_ = bot.SendToChat(user.ChatID, flag)
			}

			// update old map
			oldIssues[user.ChatID] = events
		}

		// wait to send next time
		time.Sleep(time.Duration(config.Configs.Task.IssueDuration) * time.Second)
	}
}

func Start() {
	go activityTask()
	go issueTask()
}
