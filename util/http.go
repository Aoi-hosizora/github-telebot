package util

import (
	"fmt"
	"github.com/Aoi-hosizora/ah-tgbot/model"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	GithubUserApi                  string = "https://api.github.com/users/%s"
	GithubReceivedActivityEventApi string = "https://api.github.com/users/%s/received_events"
	// GithubReceivedIssueEventApi    string = "http://206.189.236.169:10014/gh/users/%s/issues/timeline"
	GithubReceivedIssueEventApi string = "http://localhost:10014/gh/users/%s/issues/timeline"
)

func CheckGithubUser(username string, private bool, token string) (bool, error) {
	url := fmt.Sprintf(GithubUserApi, username)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	if private {
		req.Header.Add("Authorization", fmt.Sprintf("Token %s", token))
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return resp.StatusCode == 200, nil
}

func GetGithubActivityEvents(username string, private bool, token string, page int) (response string, err error) {
	url := fmt.Sprintf(GithubReceivedActivityEventApi, username)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?page=%d", url, page), nil)
	if err != nil {
		return "", err
	}
	if private {
		req.Header.Add("Authorization", fmt.Sprintf("Token %s", token))
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func RenderGithubActivityString(objs []*model.ActivityEvent) string {
	result := ""
	if len(objs) == 1 {
		return renderGithubActivityString(objs[0])
	}
	for idx, obj := range objs {
		result += fmt.Sprintf("%d. %s\n", idx+1, renderGithubActivityString(obj))
	}
	return result[:len(result)-1]
}

func GetGithubIssueEvents(username string, private bool, token string, page int) (response string, err error) {
	url := fmt.Sprintf(GithubReceivedIssueEventApi, username)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?page=%d", url, page), nil)
	if err != nil {
		return "", err
	}
	if private {
		req.Header.Add("Authorization", fmt.Sprintf("Token %s", token))
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func RenderGithubIssueString(objs []*model.IssueEvent) string {
	result := ""
	if len(objs) == 1 {
		return renderGithubIssueString(objs[0])
	}
	for idx, obj := range objs {
		result += fmt.Sprintf("%d. %s\n", idx+1, renderGithubIssueString(obj))
	}
	return result[:len(result)-1]
}
