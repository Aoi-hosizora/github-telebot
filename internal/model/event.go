package model

import (
	"github.com/Aoi-hosizora/ahlib/xgeneric/xgslice"
	"time"
)

type Repo struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type ActivityEvent struct {
	Id    string `json:"id"`
	Type  string `json:"type"`
	Actor *struct {
		Login string `json:"login"`
		Url   string `json:"url"`
	} `json:"actor"`
	Repo      *Repo            `json:"repo"`
	Payload   *ActivityPayload `json:"payload"`
	Public    bool             `json:"public"`
	CreatedAt time.Time        `json:"created_at"`
}

type ActivityPayload struct {
	Size    uint32 `json:"size"` // 1
	Commits []*struct {
		Sha string `json:"sha"` // 073f349775f412746a0494426a3c66d877c8033d
	} `json:"commits"`
	Ref     string    `json:"ref"`      // refs/heads/master null 1.1
	RefType string    `json:"ref_type"` // branch repository tag
	Forkee  *struct { // fork
		FullName string `json:"full_name"`
		HtmlUrl  string `json:"html_url"`
	} `json:"forkee"`
	Action string `json:"action"` // started added opened closed created published
	Issue  *struct {
		HtmlUrl string `json:"html_url"`
		Number  uint32 `json:"number"` // 1
	} `json:"issue"`
	Comment *struct {
		Id       uint32 `json:"id"` // 34806156
		HtmlUrl  string `json:"html_url"`
		CommitId string `json:"commit_id"` // 30aabf2faa716b5e489fca003aa3f7e7c4af4b23
	} `json:"comment"`
	PullRequest *struct { // pullRequest
		HtmlUrl string `json:"html_url"`
		Number  uint32 `json:"number"` // 1
	} `json:"pull_request"`
	Number uint32    `json:"number"` // 1
	Member *struct { // member
		HtmlUrl string `json:"html_url"`
		Login   string `json:"login"`
	} `json:"member"`
	Release *struct { // release
		HtmlUrl string `json:"html_url"`
		TagName string `json:"tag_name"` // 1.1
	} `json:"release"`
	Page []interface{} `json:"page"` // gollum
}

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

func DiffActivityEventSlice(s1 []*ActivityEvent, s2 []*ActivityEvent) []*ActivityEvent {
	return xgslice.DiffWith(s1, s2, func(e1, e2 *ActivityEvent) bool {
		// checking type and repo is dummy
		return e1.Id == e2.Id && e1.Type == e2.Type && e1.Repo.Name == e2.Repo.Name
	})
}

func DiffIssueEventSlice(s1 []*IssueEvent, s2 []*IssueEvent) []*IssueEvent {
	return xgslice.DiffWith(s1, s2, func(e1, e2 *IssueEvent) bool {
		// `id` is null when `event` is `opened` or `cross-referenced`
		return e1.Id == e2.Id && e1.Event == e2.Event && e1.Repo == e2.Repo && e1.Number == e2.Number && e1.CreatedAt.Unix() == e2.CreatedAt.Unix()
	})
}
