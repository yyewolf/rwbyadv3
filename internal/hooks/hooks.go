package hooks

import (
	"context"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/utils"
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

func listingsAfterInsert(app interfaces.App) func(ctx context.Context, exec boil.ContextExecutor, c *models.Listing) error {
	return func(ctx context.Context, exec boil.ContextExecutor, c *models.Listing) error {
		utils.App.DispatchNewListing(app, c)
		return nil
	}
}

func listingsAfterDelete(app interfaces.App) func(ctx context.Context, exec boil.ContextExecutor, c *models.Listing) error {
	return func(ctx context.Context, exec boil.ContextExecutor, c *models.Listing) error {
		utils.App.DispatchRemoveListing(app, c)
		return nil
	}
}

func auctionsAfterInsert(app interfaces.App) func(ctx context.Context, exec boil.ContextExecutor, c *models.Auction) error {
	return func(ctx context.Context, exec boil.ContextExecutor, c *models.Auction) error {
		utils.App.DispatchNewAuction(app, c)
		return nil
	}
}

func auctionsAfterUpdate(app interfaces.App) func(ctx context.Context, exec boil.ContextExecutor, c *models.Auction) error {
	return func(ctx context.Context, exec boil.ContextExecutor, c *models.Auction) error {
		utils.App.DispatchUpdateAuction(app, c)
		return nil
	}
}

func auctionsAfterDelete(app interfaces.App) func(ctx context.Context, exec boil.ContextExecutor, c *models.Auction) error {
	return func(ctx context.Context, exec boil.ContextExecutor, c *models.Auction) error {
		utils.App.DispatchRemoveAuction(app, c)
		return nil
	}
}

func bidAfterInsertOrUpdate(app interfaces.App) func(ctx context.Context, exec boil.ContextExecutor, c *models.AuctionsBid) error {
	return func(ctx context.Context, exec boil.ContextExecutor, c *models.AuctionsBid) error {
		utils.App.DispatchNewBid(app, c)
		return nil
	}
}

func RegisterHooks(app interfaces.App) {
	models.AddCardHook(boil.AfterInsertHook, cardAfterInsert)
	models.AddListingHook(boil.AfterInsertHook, listingsAfterInsert(app))
	models.AddListingHook(boil.AfterDeleteHook, listingsAfterDelete(app))
	models.AddAuctionHook(boil.AfterInsertHook, auctionsAfterInsert(app))
	models.AddAuctionHook(boil.AfterDeleteHook, auctionsAfterDelete(app))
	models.AddAuctionHook(boil.AfterUpdateHook, auctionsAfterUpdate(app))
	models.AddAuctionsBidHook(boil.AfterInsertHook, bidAfterInsertOrUpdate(app))
	models.AddAuctionsBidHook(boil.AfterUpdateHook, bidAfterInsertOrUpdate(app))
}
