package service

import (
	"fmt"
	"github.com/Aoi-hosizora/github-telebot/internal/model"
	"net/http"
)

func CheckUser(username string, token string) (bool, error) {
	url := fmt.Sprintf(UserApi, username)
	code, _, err := HttpGet(url, func(r *http.Request) {
		r.Header.Add("Authorization", fmt.Sprintf("Token %s", token))
	})
	return code == 200, err
}

func GetActivityEvents(username string, token string, page int) ([]byte, error) {
	url := fmt.Sprintf(ActivityEventApi, username)
	url = fmt.Sprintf("%s?page=%d", url, page)
	_, bs, err := HttpGet(url, func(r *http.Request) {
		r.Header.Add("Authorization", fmt.Sprintf("Token %s", token))
	})
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func GetIssueEvents(username string, token string, page int) ([]byte, error) {
	url := fmt.Sprintf(IssueEventApi, username)
	url = fmt.Sprintf("%s?page=%d", url, page)
	_, bs, err := HttpGet(url, func(r *http.Request) {
		r.Header.Add("Authorization", fmt.Sprintf("Token %s", token))
	})
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func RenderActivities(objs []*model.ActivityEvent) string {
	if len(objs) == 1 {
		return RenderActivity(objs[0]) // <<<
	}

	result := ""
	for idx, obj := range objs {
		if r := RenderActivity(obj); r != "" { // <<<
			result += fmt.Sprintf("%d\\. %s\n", idx+1, r) // <<<
		}
	}

	if result == "" {
		return ""
	}
	return result[:len(result)-1]
}

func RenderIssues(objs []*model.IssueEvent) string {
	if len(objs) == 1 {
		return RenderIssue(objs[0]) // <<<
	}

	result := ""
	for idx, obj := range objs {
		if r := RenderIssue(obj); r != "" { // <<<
			result += fmt.Sprintf("%d\\. %s\n", idx+1, r) // <<<
		}
	}

	if result == "" {
		return ""
	}
	return result[:len(result)-1]
}
