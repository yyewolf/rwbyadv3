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
	for _, t := range models.AllLootBoxesType() {
		buttons = append(buttons,
			discord.NewButton(discord.ButtonStyleSecondary, fmt.Sprintf("%d %s boxes", counts[t], t), fmt.Sprintf(componentIdFmt, p.ID, t), ""),
		)
	}

	return discord.NewActionRow(
		buttons...,
	)
}
