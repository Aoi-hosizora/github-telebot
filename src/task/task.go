package task

import (
	"github.com/Aoi-hosizora/github-telebot/src/bot"
	"github.com/Aoi-hosizora/github-telebot/src/config"
	"github.com/Aoi-hosizora/github-telebot/src/database"
	"github.com/Aoi-hosizora/github-telebot/src/model"
	"github.com/Aoi-hosizora/github-telebot/src/service"
	"gopkg.in/tucnak/telebot.v2"
	"sync"
	"time"
)

func activityTask() {
	defer func() {
		if err := recover(); err != nil {
			activityTask()
		}
	}()

	for {
		users := database.GetUsers()
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
		users := database.GetUsers()
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

		time.Sleep(time.Duration(config.Configs.Task.IssueDuration) * time.Second)
	}
}

func Start() {
	go activityTask()
	go issueTask()
}
