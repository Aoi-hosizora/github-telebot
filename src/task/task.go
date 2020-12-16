package task

import (
	"github.com/Aoi-hosizora/ahlib/xzone"
	"github.com/Aoi-hosizora/github-telebot/src/bot"
	"github.com/Aoi-hosizora/github-telebot/src/config"
	"github.com/Aoi-hosizora/github-telebot/src/database"
	"github.com/Aoi-hosizora/github-telebot/src/logger"
	"github.com/Aoi-hosizora/github-telebot/src/model"
	"github.com/Aoi-hosizora/github-telebot/src/service"
	"github.com/robfig/cron/v3"
	"gopkg.in/tucnak/telebot.v2"
	"strings"
	"sync"
	"time"
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

func checkSilent(user *model.User) bool {
	if user.Silent {
		hm, _ := xzone.MoveToZone(time.Now(), user.TimeZone)
		ss := user.SilentStart
		se := user.SilentEnd
		hour := hm.Hour()
		if ss < se { // 2 5
			if hour >= ss && hour <= se {
				return true
			}
		} else { // 22 2
			if (hour >= ss && hour <= 23) || (hour >= 0 && hour <= se) {
				return true
			}
		}
	}
	return false
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
			logger.Logger.Infof("Get old ativities: #%d | (%d %s)", len(oldEvents), user.ChatID, user.Username)
			diff := model.ActivitySliceDiff(events, oldEvents)
			logger.Logger.Infof("Get diff ativities: #%d | (%d %s)", len(diff), user.ChatID, user.Username)

			// update old events
			ok = database.SetOldActivities(user.ChatID, events)
			logger.Logger.Infof("Set new ativities: #%d | (%d %s)", len(events), user.ChatID, user.Username)
			if !ok {
				wg.Done()
				return
			}

			// render and send
			if len(diff) != 0 {
				render := service.RenderActivities(diff) // <<<
				if render != "" {
					flag := service.RenderResult(render, user.Username) + " \\(Activity events\\)" // <<<
					var sendErr error
					if checkSilent(user) {
						sendErr = bot.SendToChat(user.ChatID, flag, telebot.ModeMarkdownV2, telebot.Silent)
					} else {
						sendErr = bot.SendToChat(user.ChatID, flag, telebot.ModeMarkdownV2)
					}
					if sendErr != nil && strings.Contains(sendErr.Error(), "must be escaped") {
						flag = strings.ReplaceAll(flag, "\\", "")
						flag += "\n\nPlease contact with the developer with the message:\n" + sendErr.Error()
						if checkSilent(user) {
							_ = bot.SendToChat(user.ChatID, flag, telebot.ModeMarkdown, telebot.Silent)
						} else {
							_ = bot.SendToChat(user.ChatID, flag, telebot.ModeMarkdown)
						}
					}
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
			logger.Logger.Infof("Get old ativities: #%d | (%d %s)", len(oldEvents), user.ChatID, user.Username)
			diff := model.IssueSliceDiff(events, oldEvents)
			logger.Logger.Infof("Get diff ativities: #%d | (%d %s)", len(diff), user.ChatID, user.Username)

			// update old events
			ok = database.SetOldIssues(user.ChatID, events)
			logger.Logger.Infof("Set new ativities: #%d | (%d %s)", len(events), user.ChatID, user.Username)
			if !ok {
				wg.Done()
				return
			}

			// render and send
			if len(diff) != 0 {
				render := service.RenderIssues(diff) // <<<
				if render != "" {
					flag := service.RenderResult(render, user.Username) + " \\(Issue events\\)" // <<<
					var sendErr error
					if checkSilent(user) {
						sendErr = bot.SendToChat(user.ChatID, flag, telebot.ModeMarkdownV2, telebot.Silent)
					} else {
						sendErr = bot.SendToChat(user.ChatID, flag, telebot.ModeMarkdownV2)
					}
					if sendErr != nil && strings.Contains(sendErr.Error(), "must be escaped") {
						flag = strings.ReplaceAll(flag, "\\", "")
						flag += "\n\nPlease contact with the developer with the message:\n" + sendErr.Error()
						if checkSilent(user) {
							_ = bot.SendToChat(user.ChatID, flag, telebot.ModeMarkdown, telebot.Silent)
						} else {
							_ = bot.SendToChat(user.ChatID, flag, telebot.ModeMarkdown)
						}
					}
				}
			}

			wg.Done()
		}(user)
		// ------------------------------------------------------------------------------------------------------ //
	}
	wg.Wait()
}
