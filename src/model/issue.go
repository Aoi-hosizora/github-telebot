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

func UnmarshalIssueEvents(bs []byte) ([]*IssueEvent, error) {
	out := make([]*IssueEvent, 0)
	err := json.Unmarshal(bs, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func IssueEventEqual(e1, e2 *IssueEvent) bool {
	// `id` is null when `event` is `opened` or `cross-referenced`
	return e1.Id == e2.Id && e1.Event == e2.Event && e1.Repo == e2.Repo && e1.Number == e2.Number && e1.CreatedAt.String() == e2.CreatedAt.String()
}

func IssueSliceDiff(s1 []*IssueEvent, s2 []*IssueEvent) []*IssueEvent {
	result := make([]*IssueEvent, 0)
	for _, item1 := range s1 {
		exist := false
		for _, item2 := range s2 {
			if IssueEventEqual(item1, item2) {
				exist = true
				break
			}
		}
		if !exist {
			result = append(result, item1)
		}
	}
	return result
}
