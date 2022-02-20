package task

import (
	"context"
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
	"github.com/Aoi-hosizora/ahlib/xgopool"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/github-telebot/internal/model"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/logger"
	"github.com/Aoi-hosizora/github-telebot/internal/service"
	"github.com/Aoi-hosizora/github-telebot/internal/service/dao"
	"gopkg.in/tucnak/telebot.v2"
	"sync"
	"sync/atomic"
)

type JobSet struct {
	bw   *xtelebot.BotWrapper
	pool *xgopool.GoPool

	activityFlag int32
	issueFlag    int32
}

func NewJobSet(bw *xtelebot.BotWrapper, pool *xgopool.GoPool) *JobSet {
	return &JobSet{bw: bw, pool: pool}
}

func (j *JobSet) checkExclusive(flag *int32) (allow bool, release func()) {
	if !atomic.CompareAndSwapInt32(flag, 0, 1) {
		return false, nil
	}
	return true, func() {
		*flag = 0
	}
}

func (j *JobSet) foreachChat(chats []*model.Chat, fn func(chat *model.Chat)) {
	wg := sync.WaitGroup{}
	for _, chat := range chats {
		wg.Add(1)
		chat := chat
		ctx := context.WithValue(context.Background(), ctxFuncnameKey, "foreachChat")
		j.pool.CtxGo(ctx, func(_ context.Context) {
			defer wg.Done()
			fn(chat)
		})
	}
	wg.Wait()
}

func (j *JobSet) activityJob() {
	ok, release := j.checkExclusive(&j.activityFlag)
	if !ok {
		return
	}
	defer release()
	chats, _ := dao.QueryChats()
	if len(chats) == 0 {
		return
	}

	// foreach chat
	j.foreachChat(chats, func(chat *model.Chat) {
		// get new events
		newEvents, _ := service.GetActivityEvents(chat.Username, chat.Token, 1)
		if len(newEvents) == 0 {
			return
		}

		// get old events and calc diff
		oldEvents, err := dao.GetActivities(chat.ChatID)
		if err != nil {
			return
		}
		logger.Logger().Infof("Get old ativities: #%d | (%d %s)", len(oldEvents), chat.ChatID, chat.Username)
		diff := model.DiffActivityEventSlice(newEvents, oldEvents)
		logger.Logger().Infof("Get diff ativities: #%d | (%d %s)", len(diff), chat.ChatID, chat.Username)
		if len(diff) == 0 {
			return
		}

		// update old events
		err = dao.SetActivities(chat.ChatID, newEvents)
		if err != nil {
			return
		}
		logger.Logger().Infof("Set new ativities: #%d | (%d %s)", len(newEvents), chat.ChatID, chat.Username)

		// format and send
		formatted := service.FormatActivityEvents(diff, chat.Username, -1) // <<< MarkdownV2
		if formatted == "" {
			return
		}
		dest, err := j.bw.Bot().ChatByID(xnumber.I64toa(chat.ChatID))
		if err == nil {
			opts := []interface{}{telebot.ModeMarkdownV2}
			if chat.Silent {
				opts = append(opts, telebot.Silent)
			}
			if !chat.Preview {
				opts = append(opts, telebot.NoPreview)
			}
			j.bw.RespondSend(dest, formatted, opts...)
		}
	})
}

func (j *JobSet) issueJob() {
	ok, release := j.checkExclusive(&j.issueFlag)
	if !ok {
		return
	}
	defer release()
	chats, _ := dao.QueryChats()
	if len(chats) == 0 {
		return
	}

	// foreach chat
	j.foreachChat(chats, func(chat *model.Chat) {
		// get new events and filter
		if chat.Token == "" || !chat.Issue {
			return
		}
		newEvents, _ := service.GetIssueEvents(chat.Username, chat.Token, 1)
		if chat.FilterMe {
			newEvents = service.FilterIssueEventSlice(newEvents, chat.Username)
		}
		if len(newEvents) == 0 {
			return
		}

		// get old events and calc diff
		oldEvents, err := dao.GetIssues(chat.ChatID)
		if err != nil {
			return
		}
		logger.Logger().Infof("Get old issues: #%d | (%d %s)", len(oldEvents), chat.ChatID, chat.Username)
		diff := model.DiffIssueEventSlice(newEvents, oldEvents)
		logger.Logger().Infof("Get diff issues: #%d | (%d %s)", len(diff), chat.ChatID, chat.Username)
		if len(diff) == 0 {
			return
		}

		// update old events
		err = dao.SetIssues(chat.ChatID, newEvents)
		if err != nil {
			return
		}
		logger.Logger().Infof("Set new issues: #%d | (%d %s)", len(newEvents), chat.ChatID, chat.Username)

		// format
		formatted := service.FormatIssueEvents(diff, chat.Username, -1) // <<< MarkdownV2
		if formatted == "" {
			return
		}
		dest, err := j.bw.Bot().ChatByID(xnumber.I64toa(chat.ChatID))
		if err == nil {
			opts := []interface{}{telebot.ModeMarkdownV2}
			if chat.Silent {
				opts = append(opts, telebot.Silent)
			}
			if !chat.Preview {
				opts = append(opts, telebot.NoPreview)
			}
			j.bw.RespondSend(dest, formatted, opts...)
		}
	})
}
