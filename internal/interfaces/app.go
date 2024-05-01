package interfaces

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/handler"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/repo"
)

type App interface {
	Start()

	Shutdown() error

	// Getter
	Config() *env.Config
	Handler() *handler.Mux
	Client() bot.Client

	// Github
	// NewGithubIssue(params repo.NewIssueParams) (*github.Issue, error)
	Github() *repo.GithubClient
}
