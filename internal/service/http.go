package service

import (
	"io/ioutil"
	"net/http"
)

const (
	UserApi          string = "https://api.github.com/users/%s"
	ActivityEventApi string = "https://api.github.com/users/%s/received_events"
	IssueEventApi    string = "http://api.common.aoihosizora.top/github/users/%s/issues/timeline"
)

func HttpGet(url string, fn func(r *http.Request)) (int, []byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, nil, err
	}
	if fn != nil {
		fn(req)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}

	body := resp.Body
	defer body.Close()
	bs, err := ioutil.ReadAll(body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, bs, nil
}
