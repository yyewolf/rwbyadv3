package repo

import (
	"context"

	"github.com/google/go-github/v61/github"
	"github.com/yyewolf/rwbyadv3/internal/env"
)

type GithubClient struct {
	config *env.Config
	client *github.Client
}

func NewGithubClient(config *env.Config) *GithubClient {
	c := github.NewClient(nil).WithAuthToken(config.Github.Token)

	return &GithubClient{
		config: config,
		client: c,
	}
}

type NewIssueParams struct {
	Title       string
	Description string
}

func (c *GithubClient) NewGithubIssue(params NewIssueParams) (*github.Issue, error) {
	issue := &github.IssueRequest{
		Title: &params.Title,
		Body:  &params.Description,
	}

	i, _, err := c.client.Issues.Create(context.TODO(), c.config.Github.Username, c.config.Github.Repository, issue)
	if err != nil {
		return nil, err
	}

	return i, nil
}
