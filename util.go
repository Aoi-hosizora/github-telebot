package main

import (
	"fmt"
	"time"
)

type githubUtil struct{}

var GithubUtil githubUtil

type GithubEvent struct {
	Type  string `json:"type"`
	Actor *struct {
		Login        string `json:"login"`
		DisplayLogin string `json:"display_login"`
		Url          string `json:"url"`
	} `json:"actor"`
	Repo *struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"repo"`
	Public    bool      `json:"public"`
	CreatedAt time.Time `json:"created_at"`
}

func (g githubUtil) WrapGithubAction(obj *GithubEvent) string {
	return fmt.Sprintf("%s: %s %s", obj.Type, obj.Actor.Login, obj.Repo.Name)
}

func (g githubUtil) WrapGithubActions(objs []*GithubEvent) string {
	result := ""
	for _, obj := range objs {
		result += g.WrapGithubAction(obj) + "\n"
	}
	return result
}
