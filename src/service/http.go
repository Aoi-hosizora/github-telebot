package service

import (
	"io/ioutil"
	"net/http"
)

const (
	UserApi          string = "https://api.github.com/users/%s"
	ActivityEventApi string = "https://api.github.com/users/%s/received_events"
	IssueEventApi    string = "http://206.189.236.169:10014/gh/users/%s/issues/timeline"
	// IssueEventApi    string = "http://localhost:10014/gh/users/%s/issues/timeline"
)

func HttpGet(url string, header *http.Header) (int, []byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header = *header

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
