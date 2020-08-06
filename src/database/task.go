package database

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xredis"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/github-telebot/src/model"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"strings"
	"time"
)

const MagicToken = "$$"

func getActivityPattern(chatId string, id, t, repo string) string {
	repo = strings.ReplaceAll(repo, "-", MagicToken)
	return fmt.Sprintf("gh-activity-ev-%s-%s-%s-%s", chatId, id, t, repo)
}

func parseActivityPattern(key string) (chatId int64, id, t, repo string) {
	sp := strings.Split(key, "-")
	chatId, _ = xnumber.ParseInt64(sp[3], 10)
	id = sp[4]
	t = sp[5]
	repo = strings.ReplaceAll(sp[6], MagicToken, "-")
	return
}

func getIssuePattern(chatId string, id, event, repo, ct, num string) string {
	event = strings.ReplaceAll(event, "-", MagicToken)
	repo = strings.ReplaceAll(repo, "-", MagicToken)
	return fmt.Sprintf("gh-issue-ev-%s-%s-%s-%s-%s-%s", chatId, id, event, repo, ct, num)
}

func parseIssuePattern(key string) (chatId, id int64, event, repo string, ct time.Time, num int32) {
	sp := strings.Split(key, "-")
	chatId, _ = xnumber.ParseInt64(sp[3], 10)
	id, _ = xnumber.ParseInt64(sp[4], 10)
	event = strings.ReplaceAll(sp[5], MagicToken, "-")
	repo = strings.ReplaceAll(sp[6], MagicToken, "-")
	ctn, _ := xnumber.ParseInt64(sp[7], 10)
	ct = time.Unix(ctn, 0)
	num, _ = xnumber.ParseInt32(sp[8], 10)
	return
}

func GetOldActivities(chatId int64) ([]*model.ActivityEvent, bool) {
	pattern := getActivityPattern(strconv.FormatInt(chatId, 10), "*", "*", "*")
	redisMu.Lock()
	keys, err := redis.Strings(Conn.Do("KEYS", pattern))
	redisMu.Unlock()
	if err != nil {
		return nil, false
	}

	evs := make([]*model.ActivityEvent, len(keys))
	for idx := range evs {
		_, id, t, repo := parseActivityPattern(keys[idx])
		evs[idx] = &model.ActivityEvent{
			Id: id, Type: t, Repo: &struct {
				Name string `json:"name"`
				Url  string `json:"url"`
			}{Name: repo},
		}
	}
	return evs, true
}

func SetOldActivities(chatId int64, evs []*model.ActivityEvent) bool {
	pattern := getActivityPattern(strconv.FormatInt(chatId, 10), "*", "*", "*")
	redisMu.Lock()
	tot, del, err := xredis.WithConn(Conn).DeleteAll(pattern)
	redisMu.Unlock()
	if err != nil || (tot != 0 && del == 0) {
		return false
	}

	keys := make([]string, 0)
	values := make([]string, 0)
	for _, ev := range evs {
		id := strconv.FormatInt(chatId, 10)
		pattern := getActivityPattern(id, ev.Id, ev.Type, ev.Repo.Name)
		keys = append(keys, pattern)
		values = append(values, id)
	}

	redisMu.Lock()
	tot, add, err := xredis.WithConn(Conn).SetAll(keys, values)
	redisMu.Unlock()
	return err == nil && (tot == 0 || add >= 1)
}

func GetOldIssues(chatId int64) ([]*model.IssueEvent, bool) {
	redisMu.Lock()
	keys, err := redis.Strings(Conn.Do("KEYS", getIssuePattern(strconv.FormatInt(chatId, 10), "*", "*", "*", "*", "*")))
	redisMu.Unlock()
	if err != nil {
		return nil, false
	}

	evs := make([]*model.IssueEvent, len(keys))
	for idx := range evs {
		_, id, event, repo, ct, num := parseIssuePattern(keys[idx])
		evs[idx] = &model.IssueEvent{Id: id, Event: event, Repo: repo, CreatedAt: ct, Number: num}
	}
	return evs, true
}

func SetOldIssues(chatId int64, evs []*model.IssueEvent) bool {
	pattern := getIssuePattern(strconv.FormatInt(chatId, 10), "*", "*", "*", "*", "*")
	redisMu.Lock()
	tot, del, err := xredis.WithConn(Conn).DeleteAll(pattern)
	redisMu.Unlock()
	if err != nil || (tot != 0 && del == 0) {
		return false
	}

	keys := make([]string, 0)
	values := make([]string, 0)
	for _, ev := range evs {
		id := strconv.FormatInt(chatId, 10)
		pattern := getIssuePattern(id, strconv.FormatInt(ev.Id, 10), ev.Event, ev.Repo, strconv.FormatInt(ev.CreatedAt.Unix(), 10), strconv.Itoa(int(ev.Number)))
		keys = append(keys, pattern)
		values = append(values, id)
	}

	redisMu.Lock()
	tot, add, err := xredis.WithConn(Conn).SetAll(keys, values)
	redisMu.Unlock()
	return err == nil && (tot == 0 || add >= 1)
}
