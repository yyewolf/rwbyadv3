package confirmation

import (
	"fmt"
	"math/rand"
	"path"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/utils"
)

type Handler struct {
	app interfaces.App

	prefix   string
	callback handler.ButtonComponentHandler
}

func NewHandler(app interfaces.App, prefix string, callback handler.ButtonComponentHandler) *Handler {
	return &Handler{
		app:      app,
		prefix:   prefix,
		callback: callback,
	}
}

func (c *Handler) SetupMux(h *handler.Mux) {
	h.ButtonComponent(path.Join("/", c.prefix, "/ok"), c.HandleOk)

	h.ButtonComponent(path.Join("/", c.prefix, "/not_ok/{i}"), c.HandleNotOk)
}

var emojis = []*discord.Emoji{
	{
		Name: "🐶",
	},
	{
		Name: "🐱",
	},
	{
		Name: "🐭",
	},
	{
		Name: "🐼",
	},
}

func (c *Handler) AskForConfirmation(e *handler.CommandEvent, danger string, prefix string) error {
	var correct = rand.Intn(len(emojis))

	var buttons []discord.InteractiveComponent

	for i, emoji := range emojis {
		customID := path.Join(prefix, fmt.Sprintf("/not_ok/%d", i))
		if i == correct {
			customID = path.Join(prefix, "ok")
		}

		buttons = append(buttons,
			discord.
				NewSecondaryButton("", customID).
				WithEmoji(discord.ComponentEmoji{
					Name: emoji.Name,
				}),
		)
	}

	row := discord.NewActionRow(buttons...)

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			SetEmbeds(
				discord.NewEmbedBuilder().
					SetTitle("Confirmation").
					SetDescription(
						utils.Joinln(
							danger,
							"",
							fmt.Sprintf("You must click on %s to continue.", emojis[correct].Name),
						),
					).
					SetColor(c.app.Config().App.BotColor).
					SetEmbedFooter(c.app.Footer()).
					Build(),
			).
			AddContainerComponents(row),
	)
}

func (c *Handler) HandleOk(data discord.ButtonInteractionData, e *handler.ComponentEvent) error {
	return c.callback(data, e)
}

func (c *Handler) HandleNotOk(data discord.ButtonInteractionData, e *handler.ComponentEvent) error {
	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			SetEmbeds(
				discord.NewEmbedBuilder().
					SetTitle("Nope").
					SetDescription("You did not click the correct emoji !").
					SetColor(c.app.Config().App.BotColor).
					SetEmbedFooter(c.app.Footer()).
					Build(),
			).
			SetEphemeral(true),
	)
}
