package util

import (
	"fmt"
	"github.com/Aoi-hosizora/ah-tgbot/model"
	"github.com/Aoi-hosizora/ahlib/xcondition"
	"strings"
)

func renderGithubActivityString(obj *model.ActivityEvent) string {
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

func renderGithubIssueString(obj *model.IssueEvent) string {
	actorUrl := fmt.Sprintf("https://github.com/%s", obj.Actor.Login)
	repoUrl := fmt.Sprintf("https://github.com/%s", obj.Repo)
	issueUrl := fmt.Sprintf("https://github.com/%s/issues/%d", obj.Repo, obj.Number)
	actorMd := fmt.Sprintf("[%s](%s)", obj.Actor.Login, actorUrl)
	repoMd := fmt.Sprintf("[%s](%s)", obj.Repo, repoUrl)
	issueMd := fmt.Sprintf("[#%d](%s)", obj.Number, issueUrl)

	issueRepoMd := fmt.Sprintf("issue %s in %s", issueMd, repoMd)

	message := ""
	switch obj.Event {
	case "mentioned":
		message = fmt.Sprintf("%s is mentioned on %s", actorMd, issueRepoMd)
	case "opened":
		message = fmt.Sprintf("%s opened %s", actorMd, issueRepoMd)
	case "closed":
		message = fmt.Sprintf("%s closed %s", actorMd, issueRepoMd)
	case "reopened":
		message = fmt.Sprintf("%s reopened %s", actorMd, issueRepoMd)
	case "renamed":
		message = fmt.Sprintf("%s renamed %s to [%s](%s)", actorMd, issueRepoMd, obj.Rename.To, issueUrl)

	case "labeled":
		labelUrl := fmt.Sprintf("%s/labels/%s", repoUrl, obj.Label.Name)
		message = fmt.Sprintf("%s added label [%s](%s) to %s", actorMd, obj.Label.Name, labelUrl, issueRepoMd)
	case "unlabeled":
		labelUrl := fmt.Sprintf("%s/labels/%s", repoUrl, obj.Label.Name)
		message = fmt.Sprintf("%s removed label [%s](%s) from %s", actorMd, obj.Label.Name, labelUrl, issueRepoMd)

	case "locked":
		message = fmt.Sprintf("%s locked %s", actorMd, issueRepoMd)
	case "unlocked":
		message = fmt.Sprintf("%s unlocked %s", actorMd, issueRepoMd)

	case "milestoned":
		milestoneUrl := fmt.Sprintf("%s/milestones", repoUrl)
		message = fmt.Sprintf("%s added %s to milestone [%s](%s)", actorMd, issueRepoMd, obj.Milestone.Title, milestoneUrl)
	case "demilestoned":
		milestoneUrl := fmt.Sprintf("%s/milestones", repoUrl)
		message = fmt.Sprintf("%s removed %s from milestone [%s](%s)", actorMd, issueRepoMd, obj.Milestone.Title, milestoneUrl)

	case "pinned":
		message = fmt.Sprintf("%s pinned %s", actorMd, issueRepoMd)
	case "unpinned":
		message = fmt.Sprintf("%s unpinned %s", actorMd, issueRepoMd)

	case "assigned":
		message = fmt.Sprintf("%s is assigned to %s", actorMd, issueRepoMd)
	case "commented":
		message = fmt.Sprintf("%s added a [comment](%s) to %s", actorMd, obj.HtmlUrl, issueRepoMd)

	case "merged":
		prMd := fmt.Sprintf("[#%d](https://github.com/%s/pull/%d)", obj.Number, obj.Repo, obj.Number)
		message = fmt.Sprintf("%s merged pull request %s from %s", actorMd, prMd, issueRepoMd)
	case "cross-referenced":
		targetRepoMd := fmt.Sprintf("[%s/%s](https://github.com/%s/%s)", obj.Source.Issue.Repository.Owner.Login, obj.Source.Issue.Repository.Name, obj.Source.Issue.Repository.Owner.Login, obj.Source.Issue.Repository.Name)
		targetIssueRepoMd := fmt.Sprintf("issue [#%d](%s) in %s", obj.Source.Issue.Number, obj.Source.Issue.HtmlUrl, targetRepoMd)
		message = fmt.Sprintf("%s mentioned %s from %s", actorMd, targetIssueRepoMd, issueRepoMd)
	case "referenced":
		toRepo := obj.CommitUrl // https://api.github.com/repos/gofiber/fiber/commits/a65d5027f336339cf4fe20cda0232c56cd64212e
		if toRepo == "" { // TODO: CommitUrl may be null
			message = fmt.Sprintf("%s added a commit that referenced %s", actorMd, issueRepoMd)
		} else {
			toSp := strings.Split(toRepo, "/")
			toRepo = fmt.Sprintf("%s/%s", toSp[len(toSp)-4], toSp[len(toSp)-3])
			toRepoUrl := fmt.Sprintf("https://github.com/%s", toRepo)
			message = fmt.Sprintf("%s added a commit [%s](%s) to [%s](%s) that referenced %s", actorMd, obj.CommitId[:7], obj.CommitUrl, toRepo, toRepoUrl, issueRepoMd)
		}
	default:
		message = fmt.Sprintf("%s: %s -> %s", strings.Title(obj.Event), actorMd, issueMd)
	}

	return message
}
