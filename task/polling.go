package task

import (
	"fmt"
	"github.com/Aoi-hosizora/ah-tgbot/bot"
	"github.com/Aoi-hosizora/ah-tgbot/config"
	"github.com/Aoi-hosizora/ah-tgbot/model"
	"github.com/Aoi-hosizora/ah-tgbot/util"
	"reflect"
	"time"
)

var (
	oldActivities = make(map[int64][]*model.ActivityEvent, 0)
	oldIssues     = make(map[int64][]*model.IssueEvent, 0)
)

func sliceActivityDiff(s1 []*model.ActivityEvent, s2 []*model.ActivityEvent) []*model.ActivityEvent {
	result := make([]*model.ActivityEvent, 0)
	for _, item1 := range s1 {
		exist := false
		for _, item2 := range s2 {
			if reflect.DeepEqual(item1, item2) {
				exist = true
				break
			}
		}
		if !exist {
			result = append(result, item1)
		}
	}
	return result
}

func sliceIssueDiff(s1 []*model.IssueEvent, s2 []*model.IssueEvent) []*model.IssueEvent {
	result := make([]*model.IssueEvent, 0)
	for _, item1 := range s1 {
		exist := false
		for _, item2 := range s2 {
			if reflect.DeepEqual(item1, item2) {
				exist = true
				break
			}
		}
		if !exist {
			result = append(result, item1)
		}
	}
	return result
}

func ActivityTask() {
	for {
		users := model.GetUsers()
		for _, user := range users {
			resp, err := util.GetGithubActivityEvents(user.Username, user.Private, user.Token, 1)
			if err != nil {
				continue
			}
			events, err := model.UnmarshalActivityEvents(resp)
			if err != nil {
				continue
			}

			if _, ok := oldActivities[user.ChatID]; !ok {
				oldActivities[user.ChatID] = []*model.ActivityEvent{}
			}
			diff := sliceActivityDiff(events, oldActivities[user.ChatID])
			if len(diff) != 0 {
				render := util.RenderGithubActivityString(diff)
				flag := fmt.Sprintf("%s\n---\nFrom [%s](https://github.com/%s) updated.", render, user.Username, user.Username)
				bot.SendToChat(user.ChatID, flag)
			}
			oldActivities[user.ChatID] = events
		}
		time.Sleep(time.Duration(config.Configs.TaskConfig.PollingActivityDuration) * time.Second)
	}
}

func IssueTask() {
	for {
		users := model.GetUsers()
		for _, user := range users {
			if !user.Private {
				continue
			}

			resp, err := util.GetGithubIssueEvents(user.Username, user.Private, user.Token, 1)
			if err != nil {
				continue
			}
			events, err := model.UnmarshalIssueEvents(resp)
			if err != nil {
				continue
			}

			if _, ok := oldActivities[user.ChatID]; !ok {
				oldActivities[user.ChatID] = []*model.ActivityEvent{}
			}
			diff := sliceIssueDiff(events, oldIssues[user.ChatID])
			if len(diff) != 0 {
				render := util.RenderGithubIssueString(diff)
				flag := fmt.Sprintf("%s\n---\nFrom [%s](https://github.com/%s) updated.", render, user.Username, user.Username)
				bot.SendToChat(user.ChatID, flag)
			}
			oldIssues[user.ChatID] = events
		}
		time.Sleep(time.Duration(config.Configs.TaskConfig.PollingIssueDuration) * time.Second)
	}
}
