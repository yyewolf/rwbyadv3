package open

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
)

func (cmd *openCommand) generator(p *models.Player) discord.ContainerComponent {
	counts := utils.Players.LootBoxesCount(p)

	var buttons []discord.InteractiveComponent
	for _, lootboxType := range models.AllLootBoxesType() {
		buttons = append(buttons,
			discord.NewSecondaryButton(
				fmt.Sprintf("%d %s boxes", counts[lootboxType], lootboxType),
				fmt.Sprintf(componentIdFmt, p.ID, lootboxType),
			),
		)
	}

	return discord.NewActionRow(
		buttons...,
	)
}
