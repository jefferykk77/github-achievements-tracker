package tracker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// PRCommit represents a commit inside a PR.
type PRCommit struct {
	Commit struct {
		Message string `json:"message"`
	} `json:"commit"`
}

// PullRequest represents a GitHub PR node.
type PullRequest struct {
	Number     int       `json:"number"`
	Title      string    `json:"title"`
	MergedAt   time.Time `json:"mergedAt"`
	Repository struct {
		NameWithOwner string `json:"nameWithOwner"`
		IsPrivate     bool   `json:"isPrivate"`
	} `json:"repository"`
	Commits struct {
		Nodes []PRCommit `json:"nodes"`
	} `json:"commits"`
}

// GraphQLResponse represents the response structure of the GraphQL query.
type GraphQLResponse struct {
	Data struct {
		User struct {
			PullRequests struct {
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
				Nodes []PullRequest `json:"nodes"`
			} `json:"pullRequests"`
		} `json:"user"`
	} `json:"data"`
}

// Achievement represents a tracked GitHub Achievement.
type Achievement struct {
	Name        string
	Description string
	Levels      []int // Thresholds for each level: Level 1 (Default), Level 2 (Bronze), Level 3 (Silver), Level 4 (Gold)
	Units       string
}

// Progress represents the progress of an achievement.
type Progress struct {
	Achievement  Achievement
	PublicCount  int
	PrivateCount int
	TotalCount   int
}

// GetLevel returns the current level (0-4), the current milestone value, the next milestone value, and progress percentage.
func (p *Progress) GetLevel(usePrivate bool) (level int, currentVal int, targetVal int, percent float64) {
	count := p.PublicCount
	if usePrivate {
		count = p.TotalCount
	}
	currentVal = count

	// Check levels
	for i, threshold := range p.Achievement.Levels {
		if count >= threshold {
			level = i + 1
		} else {
			break
		}
	}

	if level == 0 {
		targetVal = p.Achievement.Levels[0]
		percent = float64(count) / float64(targetVal) * 100
	} else if level < len(p.Achievement.Levels) {
		prevThreshold := p.Achievement.Levels[level-1]
		targetVal = p.Achievement.Levels[level]
		diff := targetVal - prevThreshold
		progressDiff := count - prevThreshold
		percent = float64(progressDiff) / float64(diff) * 100
	} else {
		targetVal = p.Achievement.Levels[len(p.Achievement.Levels)-1]
		percent = 100.0
	}

	if percent > 100.0 {
		percent = 100.0
	}
	return
}

// FetchPullRequests fetches all merged PRs for a given user using the gh CLI.
func FetchPullRequests(username string, token string) ([]PullRequest, error) {
	var allPRs []PullRequest
	cursor := ""
	hasNext := true

	query := `query($username: String!, $cursor: String) {
		user(login: $username) {
			pullRequests(states: MERGED, first: 100, after: $cursor, orderBy: {field: CREATED_AT, direction: DESC}) {
				pageInfo {
					hasNextPage
					endCursor
				}
				nodes {
					number
					title
					mergedAt
					repository {
						nameWithOwner
						isPrivate
					}
					commits(first: 100) {
						nodes {
							commit {
								message
							}
						}
					}
				}
			}
		}
	}`

	for hasNext {
		// Execute gh api graphql
		args := []string{"api", "graphql", "-f", fmt.Sprintf("query=%s", query), "-F", fmt.Sprintf("username=%s", username)}
		if cursor != "" {
			args = append(args, "-F", fmt.Sprintf("cursor=%s", cursor))
		}
		cmd := exec.Command("gh", args...)
		
		// Set GITHUB_TOKEN if provided
		if token != "" {
			cmd.Env = append(cmd.Environ(), fmt.Sprintf("GITHUB_TOKEN=%s", token))
		}

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			return nil, fmt.Errorf("gh command failed: %s; error: %w", stderr.String(), err)
		}

		var resp GraphQLResponse
		if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
			return nil, fmt.Errorf("failed to parse JSON response: %w", err)
		}

		allPRs = append(allPRs, resp.Data.User.PullRequests.Nodes...)
		hasNext = resp.Data.User.PullRequests.PageInfo.HasNextPage
		cursor = resp.Data.User.PullRequests.PageInfo.EndCursor
	}

	return allPRs, nil
}

// CalculateAchievements processes the list of PRs and returns achievements progress.
func CalculateAchievements(prs []PullRequest) (pullShark Progress, pairExtra Progress) {
	pullShark = Progress{
		Achievement: Achievement{
			Name:        "Pull Shark",
			Description: "Opened pull requests that have been merged.",
			Levels:      []int{2, 16, 128, 1024},
			Units:       "PRs",
		},
	}

	pairExtra = Progress{
		Achievement: Achievement{
			Name:        "Pair Extraordinaire",
			Description: "Coauthored commits on merged pull requests.",
			Levels:      []int{1, 10, 24, 48},
			Units:       "PRs",
		},
	}

	for _, pr := range prs {
		isCoauthored := false
		for _, commitNode := range pr.Commits.Nodes {
			if strings.Contains(commitNode.Commit.Message, "Co-authored-by:") {
				isCoauthored = true
				break
			}
		}

		// Pull Shark counts
		if pr.Repository.IsPrivate {
			pullShark.PrivateCount++
		} else {
			pullShark.PublicCount++
		}
		pullShark.TotalCount++

		// Pair Extraordinaire counts
		if isCoauthored {
			if pr.Repository.IsPrivate {
				pairExtra.PrivateCount++
			} else {
				pairExtra.PublicCount++
			}
			pairExtra.TotalCount++
		}
	}

	return
}
