package general

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/middleware"
	"github.com/yyewolf/rwbyadv3/models"
)

type ProfileCommand struct {
	s *state.State
	c *env.Config

	menu *GeneralMenu

	interfaces.ContextGenerator
}

func newProfileCommand(m *GeneralMenu) interfaces.Command {
	return &ProfileCommand{
		c:                m.cr.Config(),
		s:                m.cr.State(),
		menu:             m,
		ContextGenerator: m.cr.ContextGenerator(),
	}
}

func (cmd *ProfileCommand) GetName() string {
	return "profile"
}

func (cmd *ProfileCommand) GetDescription() string {
	return "Displays your profile"
}

func (cmd *ProfileCommand) RegisterCommand() (*api.CreateCommandData, error) {
	cmd.menu.cr.CommandRouter().AddFunc(cmd.GetName(), middleware.ContextualizeCommandRouter(cmd, cmd.Func))

	return &api.CreateCommandData{
		Name:        cmd.GetName(),
		Description: cmd.GetDescription(),
		Type:        discord.ChatInputCommand,
	}, nil
}

func (cmd *ProfileCommand) Func(ctx interfaces.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	db := cmd.menu.cr.Database()
	p, err := db.Players().GetByDiscordID(data.Event.SenderID())
	if err != nil {
		p = &models.Player{
			DiscordID: data.Event.SenderID().String(),
		}
		err := db.Players().Create(p)
		if err != nil {
			return &api.InteractionResponseData{
				Content: option.NewNullableString("An error has occured"),
				Flags:   discord.EphemeralMessage,
			}
		}

		return &api.InteractionResponseData{
			Content: option.NewNullableString("Created : " + fmt.Sprintf("%+#v", p)),
		}
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString("Loaded : " + fmt.Sprintf("%+#v", p)),
	}
}
