package open

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/schema/enums"
)

func (cmd *openCommand) generator(currentPlayer *ent.Player) discord.ContainerComponent {
	counts := currentPlayer.LootBoxesCount()

	var buttons []discord.InteractiveComponent
	for _, lootboxType := range enums.LootBoxTypes() {
		buttons = append(buttons,
			discord.NewSecondaryButton(
				fmt.Sprintf("%d %s boxes", counts[lootboxType], lootboxType),
				fmt.Sprintf(componentIdFmt, currentPlayer.ID, lootboxType),
			),
		)
	}

	return discord.NewActionRow(
		buttons...,
	)
}
