package util

import (
	"fmt"
	"github.com/Aoi-hosizora/ah-tgbot/src/config"
	"github.com/Aoi-hosizora/ah-tgbot/src/model"
	"github.com/Aoi-hosizora/ahlib/xcondition"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	GithubReceivedEventApi string = "https://api.github.com/users/%s/received_events"
)

func GetActions(config *config.GithubConfig, page int) (string, error) {
	url := fmt.Sprintf(GithubReceivedEventApi, config.Username)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?page=%d", url, page), nil)
	if err != nil {
		return "", err
	}
	if config.Private {
		req.Header.Add("Authorization", fmt.Sprintf("Token %s", config.Token))
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

func WrapGithubAction(obj *model.GithubEvent) string {
	userUrl := fmt.Sprintf("https://github.com/%s", obj.Actor.Login)
	repoUrl := fmt.Sprintf("https://github.com/%s", obj.Repo.Name)
	userMd := fmt.Sprintf("[%s](%s)", obj.Actor.Login, userUrl)
	repoMd := fmt.Sprintf("[%s](%s)", obj.Repo.Name, repoUrl)
	pl := obj.Payload

	message := ""
	switch obj.Type {
	case "PushEvent":
		commitUrl := fmt.Sprintf("%s/commits/%s", repoUrl, pl.Commits[0].Sha)
		message = fmt.Sprintf("%s pushed %d commit ([%s](%s)...) to %s", userMd, pl.Size, pl.Commits[0].Sha[0:7], commitUrl, repoMd)
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
		message = fmt.Sprintf("%s %s issue [#%d](%s) in %s", userMd, pl.Action, pl.Issue.Number, pl.Issue.HtmlUrl, repoMd)
	case "IssueCommentEvent":
		message = fmt.Sprintf("%s %s comment [%d](%s) on issue [#%d](%s) in %s", userMd, pl.Action, pl.Comment.Id, pl.Comment.HtmlUrl, pl.Issue.Number, pl.Issue.HtmlUrl, repoMd)
	case "PullRequestEvent":
		message = fmt.Sprintf("%s %s pull request [#%d](%s) at %s", userMd, pl.Action, pl.Number, pl.PullRequest.HtmlUrl, repoMd)
	case "PullRequestReviewCommentEvent":
		message = fmt.Sprintf("%s %s pull request review comment [%d](%s) in pull request [#%d](%s) at %s",
			userMd, pl.Action, pl.Comment.Id, pl.Comment.HtmlUrl, pl.PullRequest.Number, pl.PullRequest.HtmlUrl, repoMd)
	case "CommitCommentEvent":
		message = fmt.Sprintf("%s %s comment [%d](%s) at commit %s in %s", userMd, pl.Action, pl.Comment.Id, pl.Comment.HtmlUrl, pl.Comment.CommitId[0:7], repoMd)
	case "MemberEvent":
		message = fmt.Sprintf("%s %s member [%s](%s) to %s", userMd, pl.Action, pl.Member.Login, pl.Member.HtmlUrl, repoMd)
	case "ReleaseEvent":
		message = fmt.Sprintf("%s release [%s](%s) at %s", userMd, pl.Release.TagName, pl.Release.HtmlUrl, repoMd)
	case "GollumEvent":
		message = fmt.Sprintf("%s updated %d wiki page at %s", userMd, len(pl.Page), repoMd)
	default:
		message = fmt.Sprintf("%s: %s %s", strings.TrimRight(obj.Type, "Event"), userMd, repoMd)
	}
	return message
}

func WrapGithubActions(objs []*model.GithubEvent) string {
	result := ""
	if len(objs) == 1 {
		return WrapGithubAction(objs[0])
	}
	for idx, obj := range objs {
		result += fmt.Sprintf("%d. %s\n", idx+1, WrapGithubAction(obj))
	}
	return result
}
