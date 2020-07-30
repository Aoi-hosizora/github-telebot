package service

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xcondition"
	"github.com/Aoi-hosizora/github-telebot/src/model"
	"strings"
)

func RenderResult(list, username string) string {
	username = fmt.Sprintf("From [%s](https://github.com/%s).", username, username)
	return fmt.Sprintf("%s\n====\n%s", list, username)
}

func Markdown(s string) string {
	s = strings.Split(s, "\\", "\\\\")
	s = strings.Split(s, "[", "\\[")
	s = strings.Split(s, "]", "\\]")
	s = strings.Split(s, "*", "\\*")
	return s
}

func RenderActivity(obj *model.ActivityEvent) string {

	userUrl := fmt.Sprintf("https://github.com/%s", obj.Actor.Login)
	repoUrl := fmt.Sprintf("https://github.com/%s", obj.Repo.Name)
	userMd := fmt.Sprintf("[%s](%s)", Markdown(obj.Actor.Login), userUrl)
	repoMd := fmt.Sprintf("[%s](%s)", Markdown(obj.Repo.Name), repoUrl)
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
		message = fmt.Sprintf("%s %s member [%s](%s) to %s", userMd, pl.Action, Markdown(pl.Member.Login), pl.Member.HtmlUrl, repoMd)
	case "ReleaseEvent":
		message = fmt.Sprintf("%s release [%s](%s) at %s", userMd, Markdown(pl.Release.TagName), pl.Release.HtmlUrl, repoMd)
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

func RenderIssue(obj *model.IssueEvent) string {
	userUrl := fmt.Sprintf("https://github.com/%s", obj.Actor.Login)
	repoUrl := fmt.Sprintf("https://github.com/%s", obj.Repo)
	issueUrl := fmt.Sprintf("https://github.com/%s/issues/%d", obj.Repo, obj.Number)
	userMd := fmt.Sprintf("[%s](%s)", Markdown(obj.Actor.Login), userUrl)
	issueMd := fmt.Sprintf("[%s#%d](%s)", Markdown(obj.Repo), obj.Number, issueUrl)

	message := ""
	switch obj.Event {
	case "mentioned":
		message = fmt.Sprintf("%s is mentioned on %s", userMd, issueMd)
	case "opened":
		message = fmt.Sprintf("%s opened %s", userMd, issueMd)
	case "closed":
		message = fmt.Sprintf("%s closed %s", userMd, issueMd)
	case "reopened":
		message = fmt.Sprintf("%s reopened %s", userMd, issueMd)
	case "renamed":
		message = fmt.Sprintf("%s renamed %s to `%s`", userMd, issueMd, obj.Rename.To)

	case "labeled":
		labelUrl := fmt.Sprintf("%s/labels/%s", repoUrl, obj.Label.Name)
		message = fmt.Sprintf("%s added label [%s](%s) to %s", userMd, Markdown(obj.Label.Name), labelUrl, issueMd)
	case "unlabeled":
		labelUrl := fmt.Sprintf("%s/labels/%s", repoUrl, obj.Label.Name)
		message = fmt.Sprintf("%s removed label [%s](%s) from %s", userMd, Markdown(obj.Label.Name), labelUrl, issueMd)

	case "locked":
		message = fmt.Sprintf("%s locked %s", userMd, issueMd)
	case "unlocked":
		message = fmt.Sprintf("%s unlocked %s", userMd, issueMd)

	case "milestoned":
		milestoneUrl := fmt.Sprintf("%s/milestones", repoUrl)
		message = fmt.Sprintf("%s added %s to milestone [%s](%s)", userMd, issueMd, Markdown(obj.Milestone.Title), milestoneUrl)
	case "demilestoned":
		milestoneUrl := fmt.Sprintf("%s/milestones", repoUrl)
		message = fmt.Sprintf("%s removed %s from milestone [%s](%s)", userMd, issueMd, Markdown(obj.Milestone.Title), milestoneUrl)

	case "pinned":
		message = fmt.Sprintf("%s pinned %s", userMd, issueMd)
	case "unpinned":
		message = fmt.Sprintf("%s unpinned %s", userMd, issueMd)

	case "assigned":
		message = fmt.Sprintf("%s is assigned to %s", userMd, issueMd)
	case "commented":
		message = fmt.Sprintf("%s added a [comment](%s) to %s", userMd, obj.HtmlUrl, issueMd)

	case "merged":
		prMd := fmt.Sprintf("[%s#%d](https://github.com/%s/pull/%d)", Markdown(obj.Repo), obj.Number, obj.Repo, obj.Number)
		message = fmt.Sprintf("%s merged pull request %s", userMd, prMd)
	case "cross-referenced":
		mdShow := fmt.Sprintf("%s/%s#%d", obj.Source.Issue.Repository.Owner.Login, obj.Source.Issue.Repository.Name, obj.Source.Issue.Number)
		mdUrl := fmt.Sprintf("https://github.com/%s/%s/issue/%d", obj.Source.Issue.Repository.Owner.Login, obj.Source.Issue.Repository.Name, obj.Source.Issue.Number)
		targetMd := fmt.Sprintf("[%s](%s)", Markdown(mdShow), mdUrl)
		message = fmt.Sprintf("%s mentioned %s from %s", userMd, issueMd, targetMd)
	case "referenced":
		toRepo := obj.CommitUrl // https://api.github.com/repos/gofiber/fiber/commits/a65d5027f336339cf4fe20cda0232c56cd64212e
		if toRepo == "" {
			message = fmt.Sprintf("%s added a commit that referenced %s", userMd, issueMd)
		} else if toSp := strings.Split(toRepo, "/"); len(toSp) >= 4 {
			commitMd := fmt.Sprintf("[%s](%s)", obj.CommitId[:7], obj.CommitUrl)
			toRepo = fmt.Sprintf("%s/%s", toSp[len(toSp)-4], toSp[len(toSp)-3])
			toRepoUrl := fmt.Sprintf("https://github.com/%s", toRepo)
			targetMd := fmt.Sprintf("[%s](%s)", Markdown(toRepo), toRepoUrl)
			message = fmt.Sprintf("%s added a commit %s of %s that referenced %s", userMd, commitMd, targetMd, issueMd)
		} else {
			message = ""
		}
	default:
		message = fmt.Sprintf("%s: %s %s", strings.Title(obj.Event), userMd, issueMd)
	}

	return message
}
