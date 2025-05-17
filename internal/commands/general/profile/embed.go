package profile

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/schema/enums"
	"github.com/yyewolf/rwbyadv3/internal/utils"
)

func (cmd *profileCommand) generator(p *ent.Player, u discord.User) discord.Embed {
	embed := discord.NewEmbedBuilder()

	embed.SetTitlef("These is your profile %s.", u.EffectiveName())
	embed.SetColor(cmd.app.Config().App.BotColor)

	embed.AddField(
		"Player :",
		utils.Joinln(
			fmt.Sprintf("Level : **%d**", p.Level),
			fmt.Sprintf("XP : **%d**/**%d**", p.ExperiencePoints, p.ExperiencePointsThreshold),
			fmt.Sprintf("Slots : **%d**", p.BackpackLevel),
			fmt.Sprintf("Boxes : **%d**/**%d**", len(p.Edges.Lootboxes), p.BackpackLevel),
			fmt.Sprintf("Liens : **%d** (**%d** locked)", p.Liens, p.LiensInAuction),
		),
		true,
	)

	counts := p.LootBoxesCount()

	embed.AddField(
		"Inventory :",
		utils.Joinln(
			fmt.Sprintf("Cards : **%d**/**%d** (**%d** reserved)", len(p.Edges.Cards), utils.Players.NewMaxSlots(p), p.BackpackReservedSlots),
			fmt.Sprintf("Classic boxes : **%d**", counts[enums.LootBoxClassic]),
			fmt.Sprintf("Rare boxes : **%d**", counts[enums.LootBoxRare]),
			fmt.Sprintf("Limited boxes : **%d**", counts[enums.LootBoxLimited]),
			fmt.Sprintf("Special boxes : **%d**", counts[enums.LootBoxSpecial]),
		),
		true,
	)

	// Activities
	dungeonActivity := p.GetDungeonState()

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
