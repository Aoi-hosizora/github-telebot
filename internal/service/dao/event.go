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

func concatActivityPattern(chatID, id, typ, repo string) string {
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

func concatIssuePattern(chatID, id, event, repo, num, ct string) string {
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

func GetActivities(chatID int64) ([]*model.ActivityEvent, error) {
	pattern := concatActivityPattern(xnumber.I64toa(chatID), "*", "*", "*")
	keys, err := database.RedisClient().Keys(context.Background(), pattern).Result()
	if err != nil {
		return nil, err
	}

	events := make([]*model.ActivityEvent, 0, len(keys))
	for _, key := range keys {
		_, id, typ, repo := parseActivityPattern(key)
		events = append(events, &model.ActivityEvent{Id: id, Type: typ, Repo: &model.Repo{Name: repo}})
	}
	return events, nil
}

func GetIssues(chatID int64) ([]*model.IssueEvent, error) {
	pattern := concatIssuePattern(xnumber.I64toa(chatID), "*", "*", "*", "*", "*")
	keys, err := database.RedisClient().Keys(context.Background(), pattern).Result()
	if err != nil {
		return nil, err
	}

	events := make([]*model.IssueEvent, 0, len(keys))
	for _, key := range keys {
		_, id, event, repo, num, ct := parseIssuePattern(key)
		events = append(events, &model.IssueEvent{Id: id, Event: event, Repo: repo, Number: num, CreatedAt: ct})
	}
	return events, nil
}

func SetActivities(chatID int64, events []*model.ActivityEvent) error {
	pattern := concatActivityPattern(xnumber.I64toa(chatID), "*", "*", "*")
	_, err := xredis.DelAll(context.Background(), database.RedisClient(), pattern)
	if err != nil {
		return err
	}

	kvs := make([]interface{}, 0, len(events)*2)
	for _, ev := range events {
		idStr := xnumber.I64toa(chatID)
		pattern = concatActivityPattern(idStr, ev.Id, ev.Type, ev.Repo.Name)
		kvs = append(kvs, pattern, idStr)
	}
	err = database.RedisClient().MSet(context.Background(), kvs...).Err()
	return err
}

func SetIssues(chatID int64, events []*model.IssueEvent) error {
	pattern := concatIssuePattern(xnumber.I64toa(chatID), "*", "*", "*", "*", "*")
	_, err := xredis.DelAll(context.Background(), database.RedisClient(), pattern)
	if err != nil {
		return err
	}

	kvs := make([]interface{}, 0, len(events)*2)
	for _, ev := range events {
		idStr := xnumber.I64toa(chatID)
		pattern = concatIssuePattern(idStr, xnumber.I64toa(ev.Id), ev.Event, ev.Repo, xnumber.I32toa(ev.Number), xnumber.I64toa(ev.CreatedAt.Unix()))
		kvs = append(kvs, pattern, idStr)
	}
	err = database.RedisClient().MSet(context.Background(), kvs...).Err()
	return err
}
