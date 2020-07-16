package model

import (
	"encoding/json"
	"time"
)

type IssueEvent struct {
	Id    int64  `json:"id"`
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
	CommitId  string `json:"commit_id"`
	CommitUrl string `json:"commit_url"`
	Body      string `json:"body"`
	HtmlUrl   string `json:"html_url"`
	Source    *struct {
		Issue *struct {
			Number     int32  `json:"number"`
			HtmlUrl    string `json:"html_url"`
			Body       string `json:"body"`
			Repository *struct {
				Name  string `json:"name"`
				Owner *struct {
					Login string `json:"login"`
				} `json:"owner"`
			} `json:"repository"`
		} `json:"issue"`
	} `json:"source"`
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

func IssueEventEqual(e1, e2 *IssueEvent) bool {
	// `id` is null when `event` is `opened`
	return e1.Id == e2.Id && e1.Event == e2.Event && e1.Repo == e2.Repo
}
