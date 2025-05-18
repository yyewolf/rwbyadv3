package hooks

import (
	"context"
	"fmt"

	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/utils"
)

func listingMutator(app interfaces.App) func(next ent.Mutator) ent.Mutator {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			listingMutation, ok := m.(*ent.ListingMutation)
			if !ok {
				return nil, fmt.Errorf("listing mutation expected, got %T", m)
			}

			id, ok := listingMutation.ID()
			if !ok {
				return nil, fmt.Errorf("listing ID expected")
			}

			switch m.Op() {
			case ent.OpCreate:
				utils.App.DispatchNewListing(app, id)
			case ent.OpDeleteOne:
				utils.App.DispatchRemoveListing(app, id)
			}

			return next.Mutate(ctx, m)
		})
	}
}

func auctionMutator(app interfaces.App) func(next ent.Mutator) ent.Mutator {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			auctionMutation, ok := m.(*ent.AuctionMutation)
			if !ok {
				return nil, fmt.Errorf("auction mutation expected, got %T", m)
			}

			id, ok := auctionMutation.ID()
			if !ok {
				return nil, fmt.Errorf("auction ID expected")
			}

			switch m.Op() {
			case ent.OpCreate:
				utils.App.DispatchNewAuction(app, id)
			case ent.OpDeleteOne:
				utils.App.DispatchRemoveAuction(app, id)
			case ent.OpUpdateOne:
				utils.App.DispatchUpdateAuction(app, id)
			}

			return next.Mutate(ctx, m)
		})
	}
}

func bidMutator(app interfaces.App) func(next ent.Mutator) ent.Mutator {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			bidMutation, ok := m.(*ent.AuctionBidMutation)
			if !ok {
				return nil, fmt.Errorf("bid mutation expected, got %T", m)
			}

			id, ok := bidMutation.ID()
			if !ok {
				return nil, fmt.Errorf("bid ID expected")
			}

			switch m.Op() {
			case ent.OpCreate:
				utils.App.DispatchNewBid(app, id)
			case ent.OpUpdateOne:
				utils.App.DispatchNewBid(app, id)
			}

			return next.Mutate(ctx, m)
		})
	}
}

func RegisterHooks(app interfaces.App) {
	app.Db().Listing.Use(listingMutator(app))
	app.Db().Auction.Use(auctionMutator(app))
	app.Db().AuctionBid.Use(bidMutator(app))
}
