package builder

import (
	"context"
	"errors"
	"reflect"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/yyewolf/rwbyadv3/models"
)

type ContextKey string

var (
	PlayerKey ContextKey = "player"
	ErrorKey  ContextKey
)

type Event interface {
	User() discord.User
	handler.CommandEvent | handler.AutocompleteEvent | handler.ComponentEvent | handler.InteractionEvent | handler.ModalEvent
}

type ContextBuilder struct {
	withPlayer            bool
	withPlayerGithubStars bool
}

type ContextOption func(a *ContextBuilder)

func FillContext[K Event](cb *ContextBuilder, event K, ctx context.Context) context.Context {
	if cb.withPlayer {
		var mods []qm.QueryMod

		if cb.withPlayerGithubStars {
			mods = append(mods, qm.Load(models.PlayerRels.GithubStar))
		}

		mods = append(mods,
			qm.Select("*"),
			qm.Where(models.PlayerColumns.ID+"=?", event.User().ID),
		)

		p, err := models.Players(mods...).OneG(ctx)
		if err != nil {
			ctx = context.WithValue(ctx, ErrorKey, err)
			return ctx
		}

		ctx = context.WithValue(ctx, PlayerKey, p)
	}

	return ctx
}

func WithContext[K Event](handler func(e *K) error, opts ...ContextOption) func(e *K) error {
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
			v = FillContext(&cb, *e, v)

			ctxVal.Set(reflect.ValueOf(v))
		}

		return handler(e)
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
