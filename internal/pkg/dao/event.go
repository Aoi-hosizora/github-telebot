package dao

import (
	"context"
	"fmt"
	"github.com/Aoi-hosizora/ahlib-db/xredis"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/github-telebot/internal/model"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/database"
	"strings"
	"time"
)

const magicToken = "$$"

func getActivityPattern(chatID, id, typ, repo string) string {
	repo = strings.ReplaceAll(repo, "-", magicToken)
	return fmt.Sprintf("gh-activity-ev-%s-%s-%s-%s", chatID, id, typ, repo)
	//                                         3  4  5  6
}

func parseActivityPattern(key string) (chatID int64, id, typ, repo string) {
	sp := strings.Split(key, "-")
	chatID, _ = xnumber.Atoi64(sp[3])
	id = sp[4]
	typ = sp[5]
	repo = strings.ReplaceAll(sp[6], magicToken, "-")
	return
}

func getIssuePattern(chatID, id, event, repo, num, ct string) string {
	event = strings.ReplaceAll(event, "-", magicToken)
	repo = strings.ReplaceAll(repo, "-", magicToken)
	return fmt.Sprintf("gh-issue-ev-%s-%s-%s-%s-%s-%s", chatID, id, event, repo, num, ct)
	//                                      3  4  5  6  7  8
}

func parseIssuePattern(key string) (chatID, id int64, event, repo string, num int32, ct time.Time) {
	sp := strings.Split(key, "-")
	chatID, _ = xnumber.Atoi64(sp[3])
	id, _ = xnumber.Atoi64(sp[4])
	event = strings.ReplaceAll(sp[5], magicToken, "-")
	repo = strings.ReplaceAll(sp[6], magicToken, "-")
	num, _ = xnumber.Atoi32(sp[7])
	ctn, _ := xnumber.Atoi64(sp[8])
	ct = time.Unix(ctn, 0)
	return
}

func GetOldActivities(chatID int64) ([]*model.ActivityEvent, bool) {
	pattern := getActivityPattern(xnumber.I64toa(chatID), "*", "*", "*")
	keys, err := database.Redis().Keys(context.Background(), pattern).Result()
	if err != nil {
		return nil, false
	}

	events := make([]*model.ActivityEvent, 0, len(keys))
	for _, key := range keys {
		_, id, typ, repo := parseActivityPattern(key)
		events = append(events, &model.ActivityEvent{Id: id, Type: typ, Repo: &model.Repo{Name: repo}})
	}
	return events, true
}

func GetOldIssues(chatID int64) ([]*model.IssueEvent, bool) {
	pattern := getIssuePattern(xnumber.I64toa(chatID), "*", "*", "*", "*", "*")
	keys, err := database.Redis().Keys(context.Background(), pattern).Result()
	if err != nil {
		return nil, false
	}

	events := make([]*model.IssueEvent, 0, len(keys))
	for _, key := range keys {
		_, id, event, repo, num, ct := parseIssuePattern(key)
		events = append(events, &model.IssueEvent{Id: id, Event: event, Repo: repo, Number: num, CreatedAt: ct})
	}
	return events, true
}

func SetOldActivities(chatID int64, events []*model.ActivityEvent) bool {
	chatIDStr := xnumber.I64toa(chatID)
	pattern := getActivityPattern(chatIDStr, "*", "*", "*")
	_, err := xredis.DelAll(database.Redis(), context.Background(), pattern)
	if err != nil {
		return false
	}

	keys := make([]string, 0)
	values := make([]string, 0)
	for _, ev := range events {
		pattern = getActivityPattern(chatIDStr, ev.Id, ev.Type, ev.Repo.Name)
		keys = append(keys, pattern)
		values = append(values, chatIDStr)
	}

	_, err = xredis.SetAll(database.Redis(), context.Background(), keys, values)
	return err == nil
}

func SetOldIssues(chatID int64, events []*model.IssueEvent) bool {
	chatIDStr := xnumber.I64toa(chatID)
	pattern := getIssuePattern(chatIDStr, "*", "*", "*", "*", "*")
	_, err := xredis.DelAll(database.Redis(), context.Background(), pattern)
	if err != nil {
		return false
	}

	keys := make([]string, 0)
	values := make([]string, 0)
	for _, ev := range events {
		pattern = getIssuePattern(chatIDStr, xnumber.I64toa(ev.Id), ev.Event, ev.Repo, xnumber.I32toa(ev.Number), xnumber.I64toa(ev.CreatedAt.Unix()))
		keys = append(keys, pattern)
		values = append(values, chatIDStr)
	}

	_, err = xredis.SetAll(database.Redis(), context.Background(), keys, values)
	return err == nil
}
