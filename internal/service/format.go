package service

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xgeneric/xsugar"
	"github.com/Aoi-hosizora/github-telebot/internal/model"
	"strings"
)

func markdownV2(s string) string {
	// https://core.telegram.org/bots/api#markdownv2-style

	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `*`, `\*`)
	s = strings.ReplaceAll(s, `[`, `\[`)
	s = strings.ReplaceAll(s, "`", "\\`")

	s = strings.ReplaceAll(s, `_`, `\_`)
	s = strings.ReplaceAll(s, `]`, `\]`)
	s = strings.ReplaceAll(s, `~`, `\~`)
	s = strings.ReplaceAll(s, `>`, `\>`)
	s = strings.ReplaceAll(s, `#`, `\#`)
	s = strings.ReplaceAll(s, `+`, `\+`)
	s = strings.ReplaceAll(s, `-`, `\-`)
	s = strings.ReplaceAll(s, `=`, `\=`)
	s = strings.ReplaceAll(s, `|`, `\|`)
	s = strings.ReplaceAll(s, `{`, `\{`)
	s = strings.ReplaceAll(s, `}`, `\}`)
	s = strings.ReplaceAll(s, `.`, `\.`)
	s = strings.ReplaceAll(s, `!`, `\!`)

	// [google\-test\_test\+test\=test\[test\]\!\(test\)\`test](www.google.co.jp)
	// google-test_test+test=test[test]!(test)`test (http://www.google.co.jp/)
	return s
}

