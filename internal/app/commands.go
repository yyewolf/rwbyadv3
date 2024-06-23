package app

import (
	"github.com/disgoorg/disgo/discord"
)

func (a *App) loadCommandMentions() error {
	commands, err := a.Client().Rest().GetGlobalCommands(a.Config().Discord.AppID, false)
	if err != nil {
		return err
	}

	a.commandMentions = make(map[string]string)

	for _, c := range commands {
		if c.Type() != discord.ApplicationCommandTypeSlash {
			continue
		}

		command := c.(discord.SlashCommand)
		a.commandMentions[command.Name()] = command.Mention()

		for _, opt := range command.Options {
			switch opt.Type() {
			case discord.ApplicationCommandOptionTypeSubCommandGroup:
				a.commandMentions[command.Name()+" "+opt.OptionName()] = discord.SlashCommandMention(command.ID(), command.Name()+" "+opt.OptionName())

				for _, opt := range command.Options {
					if opt.Type() != discord.ApplicationCommandOptionTypeSubCommand {
						continue
					}
					a.commandMentions[command.Name()+" "+opt.OptionName()] = discord.SlashCommandMention(command.ID(), command.Name()+" "+opt.OptionName())
				}
			case discord.ApplicationCommandOptionTypeSubCommand:
				a.commandMentions[command.Name()+" "+opt.OptionName()] = discord.SlashCommandMention(command.ID(), command.Name()+" "+opt.OptionName())
			}
		}
	}
	return nil
}

func (a *App) Footer() *discord.EmbedFooter {
	u, _ := a.Client().Rest().GetCurrentUser("")
	return &discord.EmbedFooter{
		Text:    "Made by Yewolf - Support: https://discord.gg/adJGyVxv7H",
		IconURL: u.EffectiveAvatarURL(),
	}
}
