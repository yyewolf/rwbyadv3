package profile

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/schema/enums"
	"github.com/yyewolf/rwbyadv3/internal/utils"
)

func (cmd *profileCommand) generator(currentPlayer *ent.Player, discordUser discord.User) discord.Embed {
	embed := discord.NewEmbedBuilder()

	embed.SetTitlef("These is your profile %s.", discordUser.EffectiveName())
	embed.SetColor(cmd.app.Config().App.BotColor)

	embed.AddField(
		"Player :",
		utils.Joinln(
			fmt.Sprintf("Level : **%d**", currentPlayer.Level),
			fmt.Sprintf("XP : **%d**/**%d**", currentPlayer.ExperiencePoints, currentPlayer.ExperiencePointsThreshold),
			fmt.Sprintf("Slots : **%d**", currentPlayer.BackpackLevel),
			fmt.Sprintf("Boxes : **%d**/**%d**", len(currentPlayer.Edges.Lootboxes), currentPlayer.BackpackLevel),
			fmt.Sprintf("Liens : **%d** (**%d** locked)", currentPlayer.Liens, currentPlayer.LiensInAuction),
		),
		true,
	)

	counts := currentPlayer.LootBoxesCount()

	embed.AddField(
		"Inventory :",
		utils.Joinln(
			fmt.Sprintf("Cards : **%d**/**%d** (**%d** reserved)", len(currentPlayer.Edges.Cards), utils.Players.MaxSlots(currentPlayer), currentPlayer.BackpackReservedSlots),
			fmt.Sprintf("Classic boxes : **%d**", counts[enums.LootBoxClassic]),
			fmt.Sprintf("Rare boxes : **%d**", counts[enums.LootBoxRare]),
			fmt.Sprintf("Limited boxes : **%d**", counts[enums.LootBoxLimited]),
			fmt.Sprintf("Special boxes : **%d**", counts[enums.LootBoxSpecial]),
		),
		true,
	)

	// Activities
	dungeonActivity := currentPlayer.GetDungeonState()

	embed.AddField(
		"Activities :",
		utils.Joinln(
			dungeonActivity.RenderUser(),
		),
		true,
	)

	embed.SetEmbedFooter(cmd.app.Footer())

	return embed.Build()
}
