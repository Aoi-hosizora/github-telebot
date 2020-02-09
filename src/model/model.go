package model

import (
	"time"
)

type GithubEvent struct {
	Type  string `json:"type"`
	Actor *struct {
		Login string `json:"login"`
		Url   string `json:"url"`
	} `json:"actor"`
	Repo *struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"repo"`
	Payload   *Payload  `json:"payload"`
	Public    bool      `json:"public"`
	CreatedAt time.Time `json:"created_at"`
}

type Payload struct {
	Size    uint32 `json:"size"` // 1
	Commits []*struct {
		Sha string `json:"sha"` // 073f349775f412746a0494426a3c66d877c8033d
	}
	Ref     string `json:"ref"`      // refs/heads/master null 1.1
	RefType string `json:"ref_type"` // branch repository tag
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
	Number uint32 `json:"number"` // 1
	Member *struct { // member
		HtmlUrl string `json:"html_url"`
		Login   string `json:"login"`
	}
	Release *struct { // release
		HtmlUrl string `json:"html_url"`
		TagName string `json:"tag_name"` // 1.1
	}
	Page []interface{} // gollum
}
