package general

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/utils"
)

type HelpCommand struct {
	s *state.State
	c *env.Config

	embeds map[string]*discord.Embed

	menu *GeneralMenu
}

func newHelpCommand(m *GeneralMenu) interfaces.Command {
	return &HelpCommand{
		c:    m.cr.Config(),
		s:    m.cr.State(),
		menu: m,
	}
}

func (cmd *HelpCommand) GetName() string {
	return "help"
}

func (cmd *HelpCommand) GetDescription() string {
	return "Displays the help menu"
}

func (cmd *HelpCommand) RegisterCommand() (*api.CreateCommandData, error) {
	cmd.menu.cr.CommandRouter().AddFunc(cmd.GetName(), cmd.Func())

	return &api.CreateCommandData{
		Name:        cmd.GetName(),
		Description: cmd.GetDescription(),
		Type:        discord.ChatInputCommand,
	}, nil
}

func (cmd *HelpCommand) GenerateEmbed() {
	commands, err := cmd.s.Commands(discord.AppID(cmd.c.Discord.AppIDSnowflake))
	if err != nil {
		logrus.WithError(err).Error("Failed to get commands from discord")
	}

	cmd.embeds = make(map[string]*discord.Embed)
	menus := cmd.menu.cr.GetMenus()

	for _, menu := range menus {
		embed := &discord.Embed{
			Title: fmt.Sprintf("%s %s commands :", menu.GetEmoji().String(), menu.GetName()),
			Color: 0x00ff00,
		}

		for _, command := range menu.GetSubcommands() {
			discordCmd := utils.FindCommandByName(commands, command.GetName())

			embed.Description += fmt.Sprintf("</%s:%s> - `%s`\n", discordCmd.Name, discordCmd.ID, command.GetDescription())
		}

		cmd.embeds[menu.GetName()] = embed
	}
}

func (cmd *HelpCommand) GetEmbed(menu string) discord.Embed {
	return *cmd.embeds[menu]
}

func (cmd *HelpCommand) Func() cmdroute.CommandHandlerFunc {
	return func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
		if cmd.embeds == nil {
			cmd.GenerateEmbed()
		}

		return &api.InteractionResponseData{
			Embeds: &[]discord.Embed{cmd.GetEmbed("General")},
		}
	}
}
