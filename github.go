package main

import (
	"log"
	"time"

	"github.com/google/go-github/github"
)

// Func to check if a repository has open issues or not
func hasOpenIssues(client *github.Client, owner, repo string) (int, error) {
	repo_details, _, err := client.Repositories.Get(owner, repo)
	if err != nil {
		log.Fatal(err)
	}
	return *repo_details.OpenIssuesCount, nil
}

// Func to get issue count based on difference in days
func getIssueCount(client *github.Client, owner, repo string) (int, int, int, int, error) {
	num_issues_24h := 0  // Number of open issues that were opened in the last 24 hours
	num_issues_7d := 0   // Number of open issues that were opened more than 24 hours ago but less than 7 days ago
	num_issues_gt7d := 0 // Number of open issues that were opened more than 7 days ago

	opt := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	var allissues []*github.Issue

	for {
		issues, resp, err := client.Issues.ListByRepo(owner, repo, opt)
		if err != nil {
			log.Fatal(err)
		}
		allissues = append(allissues, issues...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}

	for _, issue := range allissues {
		if !ifPullRequest(issue) {
			daydiff := getDaysDiff(issue.CreatedAt)
			switch {
			case daydiff <= 1:
				num_issues_24h = num_issues_24h + 1
			case 1 < daydiff && daydiff <= 7:
				num_issues_7d = num_issues_7d + 1
			case daydiff > 7:
				num_issues_gt7d = num_issues_gt7d + 1
			}
		}
	}
	total_issues := num_issues_24h + num_issues_7d + num_issues_gt7d
	return num_issues_24h, num_issues_7d, num_issues_gt7d, total_issues, nil
}

// Func to get difference in number of days.
func getDaysDiff(issuetime *time.Time) int {
	now := time.Now().UTC()
	diff := now.Sub(*issuetime)
	days := int(diff.Hours() / 24)
	return days
}

// Func to check if the listed issue is a pull request or issue (Github api lists pull request along with issues in https://developer.github.com/v3/issues/#list-issues)
func ifPullRequest(issue *github.Issue) bool {
	if issue.PullRequestLinks != nil {
		return true
	}
	return false
}
