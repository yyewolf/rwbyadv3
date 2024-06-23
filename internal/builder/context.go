package builder

import (
	"context"
	"errors"
	"reflect"
	"runtime"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
	"github.com/sirupsen/logrus"
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

func FillPlayerContext(cb *ContextBuilder, userID snowflake.ID, ctx context.Context) (context.Context, error) {
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
		qm.Where(models.PlayerColumns.ID+"=?", userID),
	)

	p, err := models.Players(mods...).OneG(ctx)
	if err != nil {
		logrus.WithError(err).Error("error when fetching player")
		return ctx, errors.New("auth error")
	}

	ctx = context.WithValue(ctx, PlayerKey, p)
	return ctx, nil
}

func FillContextReply[K Event](cb *ContextBuilder, event K, ctx context.Context) (context.Context, error) {
	var err error
	if cb.withPlayer {
		ctx, err = FillPlayerContext(cb, event.User().ID, ctx)
		if err != nil {
			event.CreateMessage(
				discord.NewMessageCreateBuilder().
					SetContentf("You cannot use this command yet... Try using %s first !", cb.app.CommandMention("begin")).
					SetEphemeral(true).
					Build(),
			)
			return ctx, errors.New("auth error")
		}
	}

	return ctx, nil
}

func getFuncClear(i interface{}) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	unclearPart := strings.Split(fullName, "/commands")[1]
	prefix := strings.Split(unclearPart, ".(")[0]
	suffix := strings.Split(unclearPart, ").")[1]
	suffix = strings.Split(suffix, "-")[0]

	return prefix + "/" + suffix
}

func WithContext[K Event](app interfaces.App, handler func(e *K) error, opts ...ContextOption) func(e *K) error {
	// Context builder
	var cb ContextBuilder
	cb.app = app

	for _, opt := range opts {
		opt(&cb)
	}

	funcName := getFuncClear(handler)

	return func(e *K) error {
		// Firstly, we extract the context
		// Try not to use to much reflection
		ctxVal := reflect.ValueOf(e).Elem().FieldByName("Ctx")
		if !ctxVal.IsValid() {
			return errors.New("invalid handler passed")
		}

		logrus.
			WithField("func", funcName).
			WithField("user_id", (*e).User().ID).
			Info("command")

		switch v := ctxVal.Interface().(type) {
		default:
			return errors.New("invalid handler passed")
		case context.Context:
			v, err := FillContextReply(&cb, *e, v)
			if err != nil {
				logrus.
					WithField("func", funcName).
					WithField("user_id", (*e).User().ID).
					WithError(err).
					Error("command")
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
	cb.app = app

	for _, opt := range opts {
		opt(&cb)
	}

	funcName := getFuncClear(handler)

	return func(d D, e *K) error {
		// Firstly, we extract the context
		// Try not to use to much reflection
		ctxVal := reflect.ValueOf(e).Elem().FieldByName("Ctx")
		if !ctxVal.IsValid() {
			return errors.New("invalid handler passed")
		}

		logrus.
			WithField("func", funcName).
			WithField("user_id", (*e).User().ID).
			Info("command")

		switch v := ctxVal.Interface().(type) {
		default:
			return errors.New("invalid handler passed")
		case context.Context:
			v, err := FillContextReply(&cb, *e, v)
			if err != nil {
				logrus.
					WithField("func", funcName).
					WithField("user_id", (*e).User().ID).
					WithError(err).
					Error("command")
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

func GetContext(ctx context.Context, app interfaces.App, userID snowflake.ID, opts ...ContextOption) (context.Context, error) {
	// Context builder
	var cb ContextBuilder
	cb.app = app

	for _, opt := range opts {
		opt(&cb)
	}

	return FillPlayerContext(&cb, userID, ctx)
}
