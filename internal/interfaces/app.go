package interfaces

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/handler"
	"github.com/google/go-github/v61/github"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/repo"
)

type App interface {
	Start()

	Shutdown() error

	// Getter
	Config() *env.Config
	Database() Database
	Handler() *handler.Mux
	Client() bot.Client

	// Github
	NewGithubIssue(params repo.NewIssueParams) (*github.Issue, error)
}
