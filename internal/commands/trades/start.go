package trades

import (
	"net/url"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/internal/utils"
)

func (cmd *tradesCommand) Start(logger *logrus.Entry, event *handler.CommandEvent) error {
	// currentPlayer := event.Ctx.Value(builder.NewPlayerKey).(*ent.Player)
	target := event.SlashCommandInteractionData().User("with")

	uri, err := url.Parse(cmd.app.Config().App.BaseURI)
	if err != nil {
		return utils.CommandError(logger, event, err)
	}
	uri.Path, _ = url.JoinPath(uri.Path, "/trade/")
	query := uri.Query()
	query.Set("playerId", target.ID.String())
	uri.RawQuery = query.Encode()

	return event.CreateMessage(
		discord.NewMessageCreateBuilder().
			SetContent("Click the button to trade with " + target.Mention() + "!").
			AddActionRow(
				discord.NewLinkButton("View Trade", uri.String()),
			).
			SetEphemeral(true).
			Build(),
	)
}
