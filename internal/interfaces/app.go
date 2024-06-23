package interfaces

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/repo"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type App interface {
	Start()
	Shutdown() error

	// Getter
	Config() *env.Config
	Handler() *handler.Mux
	Client() bot.Client
	EventHandler() JobHandler

	// Temporal
	Temporal() client.Client
	Worker() worker.Worker

	// Commands
	CommandMention(c string) string

	// Embeds
	Footer() *discord.EmbedFooter

	// Github
	// NewGithubIssue(params repo.NewIssueParams) (*github.Issue, error)
	Github() *repo.GithubClient
}
