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
	cmd.menu.cr.CommandRouter().AddComponentFunc("help_menu", cmd.InteractionHandler())

	return &api.CreateCommandData{
		Name:        cmd.GetName(),
		Description: cmd.GetDescription(),
		Type:        discord.ChatInputCommand,
	}, nil
}

func (cmd *HelpCommand) GenerateEmbed() {
	commands, err := cmd.s.Commands(cmd.c.Discord.AppIDSnowflake)
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
	e, ok := cmd.embeds[menu]
	if !ok {
		return discord.Embed{}
	}
	return *e
}

func (cmd *HelpCommand) GetSelect(selectedMenu string) *discord.ContainerComponents {
	var options []discord.SelectOption

	for _, menu := range cmd.menu.cr.GetMenus() {
		options = append(options, discord.SelectOption{
			Label:   fmt.Sprintf("%s %s", menu.GetEmoji().String(), menu.GetName()),
			Value:   menu.GetName(),
			Default: menu.GetName() == selectedMenu,
		})
	}

	return discord.ComponentsPtr(
		&discord.ActionRowComponent{
			&discord.StringSelectComponent{
				CustomID: "help_menu",
				Options:  options,
			},
		},
	)
}

func (cmd *HelpCommand) Func() cmdroute.CommandHandlerFunc {
	return func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
		if cmd.embeds == nil {
			cmd.GenerateEmbed()
		}

		var defaultMenu = "General"

		return &api.InteractionResponseData{
			Embeds:     &[]discord.Embed{cmd.GetEmbed(defaultMenu)},
			Components: cmd.GetSelect(defaultMenu),
		}
	}
}

func (cmd *HelpCommand) InteractionHandler() cmdroute.ComponentHandlerFunc {
	return func(ctx context.Context, data cmdroute.ComponentData) *api.InteractionResponse {
		if cmd.embeds == nil {
			cmd.GenerateEmbed()
		}

		menuName := data.Event.Data.(*discord.StringSelectInteraction).Values[0]

		err := cmd.s.RespondInteraction(data.Event.ID, data.Event.Token, api.InteractionResponse{
			Type: api.DeferredMessageUpdate,
		})
		if err != nil {
			logrus.WithError(err).WithField("menu", "help").Error("Failed to update message")
		}

		// Update the embed of the original message
		_, err = cmd.s.EditInteractionResponse(cmd.c.Discord.AppIDSnowflake, data.Event.Token, api.EditInteractionResponseData{
			Embeds:     &[]discord.Embed{cmd.GetEmbed(menuName)},
			Components: cmd.GetSelect(menuName),
		})
		if err != nil {
			logrus.WithError(err).WithField("menu", "help").Error("Failed to update message")
		}

		return nil
	}
}
