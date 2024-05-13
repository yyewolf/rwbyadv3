package help

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/yyewolf/rwbyadv3/internal/utils"
)

func (cmd *helpCommand) generateEmbed() {
	commands, _ := cmd.app.Client().Rest().GetGlobalCommands(cmd.c.Discord.AppID, false)

	cmd.embeds = make(map[string]*discord.Embed)
	menus := cmd.ms.Menus

	for _, menu := range menus {
		embed := &discord.Embed{
			Title:  fmt.Sprintf("%s %s commands :", menu.Emoji.Name, menu.Name),
			Color:  cmd.app.Config().App.BotColor,
			Footer: cmd.app.Footer(),
		}

		for _, command := range menu.Commands {
			discordCmd := utils.FindCommandByName(commands, command.Name)

			embed.Description += fmt.Sprintf("</%s:%s> - `%s`\n", discordCmd.Name(), discordCmd.ID(), command.Description)
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
