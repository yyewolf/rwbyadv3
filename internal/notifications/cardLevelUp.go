package notifications

import (
	"context"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

type CardLevelUpParams struct {
	Player *models.Player
	Card   *models.Card
}

func DispatchCardLevelUp(app interfaces.App, p *models.Player, c *models.Card) {
	workflowOptions := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("card_level_up_%s_%d", c.ID, c.Level),
		TaskQueue: app.Config().Temporal.TaskQueue,
	}
	app.Temporal().ExecuteWorkflow(context.Background(), workflowOptions, Repository.NotifyCardLevelUpWorkflow, &CardLevelUpParams{
		Player: p,
		Card:   p.R.SelectedCard,
	})
}

func (n *NotificationsRepository) NotifyCardLevelUpWorkflow(ctx workflow.Context, params *CardLevelUpParams) error {
	// TODO : Add check for DMs, GuildChannels, and if the user wants the notification at all

	c := n.app.Client()
	ch, err := c.Rest().CreateDMChannel(snowflake.MustParse(params.Player.ID))
	if err != nil {
		return err
	}
	primitive := utils.Cards.Primitive(params.Card)
	_, err = c.Rest().CreateMessage(ch.ID(), discord.NewMessageCreateBuilder().
		SetEmbeds(
			discord.NewEmbedBuilder().
				SetTitle("Congratulations !").
				SetDescriptionf(
					utils.Joinln(
						"<@%s>, your **%s** has leveled up!",
						"Level : **%d**.",
						"XP : **%d/%d**",
					),
					params.Player.ID, primitive.Name,
					params.Card.Level,
					params.Card.XP, params.Card.NextLevelXP,
				).
				SetThumbnail(utils.Cards.IconURI(params.Card)).
				SetColor(n.app.Config().App.BotColor).
				Build(),
		).
		Build(),
	)
	return err
}
