package system

import (
	"context"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

type PingCommand struct {
	s *state.State
	c *env.Config

	menu *SystemMenu
}

func newPingCommand(m *SystemMenu) interfaces.Command {
	return &PingCommand{
		c:    m.cr.Config(),
		s:    m.cr.State(),
		menu: m,
	}
}

func (cmd *PingCommand) GetName() string {
	return "ping"
}

func (cmd *PingCommand) GetDescription() string {
	return "pong!"
}

func (cmd *PingCommand) RegisterCommand() (*api.CreateCommandData, error) {
	cmd.menu.cr.CommandRouter().AddFunc(cmd.GetName(), cmd.Func())

	return &api.CreateCommandData{
		Name:        cmd.GetName(),
		Description: cmd.GetDescription(),
		Type:        discord.ChatInputCommand,
	}, nil
}

func (cmd *PingCommand) Func() cmdroute.CommandHandlerFunc {
	return func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
		return &api.InteractionResponseData{
			Content: option.NewNullableString("Pong!"),
		}
	}
}
