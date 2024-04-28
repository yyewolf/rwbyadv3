package interfaces

import (
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/google/go-github/v61/github"
	"github.com/yyewolf/rwbyadv3/internal/database"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/repo"
)

type App interface {
	Start()

	Shutdown() error

	// Getter
	State() *state.State
	CommandRouter() *cmdroute.Router
	Config() *env.Config
	Database() *database.Database

	// Github
	NewGithubIssue(params repo.NewIssueParams) (*github.Issue, error)
}
