package util

import (
	"fmt"
	"github.com/Aoi-hosizora/ah-tgbot/src/config"
	"github.com/Aoi-hosizora/ah-tgbot/src/model"
	"io/ioutil"
	"net/http"
)

const (
	GithubReceivedEventApi string = "https://api.github.com/users/%s/received_events"
)

func GetActions(config *config.GithubConfig, page int) (string, error) {
	url := fmt.Sprintf(GithubReceivedEventApi, config.Username)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?page=%d", url, page), nil)
	if err != nil {
		return "", err
	}
	if config.Private {
		req.Header.Add("Authorization", fmt.Sprintf("Token %s", config.Token))
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

func WrapGithubActions(objs []*model.GithubEvent) string {
	result := ""
	if len(objs) == 1 {
		return WrapGithubAction(objs[0])
	}
	for idx, obj := range objs {
		result += fmt.Sprintf("%d. %s\n", idx+1, WrapGithubAction(obj))
	}
	return result
}