func formatActivityEvent(obj *model.ActivityEvent) string {
	userUrl := fmt.Sprintf("https://github.com/%s", obj.Actor.Login)
	repoUrl := fmt.Sprintf("https://github.com/%s", obj.Repo.Name)
	userMd := fmt.Sprintf("[%s](%s)", markdownV2(obj.Actor.Login), userUrl)
	repoMd := fmt.Sprintf("[%s](%s)", markdownV2(obj.Repo.Name), repoUrl)

	message := ""
	switch pl := obj.Payload; obj.Type {
	case "PushEvent":
		count := xsugar.IfThenElse(pl.Size <= 1, "1 commit", fmt.Sprintf("%d commits", pl.Size))
		commitUrl := fmt.Sprintf("%s/commits/%s", repoUrl, pl.Commits[0].Sha)
		commitMd := fmt.Sprintf("[%s](%s)", pl.Commits[0].Sha[0:7], commitUrl)
		commitMd += xsugar.IfThenElse(pl.Size <= 1, "", "\\.\\.\\.")
		message = fmt.Sprintf("%s pushed %s \\(%s\\) to %s", userMd, count, commitMd, repoMd)
	case "WatchEvent":
		message = fmt.Sprintf("%s starred %s", userMd, repoMd)
	case "CreateEvent":
		switch pl.RefType {
		case "branch":
			branchUrl := fmt.Sprintf("%s/tree/%s", repoUrl, pl.Ref)
			message = fmt.Sprintf("%s created branch [%s](%s) at %s", userMd, markdownV2(pl.Ref), branchUrl, repoMd)
		case "tag":
			tagUrl := fmt.Sprintf("%s/tree/%s", repoUrl, pl.Ref)
			message = fmt.Sprintf("%s created tag [%s](%s) at %s", userMd, markdownV2(pl.Ref), tagUrl, repoMd)
		case "repository":
			pub := xsugar.IfThenElse(obj.Public, "public", "private")
			message = fmt.Sprintf("%s created %s repository %s", userMd, pub, repoMd)
		default:
			message = fmt.Sprintf("%s created %s at %s", userMd, markdownV2(pl.RefType), repoMd)
		}
	case "ForkEvent":
		message = fmt.Sprintf("%s forked %s to [%s](%s)", userMd, repoMd, markdownV2(pl.Forkee.FullName), pl.Forkee.HtmlUrl)
	case "DeleteEvent":
		message = fmt.Sprintf("%s deleted %s %s at %s", userMd, markdownV2(pl.RefType), markdownV2(pl.Ref), repoMd)
	case "PublicEvent":
		message = fmt.Sprintf("%s made %s %s", userMd, repoMd, xsugar.IfThenElse(obj.Public, "public", "private"))

	case "IssuesEvent":
		message = fmt.Sprintf("%s %s issue [\\#%d](%s) in %s", userMd, pl.Action, pl.Issue.Number, pl.Issue.HtmlUrl, repoMd)
	case "IssueCommentEvent":
		message = fmt.Sprintf("%s %s comment [%d](%s) on issue [\\#%d](%s) in %s", userMd, pl.Action, pl.Comment.Id, pl.Comment.HtmlUrl, pl.Issue.Number, pl.Issue.HtmlUrl, repoMd)

	case "PullRequestEvent":
		message = fmt.Sprintf("%s %s pull request [\\#%d](%s) at %s", userMd, pl.Action, pl.Number, pl.PullRequest.HtmlUrl, repoMd)
	case "PullRequestReviewEvent":
		message = fmt.Sprintf("%s %s pull a request review in pull request [\\#%d](%s) at %s", userMd, pl.Action, pl.PullRequest.Number, pl.PullRequest.HtmlUrl, repoMd)
	case "PullRequestReviewCommentEvent":
		message = fmt.Sprintf("%s %s pull request review comment [%d](%s) in pull request [\\#%d](%s) at %s", userMd, pl.Action, pl.Comment.Id, pl.Comment.HtmlUrl, pl.PullRequest.Number, pl.PullRequest.HtmlUrl, repoMd)

	case "CommitCommentEvent":
		message = fmt.Sprintf("%s %s comment [%d](%s) at commit %s in %s", userMd, pl.Action, pl.Comment.Id, pl.Comment.HtmlUrl, pl.Comment.CommitId[0:7], repoMd)

	case "MemberEvent":
		message = fmt.Sprintf("%s %s member [%s](%s) to %s", userMd, pl.Action, markdownV2(pl.Member.Login), pl.Member.HtmlUrl, repoMd)
	case "ReleaseEvent":
		message = fmt.Sprintf("%s release [%s](%s) at %s", userMd, markdownV2(pl.Release.TagName), pl.Release.HtmlUrl, repoMd)
	case "GollumEvent":
		count := xsugar.IfThenElse(len(pl.Page) <= 1, "1 wiki page", fmt.Sprintf("%d wiki pages", len(pl.Page)))
		message = fmt.Sprintf("%s updated %s at %s", userMd, count, repoMd)
	default:
		event := strings.TrimRight(obj.Type, "Event")
		message = fmt.Sprintf("%s: %s %s", markdownV2(event), userMd, repoMd)
	}

	if !obj.Public {
		message += " (private)"
	}
	return message
}

