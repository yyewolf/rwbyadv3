package builder

import (
	"context"
	"errors"
	"reflect"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/rest"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/models"
)

type ContextKey string

var (
	PlayerKey ContextKey = "player"
	ErrorKey  ContextKey
)

type Event interface {
	User() discord.User
	CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error
	handler.CommandEvent | handler.AutocompleteEvent | handler.ComponentEvent | handler.InteractionEvent | handler.ModalEvent
}

type ContextBuilder struct {
	app interfaces.App

	withPlayer             bool
	withPlayerGithubStars  bool
	withPlayerCards        bool
	withPlayerLootBoxes    bool
	withPlayerSelectedCard bool
}

type ContextOption func(a *ContextBuilder)

func FillContext[K Event](cb *ContextBuilder, event K, ctx context.Context) (context.Context, error) {
	if cb.withPlayer {
		var mods []qm.QueryMod

		if cb.withPlayerGithubStars {
			mods = append(mods, qm.Load(models.PlayerRels.GithubStar))
		}

		if cb.withPlayerCards {
			mods = append(mods,
				qm.Load(
					models.PlayerRels.PlayerCards,
					qm.OrderBy(models.PlayerCardColumns.Position),
				),
				qm.Load(
					qm.Rels(models.PlayerRels.PlayerCards, models.PlayerCardRels.Card, models.CardRels.CardsStat),
				),
			)
		}

		if cb.withPlayerLootBoxes {
			mods = append(mods, qm.Load(models.PlayerRels.LootBoxes))
		}

		if cb.withPlayerSelectedCard {
			mods = append(mods, qm.Load(models.PlayerRels.SelectedCard))
		}

		mods = append(mods,
			qm.Select("*"),
			qm.Where(models.PlayerColumns.ID+"=?", event.User().ID),
		)

		p, err := models.Players(mods...).OneG(ctx)
		if err != nil {
			event.CreateMessage(
				discord.NewMessageCreateBuilder().
					SetContentf("You cannot use this command yet... Try using %s first !", cb.app.CommandMention("begin")).
					SetEphemeral(true).
					Build(),
			)
			return ctx, errors.New("auth error")
		}

		ctx = context.WithValue(ctx, PlayerKey, p)
	}

	return ctx, nil
}

func WithContext[K Event](app interfaces.App, handler func(e *K) error, opts ...ContextOption) func(e *K) error {
	// Context builder
	var cb ContextBuilder

	for _, opt := range opts {
		opt(&cb)
	}

	return func(e *K) error {
		// Firstly, we extract the context
		// Try not to use to much reflection
		ctxVal := reflect.ValueOf(e).Elem().FieldByName("Ctx")
		if !ctxVal.IsValid() {
			return errors.New("invalid handler passed")
		}

		switch v := ctxVal.Interface().(type) {
		default:
			return errors.New("invalid handler passed")
		case context.Context:
			v, err := FillContext(&cb, *e, v)
			if err != nil {
				return nil
			}

			ctxVal.Set(reflect.ValueOf(v))
		}

		return handler(e)
	}
}

func WithContextD[D any, K Event](app interfaces.App, handler func(d D, e *K) error, opts ...ContextOption) func(d D, e *K) error {
	// Context builder
	var cb ContextBuilder

	for _, opt := range opts {
		opt(&cb)
	}

	return func(d D, e *K) error {
		// Firstly, we extract the context
		// Try not to use to much reflection
		ctxVal := reflect.ValueOf(e).Elem().FieldByName("Ctx")
		if !ctxVal.IsValid() {
			return errors.New("invalid handler passed")
		}

		switch v := ctxVal.Interface().(type) {
		default:
			return errors.New("invalid handler passed")
		case context.Context:
			v, err := FillContext(&cb, *e, v)
			if err != nil {
				return nil
			}

			ctxVal.Set(reflect.ValueOf(v))
		}

		return handler(d, e)
	}
}

func WithPlayer() func(a *ContextBuilder) {
	return func(a *ContextBuilder) {
		a.withPlayer = true
	}
}

func WithPlayerGithubStars() func(a *ContextBuilder) {
	return func(a *ContextBuilder) {
		a.withPlayerGithubStars = true
	}
}

func WithPlayerCards() func(a *ContextBuilder) {
	return func(a *ContextBuilder) {
		a.withPlayerCards = true
	}
}

func WithPlayerSelectedCard() func(a *ContextBuilder) {
	return func(a *ContextBuilder) {
		a.withPlayerSelectedCard = true
	}
}

func WithPlayerLootBoxes() func(a *ContextBuilder) {
	return func(a *ContextBuilder) {
		a.withPlayerLootBoxes = true
	}
}
