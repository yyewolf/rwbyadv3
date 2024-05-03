package repo

import (
	"context"
	"fmt"
	"net/http"

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

func (c *GithubClient) GetTokenUser(token string) (*github.User, error) {
	tempClient := github.NewClient(nil).WithAuthToken(token)

	i, _, err := tempClient.Users.Get(context.TODO(), c.config.Github.Username)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (c *GithubClient) CheckTokenUserStar(token, owner, repo string) (bool, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/user/starred/%s/%s", owner, repo), nil)
	if err != nil {
		return false, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	return response.StatusCode == http.StatusNoContent, nil
}
