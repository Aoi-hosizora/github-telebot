package service

import (
	"encoding/json"
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xgin/headers"
	"github.com/Aoi-hosizora/github-telebot/internal/model"
	"net/http"
	"strings"
)

const (
	UserApi          string = "https://api.github.com/users/%s"
	ActivityEventApi string = "https://api.github.com/users/%s/received_events?page=%d"
	IssueEventApi    string = "http://api.common.aoihosizora.top/github/users/%s/issues/timeline?page=%d"
)

func githubToken(token string) func(*http.Request) {
	return func(r *http.Request) {
		if token != "" {
			r.Header.Add(headers.Authorization, "Token "+token)
		}
	}
}

func CheckUserExistence(username string, token string) (bool, error) {
	url := fmt.Sprintf(UserApi, username)
	_, resp, err := httpGet(url, githubToken(token))
	if err != nil {
		return false, err
	}
	return resp.StatusCode == 200, nil
}

func GetActivityEvents(username string, token string, page int) ([]*model.ActivityEvent, error) {
	url := fmt.Sprintf(ActivityEventApi, username, page)
	bs, _, err := httpGet(url, githubToken(token))
	if err != nil {
		return nil, err
	}

	out := make([]*model.ActivityEvent, 0)
	err = json.Unmarshal(bs, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func GetIssueEvents(username string, token string, page int) ([]*model.IssueEvent, error) {
	url := fmt.Sprintf(IssueEventApi, username, page)
	bs, _, err := httpGet(url, githubToken(token))
	if err != nil {
		return nil, err
	}

	out := make([]*model.IssueEvent, 0)
	err = json.Unmarshal(bs, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// func FilterActivityEventSlice(events []*model.ActivityEvent) []*model.ActivityEvent {
// 	out := make([]*model.ActivityEvent, 0, len(events))
// 	for _, e := range events {
// 		// ...
// 		out = append(out, e)
// 	}
// 	return out
// }

func FilterIssueEventSlice(events []*model.IssueEvent, username string) []*model.IssueEvent {
	out := make([]*model.IssueEvent, 0, len(events))
	for _, e := range events {
		if e.Actor.Login != username {
			out = append(out, e)
		}
	}
	return out
}

func FormatActivityEvents(events []*model.ActivityEvent, username string, page int) string {
	if len(events) == 0 {
		return ""
	}

	sb := strings.Builder{}
	if page <= 0 {
		sb.WriteString(fmt.Sprintf("*New activity events*"))
	} else {
		sb.WriteString(fmt.Sprintf("*Activity events from page %d*", page))
	}

	if len(events) == 1 {
		sb.WriteString(formatActivityEvent(events[0])) // <<<
	} else {
		for idx, ev := range events {
			sb.WriteString(fmt.Sprintf("%d\\. %s\n", idx+1, formatActivityEvent(ev))) // <<<
		}
	}

	sb.WriteString(`\=\=\=\=`)
	sb.WriteString(fmt.Sprintf("\nFrom [%s](https://github.com/%s)\\.", Markdown(username), username))
	return sb.String()
}

func FormatIssueEvents(events []*model.IssueEvent, username string, page int) string {
	if len(events) == 0 {
		return ""
	}

	sb := strings.Builder{}
	if page <= 0 {
		sb.WriteString(fmt.Sprintf("*New issue events*"))
	} else {
		sb.WriteString(fmt.Sprintf("*Issue events from page %d*", page))
	}

	if len(events) == 1 {
		sb.WriteString(formatIssueEvent(events[0])) // <<<
	} else {
		for idx, ev := range events {
			sb.WriteString(fmt.Sprintf("%d\\. %s\n", idx+1, formatIssueEvent(ev))) // <<<
		}
	}

	sb.WriteString(`\=\=\=\=`)
	sb.WriteString(fmt.Sprintf("\nFrom [%s](https://github.com/%s)\\.", Markdown(username), username))
	return sb.String()
}
