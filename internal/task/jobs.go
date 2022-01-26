package task

import (
	"context"
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
	"github.com/Aoi-hosizora/ahlib/xgopool"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/server"
	"github.com/Aoi-hosizora/github-telebot/internal/model"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/logger"
	"github.com/Aoi-hosizora/github-telebot/internal/service"
	"github.com/Aoi-hosizora/github-telebot/internal/service/dao"
	"gopkg.in/tucnak/telebot.v2"
	"sync"
)

type JobSet struct {
	bw   *xtelebot.BotWrapper
	pool *xgopool.GoPool
}

func NewJobSet(bw *xtelebot.BotWrapper, pool *xgopool.GoPool) *JobSet {
	return &JobSet{bw: bw, pool: pool}
}

func (j *JobSet) foreachUser(users []*model.Chat, fn func(user *model.Chat)) {
	wg := sync.WaitGroup{}
	for _, user := range users {
		wg.Add(1)
		user := user
		ctx := context.WithValue(context.Background(), ctxFuncnameKey, "foreachUser")
		j.pool.CtxGo(ctx, func(_ context.Context) {
			defer wg.Done()
			fn(user)
		})
	}
	wg.Wait()
}

func (j *JobSet) activityJob() {
	users, _ := dao.QueryChats()
	if len(users) == 0 {
		return
	}

	// foreach user
	j.foreachUser(users, func(user *model.Chat) {
		// get new events
		newEvents, _ := service.GetActivityEvents(user.Username, user.Token, 1)
		if len(newEvents) == 0 {
			return
		}

		// get old events and calc diff
		oldEvents, err := dao.GetActivities(user.ChatID)
		if err != nil {
			return
		}
		logger.Logger().Infof("Get old ativities: #%d | (%d %s)", len(oldEvents), user.ChatID, user.Username)
		diff := model.DiffActivityEventSlice(newEvents, oldEvents)
		logger.Logger().Infof("Get diff ativities: #%d | (%d %s)", len(diff), user.ChatID, user.Username)
		if len(diff) == 0 {
			return
		}

		// update old events
		err = dao.SetActivities(user.ChatID, newEvents)
		if err != nil {
			return
		}
		logger.Logger().Infof("Set new ativities: #%d | (%d %s)", len(newEvents), user.ChatID, user.Username)

		// format and send
		format := service.FormatActivityEvents(diff, user.Username, -1) // <<<
		if format == "" {
			return
		}
		opt := []interface{}{telebot.ModeMarkdownV2}
		if user.Silent {
			opt = append(opt, telebot.Silent)
		}
		if !user.Preview {
			opt = append(opt, telebot.NoPreview)
		}
		_ = server.Bot().SendToChat(user.ChatID, format, opt...)
	})
}

func (j *JobSet) issueJob() {
	users, _ := dao.QueryChats()
	if len(users) == 0 {
		return
	}

	// foreach user
	j.foreachUser(users, func(user *model.Chat) {
		// get new events and filter
		if user.Token == "" || !user.Issue {
			return
		}
		newEvents, _ := service.GetIssueEvents(user.Username, user.Token, 1)
		if user.FilterMe {
			newEvents = service.FilterIssueEventSlice(newEvents, user.Username)
		}
		if len(newEvents) == 0 {
			return
		}

		// get old events and calc diff
		oldEvents, err := dao.GetIssues(user.ChatID)
		if err != nil {
			return
		}
		logger.Logger().Infof("Get old issues: #%d | (%d %s)", len(oldEvents), user.ChatID, user.Username)
		diff := model.DiffIssueEventSlice(newEvents, oldEvents)
		logger.Logger().Infof("Get diff issues: #%d | (%d %s)", len(diff), user.ChatID, user.Username)
		if len(diff) == 0 {
			return
		}

		// update old events
		err = dao.SetIssues(user.ChatID, newEvents)
		if err != nil {
			return
		}
		logger.Logger().Infof("Set new issues: #%d | (%d %s)", len(newEvents), user.ChatID, user.Username)

		// format
		format := service.FormatIssueEvents(diff, user.Username, -1) // <<<
		if format == "" {
			return
		}
		opt := []interface{}{telebot.ModeMarkdownV2}
		if user.Silent {
			opt = append(opt, telebot.Silent)
		}
		if !user.Preview {
			opt = append(opt, telebot.NoPreview)
		}
		_ = server.Bot().SendToChat(user.ChatID, format, opt...)
	})
}