func formatIssueEvent(obj *model.IssueEvent) string {
	userUrl := fmt.Sprintf("https://github.com/%s", obj.Actor.Login)
	repoUrl := fmt.Sprintf("https://github.com/%s", obj.Repo)
	issueUrl := fmt.Sprintf("https://github.com/%s/issues/%d", obj.Repo, obj.Number)
	userMd := fmt.Sprintf("[%s](%s)", markdownV2(obj.Actor.Login), userUrl)
	issueMd := fmt.Sprintf("[%s\\#%d](%s)", markdownV2(obj.Repo), obj.Number, issueUrl)

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
		message = fmt.Sprintf("%s renamed %s to `%s`", userMd, issueMd, markdownV2(obj.Rename.To))

	case "labeled":
		labelUrl := fmt.Sprintf("%s/labels/%s", repoUrl, obj.Label.Name)
		message = fmt.Sprintf("%s added label [%s](%s) to %s", userMd, markdownV2(obj.Label.Name), labelUrl, issueMd)
	case "unlabeled":
		labelUrl := fmt.Sprintf("%s/labels/%s", repoUrl, obj.Label.Name)
		message = fmt.Sprintf("%s removed label [%s](%s) from %s", userMd, markdownV2(obj.Label.Name), labelUrl, issueMd)

	case "locked":
		message = fmt.Sprintf("%s locked %s", userMd, issueMd)
	case "unlocked":
		message = fmt.Sprintf("%s unlocked %s", userMd, issueMd)

	case "milestoned":
		milestoneUrl := fmt.Sprintf("%s/milestones", repoUrl)
		message = fmt.Sprintf("%s added %s to milestone [%s](%s)", userMd, issueMd, markdownV2(obj.Milestone.Title), milestoneUrl)
	case "demilestoned":
		milestoneUrl := fmt.Sprintf("%s/milestones", repoUrl)
		message = fmt.Sprintf("%s removed %s from milestone [%s](%s)", userMd, issueMd, markdownV2(obj.Milestone.Title), milestoneUrl)

	case "pinned":
		message = fmt.Sprintf("%s pinned %s", userMd, issueMd)
	case "unpinned":
		message = fmt.Sprintf("%s unpinned %s", userMd, issueMd)

	case "assigned":
		message = fmt.Sprintf("%s is assigned to %s", userMd, issueMd)
	case "commented":
		message = fmt.Sprintf("%s added a [comment](%s) to %s", userMd, obj.HtmlUrl, issueMd)

	case "merged":
		message = fmt.Sprintf("%s merged pull request %s", userMd, issueMd)
	case "head_ref_deleted":
		message = fmt.Sprintf("%s deleted the head branch of %s", userMd, issueMd)
	case "head_ref_restored":
		message = fmt.Sprintf("%s restored the head branch of %s", userMd, issueMd)
	case "added_to_project":
		message = fmt.Sprintf("%s added %s to a project", userMd, issueMd)
	case "removed_from_project":
		message = fmt.Sprintf("%s removed %s from a project", userMd, issueMd)
	case "moved_columns_in_project":
		message = fmt.Sprintf("%s repository %s was moved to column in a project", userMd, issueMd)

	case "cross-referenced":
		mdShow := fmt.Sprintf("%s/%s#%d", obj.Source.Issue.Repository.Owner.Login, obj.Source.Issue.Repository.Name, obj.Source.Issue.Number)
		mdUrl := fmt.Sprintf("https://github.com/%s/%s/issues/%d", obj.Source.Issue.Repository.Owner.Login, obj.Source.Issue.Repository.Name, obj.Source.Issue.Number)
		targetMd := fmt.Sprintf("[%s](%s)", markdownV2(mdShow), mdUrl)
		message = fmt.Sprintf("%s mentioned %s from %s", userMd, issueMd, targetMd)
	case "referenced":
		toRepo := obj.CommitUrl // https://api.github.com/repos/gofiber/fiber/commits/a65d5027f336339cf4fe20cda0232c56cd64212e
		if toRepo == "" {
			message = fmt.Sprintf("%s added a commit that referenced %s", userMd, issueMd)
		} else if toSp := strings.Split(toRepo, "/"); len(toSp) >= 4 {
			commitMd := fmt.Sprintf("[%s](%s)", obj.CommitId[:7], obj.CommitUrl)
			toRepo = fmt.Sprintf("%s/%s", toSp[len(toSp)-4], toSp[len(toSp)-3])
			toRepoUrl := fmt.Sprintf("https://github.com/%s", toRepo)
			targetMd := fmt.Sprintf("[%s](%s)", markdownV2(toRepo), toRepoUrl)
			message = fmt.Sprintf("%s added a commit %s of %s that referenced %s", userMd, commitMd, targetMd, issueMd)
		} else {
			message = ""
		}
	default:
		message = fmt.Sprintf("%s: %s %s", markdownV2(obj.Event), userMd, issueMd)
	}

	return message
}
