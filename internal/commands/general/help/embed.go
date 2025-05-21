package help

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/yyewolf/rwbyadv3/internal/utils"
)

func (cmd *helpCommand) generateEmbed() {
	commands, _ := cmd.app.Client().Rest().GetGlobalCommands(cmd.config.Discord.AppID, false)

	cmd.embeds = make(map[string]*discord.Embed)
	menus := cmd.menus.Menus

	for _, menu := range menus {
		embed := &discord.Embed{
			Title:  fmt.Sprintf("%s %s commands :", menu.Emoji.Name, menu.Name),
			Color:  cmd.app.Config().App.BotColor,
			Footer: cmd.app.Footer(),
		}

		for _, command := range menu.Commands {
			discordCmd := utils.FindCommandByName(commands, command.Name)

			subcommands := command.GetSubCommands()

			// This is not clickable if there are subcommands below
			if len(subcommands) == 0 {
				embed.Description += fmt.Sprintf("</%s:%s> - `%s`\n", discordCmd.Name(), discordCmd.ID(), command.Description)
			}

			for _, subcommand := range subcommands {
				embed.Description += fmt.Sprintf("</%s %s:%s> - `%s`\n", subcommand.Prefix, subcommand.Name, discordCmd.ID(), subcommand.Description)
			}
		}

		cmd.embeds[menu.Name] = embed
	}
}

func (cmd *helpCommand) getEmbed(menu string) discord.Embed {
	e, ok := cmd.embeds[menu]
	if !ok {
		return discord.Embed{}
	}
	return *e
}
