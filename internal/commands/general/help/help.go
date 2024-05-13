package help

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

const (
	commandName        = "help"
	commandDescription = "Displays the help menu"

	componentId = "help_menu"
)

type helpCommand struct {
	c   *env.Config
	app interfaces.App
	ms  *builder.MenuStore

	embeds map[string]*discord.Embed
}

func HelpCommand(ms *builder.MenuStore, app interfaces.App) *builder.Command {
	var cmd helpCommand

	cmd.app = app
	cmd.ms = ms
	cmd.c = app.Config()

	return builder.NewCommand(
		builder.WithCommandName(commandName),
		builder.WithDescription(commandDescription),
		builder.WithRegisterFunc(func(h *handler.Mux) error {
			h.Command("/"+commandName, builder.WithContext(app, cmd.HandleCommand))
			h.SelectMenuComponent("/"+componentId, builder.WithContextD(app, cmd.HandleInteraction))
			return nil
		}),
		builder.WithSlashCommand(discord.SlashCommandCreate{
			Name:        commandName,
			Description: commandDescription,
		}),
	)
}

func (cmd *helpCommand) GetSelect(selectedMenu string) []discord.ContainerComponent {
	var options []discord.StringSelectMenuOption

	for _, menu := range cmd.ms.Menus {
		opt := discord.NewStringSelectMenuOption(
			fmt.Sprintf("%s %s", menu.Emoji.Name, menu.Name),
			menu.Name,
		)
		opt.Default = menu.Name == selectedMenu
		options = append(options, opt)
	}

	return []discord.ContainerComponent{
		&discord.ActionRowComponent{
			&discord.StringSelectMenuComponent{
				CustomID: componentId,
				Options:  options,
			},
		},
	}
}

func (cmd *helpCommand) HandleCommand(e *handler.CommandEvent) error {
	if cmd.embeds == nil {
		cmd.generateEmbed()
	}

	var defaultMenu = "General"

	var s = cmd.GetSelect(defaultMenu)

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			AddEmbeds(cmd.getEmbed(defaultMenu)).
			AddContainerComponents(s...),
	)
}

func (cmd *helpCommand) HandleInteraction(data discord.SelectMenuInteractionData, e *handler.ComponentEvent) error {
	if cmd.embeds == nil {
		cmd.generateEmbed()
	}

	menuName := data.(discord.StringSelectMenuInteractionData).Values[0]

	err := e.DeferUpdateMessage()
	if err != nil {
		logrus.WithError(err).WithField("menu", "help").Error("Failed to update message")
	}

	var s = cmd.GetSelect(menuName)

	// Update the embed of the original message
	_, err = e.UpdateInteractionResponse(
		discord.NewMessageUpdateBuilder().
			AddEmbeds(cmd.getEmbed(menuName)).
			AddContainerComponents(s...).
			Build(),
	)
	if err != nil {
		logrus.WithError(err).WithField("menu", "help").Error("Failed to update message")
	}

	return nil
}
