package model

import (
	"encoding/json"
	"time"
)

type IssueEvent struct {
	Event string `json:"event"`
	Actor *struct {
		Login string `json:"login"`
		Url   string `json:"url"`
	} `json:"actor"`
	Repo    string `json:"repo"`
	Number  int32  `json:"number"`
	Involve string `json:"involve"`
	Rename  *struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"rename"`
	Label *struct {
		Color string `json:"color"`
		Name  string `json:"name"`
	} `json:"label"`
	Milestone *struct {
		Title string `json:"title"`
	} `json:"milestone"`
	CommitId  string    `json:"commit_id"`
	CommitUrl string    `json:"commit_url"`
	CreatedAt time.Time `json:"created_at"`
}

func UnmarshalIssueEvents(response string) ([]*IssueEvent, error) {
	out := make([]*IssueEvent, 0)
	err := json.Unmarshal([]byte(response), &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
