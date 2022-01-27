package service

import (
	"errors"
	"io/ioutil"
	"net/http"
)

func httpGet(url string, fn func(r *http.Request)) ([]byte, *http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	if fn != nil {
		fn(req)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != 200 {
		return nil, nil, errors.New("service: get non-200 response")
	}

	body := resp.Body
	defer body.Close()
	bs, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, nil, err
	}

	return bs, resp, nil
}
