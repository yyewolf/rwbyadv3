package profile

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
)

func (cmd *profileCommand) generator(p *models.Player, u discord.User) discord.Embed {
	embed := discord.NewEmbedBuilder()

	embed.SetTitlef("These is your profile %s.", u.EffectiveName())
	embed.SetColor(cmd.app.Config().App.BotColor)

	embed.AddField(
		"Player :",
		utils.Joinln(
			fmt.Sprintf("Level : **%d**", p.Level),
			fmt.Sprintf("XP : **%d**/**%d**", p.XP, p.NextLevelXP),
			fmt.Sprintf("Slots : **%d**", p.BackpackLevel),
			fmt.Sprintf("Boxes : **%d**/**%d**", len(p.R.LootBoxes), p.BackpackLevel),
			fmt.Sprintf("Liens : **%d**", p.Liens),
		),
		true,
	)

	counts := utils.Players.LootBoxesCount(p)

	embed.AddField(
		"Inventory :",
		utils.Joinln(
			fmt.Sprintf("Cards : **%d**/**%d**", len(p.R.PlayerCards), p.BackpackLevel),
			fmt.Sprintf("Classic boxes : **%d**", counts[models.LootBoxesTypeClassic]),
			fmt.Sprintf("Rare boxes : **%d**", counts[models.LootBoxesTypeRare]),
			fmt.Sprintf("Limited boxes : **%d**", counts[models.LootBoxesTypeLimited]),
			fmt.Sprintf("Special boxes : **%d**", counts[models.LootBoxesTypeSpecial]),
		),
		true,
	)

	embed.SetEmbedFooter(cmd.app.Footer())

	return embed.Build()
}
