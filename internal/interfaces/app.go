package interfaces

import (
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/yyewolf/rwbyadv3/internal/env"
)

type App interface {
	Start()

	Shutdown() error

	// Getter
	State() *state.State
	CommandRouter() *cmdroute.Router
	Config() *env.Config
}
