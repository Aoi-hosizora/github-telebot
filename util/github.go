package util

import (
	"fmt"
	"github.com/Aoi-hosizora/ah-tgbot/model"
	"github.com/Aoi-hosizora/ahlib/xcondition"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	GithubUserApi          string = "https://api.github.com/users/%s"
	GithubReceivedEventApi string = "https://api.github.com/users/%s/received_events"
)

func CheckGithubUser(username string, private bool, token string) (bool, error) {
	url := fmt.Sprintf(GithubUserApi, username)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	if private {
		req.Header.Add("Authorization", fmt.Sprintf("Token %s", token))
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return resp.StatusCode == 200, nil
}

func GetGithubEvents(username string, private bool, token string, page int) (response string, err error) {
	url := fmt.Sprintf(GithubReceivedEventApi, username)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?page=%d", url, page), nil)
	if err != nil {
		return "", err
	}
	if private {
		req.Header.Add("Authorization", fmt.Sprintf("Token %s", token))
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

func RenderGithubActions(objs []*model.GithubEvent) string {
	result := ""
	if len(objs) == 1 {
		return renderGithubAction(objs[0])
	}
	for idx, obj := range objs {
		result += fmt.Sprintf("%d. %s\n", idx+1, renderGithubAction(obj))
	}
	return result[:len(result)-1]
}

func renderGithubAction(obj *model.GithubEvent) string {
	userUrl := fmt.Sprintf("https://github.com/%s", obj.Actor.Login)
	repoUrl := fmt.Sprintf("https://github.com/%s", obj.Repo.Name)
	userMd := fmt.Sprintf("[%s](%s)", obj.Actor.Login, userUrl)
	repoMd := fmt.Sprintf("[%s](%s)", obj.Repo.Name, repoUrl)
	pl := obj.Payload

	message := ""
	switch obj.Type {
	case "PushEvent":
		cnt := xcondition.IfThenElse(pl.Size <= 1, "1 commit", fmt.Sprintf("%d commits", pl.Size)).(string)
		commitUrl := fmt.Sprintf("%s/commits/%s", repoUrl, pl.Commits[0].Sha)
		detail := fmt.Sprintf("[%s](%s)", pl.Commits[0].Sha[0:7], commitUrl)
		detail = xcondition.IfThenElse(pl.Size <= 1, detail, detail+"...").(string)
		message = fmt.Sprintf("%s pushed %s (%s) to %s", userMd, cnt, detail, repoMd)
	case "WatchEvent":
		message = fmt.Sprintf("%s starred %s", userMd, repoMd)
	case "CreateEvent":
		switch pl.RefType {
		case "branch":
			branchUrl := fmt.Sprintf("%s/tree/%s", repoUrl, pl.Ref)
			message = fmt.Sprintf("%s created branch [%s](%s) at %s", userMd, pl.Ref, branchUrl, repoMd)
		case "tag":
			tagUrl := fmt.Sprintf("%s/tree/%s", repoUrl, pl.Ref)
			message = fmt.Sprintf("%s created tag [%s](%s) at %s", userMd, pl.Ref, tagUrl, repoMd)
		case "repository":
			message = fmt.Sprintf("%s created %s repository %s", userMd, xcondition.IfThenElse(obj.Public, "public", "private").(string), repoMd)
		default:
			message = fmt.Sprintf("%s created %s at %s", userMd, pl.RefType, repoMd)
		}
	case "ForkEvent":
		message = fmt.Sprintf("%s forked %s to [%s](%s)", userMd, repoMd, pl.Forkee.FullName, pl.Forkee.HtmlUrl)
	case "DeleteEvent":
		message = fmt.Sprintf("%s delete %s %s at %s", userMd, pl.RefType, pl.Ref, repoMd)
	case "PublicEvent":
		message = fmt.Sprintf("%s made %s %s", userMd, repoMd, xcondition.IfThenElse(obj.Public, "public", "private").(string))

	case "IssuesEvent":
		message = fmt.Sprintf("%s %s issue [#%d](%s) in %s",
			userMd, pl.Action, pl.Issue.Number, pl.Issue.HtmlUrl, repoMd)
	case "IssueCommentEvent":
		message = fmt.Sprintf("%s %s comment [%d](%s) on issue [#%d](%s) in %s",
			userMd, pl.Action, pl.Comment.Id, pl.Comment.HtmlUrl, pl.Issue.Number, pl.Issue.HtmlUrl, repoMd)

	case "PullRequestEvent":
		message = fmt.Sprintf("%s %s pull request [#%d](%s) at %s",
			userMd, pl.Action, pl.Number, pl.PullRequest.HtmlUrl, repoMd)
	case "PullRequestReviewCommentEvent":
		message = fmt.Sprintf("%s %s pull request review comment [%d](%s) in pull request [#%d](%s) at %s",
			userMd, pl.Action, pl.Comment.Id, pl.Comment.HtmlUrl, pl.PullRequest.Number, pl.PullRequest.HtmlUrl, repoMd)

	case "CommitCommentEvent":
		message = fmt.Sprintf("%s %s comment [%d](%s) at commit %s in %s",
			userMd, pl.Action, pl.Comment.Id, pl.Comment.HtmlUrl, pl.Comment.CommitId[0:7], repoMd)

	case "MemberEvent":
		message = fmt.Sprintf("%s %s member [%s](%s) to %s", userMd, pl.Action, pl.Member.Login, pl.Member.HtmlUrl, repoMd)
	case "ReleaseEvent":
		message = fmt.Sprintf("%s release [%s](%s) at %s", userMd, pl.Release.TagName, pl.Release.HtmlUrl, repoMd)
	case "GollumEvent":
		cnt := xcondition.IfThenElse(len(pl.Page) <= 1, "1 wiki page", fmt.Sprintf("%d wiki pages", len(pl.Page))).(string)
		message = fmt.Sprintf("%s updated %s at %s", userMd, cnt, repoMd)
	default:
		message = fmt.Sprintf("%s: %s %s", strings.TrimRight(obj.Type, "Event"), userMd, repoMd)
	}

	if !obj.Public {
		message += " (private)"
	}
	return message
}
