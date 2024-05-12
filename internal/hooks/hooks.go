package hooks

import (
	"context"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/yyewolf/rwbyadv3/models"
)

func cardAfterInsert(ctx context.Context, exec boil.ContextExecutor, c *models.Card) error {
	amount, err := models.PlayerCards(models.PlayerCardWhere.PlayerID.EQ(c.PlayerID)).Count(ctx, exec)
	if err != nil {
		return err
	}

	playerCard := models.PlayerCard{
		PlayerID: c.PlayerID,
		CardID:   c.ID,
		Position: int(amount),
	}

	return playerCard.Insert(ctx, exec, boil.Infer())
}

func RegisterHooks() {
	// Register my before insert hook for pilots
	models.AddCardHook(boil.AfterInsertHook, cardAfterInsert)
}
