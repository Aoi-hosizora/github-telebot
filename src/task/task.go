package task

import (
	"github.com/Aoi-hosizora/github-telebot/src/bot"
	"github.com/Aoi-hosizora/github-telebot/src/config"
	"github.com/Aoi-hosizora/github-telebot/src/model"
	"github.com/Aoi-hosizora/github-telebot/src/service"
	"sync"
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
		if len(users) > 0 {
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
					events = model.ReverseActivitySlice(events)

					// check map and get diff
					if _, ok := oldActivities[user.ChatID]; !ok {
						oldActivities[user.ChatID] = []*model.ActivityEvent{}
					}
					diff := model.ActivitySliceDiff(events, oldActivities[user.ChatID])

					// render and send
					if len(diff) != 0 {
						render := service.RenderActivities(diff)
						flag := service.RenderResult(render, user.Username)
						_ = bot.SendToChat(user.ChatID, flag)
					}

					// update old map
					oldActivities[user.ChatID] = events
					wg.Done()
				}(user)
				// ------------------------------------------------------------------------------------------------------ //
			}
			wg.Wait()
		}

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
		if len(users) > 0 {
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
					events = model.ReverseIssueSlice(events)

					// check map and get diff
					if _, ok := oldIssues[user.ChatID]; !ok {
						oldIssues[user.ChatID] = []*model.IssueEvent{}
					}
					diff := model.IssueSliceDiff(events, oldIssues[user.ChatID])

					// render and send
					if len(diff) != 0 {
						render := service.RenderIssues(diff)
						flag := service.RenderResult(render, user.Username)
						_ = bot.SendToChat(user.ChatID, flag)
					}

					// update old map
					oldIssues[user.ChatID] = events
					wg.Done()
				}(user)
				// ------------------------------------------------------------------------------------------------------ //
			}
			wg.Wait()
		}

		// wait to send next time
		time.Sleep(time.Duration(config.Configs.Task.IssueDuration) * time.Second)
	}
}

func Start() {
	go activityTask()
	go issueTask()
}
