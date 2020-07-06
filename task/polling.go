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
	old = make(map[int64][]*model.GithubEvent, 0)
)

func sliceDiff(s1 []*model.GithubEvent, s2 []*model.GithubEvent) []*model.GithubEvent {
	result := make([]*model.GithubEvent, 0)
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

func Task() {
	for {
		users := model.GetUsers()
		for _, user := range users {
			resp, err := util.GetGithubEvents(user.Username, user.Private, user.Token, 1)
			if err != nil {
				continue
			}
			events, err := model.UnmarshalEvents(resp)
			if err != nil {
				continue
			}

			if _, ok := old[user.ChatID]; !ok {
				old[user.ChatID] = []*model.GithubEvent{}
			}
			diff := sliceDiff(events, old[user.ChatID])
			if len(diff) != 0 {
				render := util.RenderGithubActions(diff)
				flag := fmt.Sprintf("%s\n---\nFrom [%s](https://github.com/%s) updated.", render, user.Username, user.Username)
				bot.SendToChat(user.ChatID, flag)
			}
			old[user.ChatID] = events
		}
		time.Sleep(time.Duration(config.Configs.TaskConfig.PollingDuration) * time.Second)
	}
}
