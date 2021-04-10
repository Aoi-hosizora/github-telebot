package service

import (
	"fmt"
	"github.com/Aoi-hosizora/github-telebot/internal/model"
	"strings"
)

const (
	UserApi          string = "https://api.github.com/users/%s"
	ActivityEventApi string = "https://api.github.com/users/%s/received_events?page=%d"
	IssueEventApi    string = "http://api.common.aoihosizora.top/github/users/%s/issues/timeline?page=%d"
)

func CheckUserExist(username string, token string) (bool, error) {
	url := fmt.Sprintf(UserApi, username)
	_, resp, err := httpGet(url, githubToken(token))
	if err != nil {
		return false, err
	}
	return resp.StatusCode == 200, nil
}

func GetActivityEvents(username string, token string, page int) ([]byte, error) {
	url := fmt.Sprintf(ActivityEventApi, username, page)
	bs, _, err := httpGet(url, githubToken(token))
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func GetIssueEvents(username string, token string, page int) ([]byte, error) {
	url := fmt.Sprintf(IssueEventApi, username, page)
	bs, _, err := httpGet(url, githubToken(token))
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func RenderActivityEvents(events []*model.ActivityEvent) string {
	if len(events) == 0 {
		return ""
	}
	if len(events) == 1 {
		return RenderActivity(events[0]) // <<<
	}

	sb := strings.Builder{}
	for idx, obj := range events {
		if r := RenderActivity(obj); r != "" { // <<<
			sb.WriteString(fmt.Sprintf("%d\\. %s\n", idx+1, r)) // <<<
		}
	}
	if sb.Len() == 0 {
		return ""
	}
	return sb.String()[:sb.Len()-1]
}

func RenderIssueEvents(events []*model.IssueEvent) string {
	if len(events) == 0 {
		return ""
	}
	if len(events) == 1 {
		return RenderIssue(events[0]) // <<<
	}

	sb := strings.Builder{}
	for idx, obj := range events {
		if r := RenderIssue(obj); r != "" { // <<<
			sb.WriteString(fmt.Sprintf("%d\\. %s\n", idx+1, r)) // <<<
		}
	}
	if sb.Len() == 0 {
		return ""
	}
	return sb.String()[:sb.Len()-1]
}

func ConcatListAndUsername(list, username string) string {
	res := fmt.Sprintf("From [%s](https://github.com/%s)\\.", Markdown(username), username)
	return fmt.Sprintf("%s\n%s\n%s", list, `\=\=\=\=`, res)
}
