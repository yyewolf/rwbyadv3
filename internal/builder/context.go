package builder

import (
	"context"
	"errors"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/card"
	"github.com/yyewolf/rwbyadv3/ent/player"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

type ContextKey string

var (
	NewPlayerKey ContextKey = "new_player"
	ErrorKey     ContextKey = "error"
	ContextIdKey ContextKey = "context_id"
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
	withPlayerAuctions     bool
	withPlayerLootBoxes    bool
	withPlayerSelectedCard bool
	withPlayerLimits       bool
	withPlayerDungeons     bool
	withPlayerDaily        bool
}

type ContextOption func(a *ContextBuilder)

func FillPlayerContext(builder *ContextBuilder, userID snowflake.ID, ctx context.Context) (context.Context, error) {
	var query = builder.app.Db().Player.Query().Where(player.ID(userID.String()))

	if builder.withPlayerGithubStars {
		query.WithGithubStar()
	}

	if builder.withPlayerCards {
		query.WithCards(func(q *ent.CardQuery) {
			q.Order(card.ByPosition())
			q.WithStats()
			q.WithType()
		})
	}

	if builder.withPlayerAuctions {
		query.WithAuctions(func(q *ent.AuctionQuery) {
			q.WithCard(func(q *ent.CardQuery) {
				q.WithType()
				q.WithStats()
			})
		})
	}

	if builder.withPlayerLootBoxes {
		query.WithLootboxes()
	}

	if builder.withPlayerSelectedCard {
		query.WithSelectedCard(func(q *ent.CardQuery) {
			q.WithType()
			q.WithStats()
		})
	}

	if builder.withPlayerLimits {
		query.WithLimits()
	}

	if builder.withPlayerDungeons {
		query.WithDungeons()
	}

	if builder.withPlayerDaily {
		query.WithDaily()
	}

	np, err := query.First(ctx)
	if err != nil {
		logrus.WithError(err).Error("error when fetching player")
		return ctx, errors.New("auth error")
	}

	ctx = context.WithValue(ctx, NewPlayerKey, np)
	return ctx, nil
}

func FillContextReply[K Event](builder *ContextBuilder, event K, ctx context.Context) (context.Context, error) {
	var err error
	if builder.withPlayer {
		ctx, err = FillPlayerContext(builder, event.User().ID, ctx)
		if err != nil {
			event.CreateMessage(
				discord.NewMessageCreateBuilder().
					SetContentf("You cannot use this command yet... Try using %s first !", builder.app.CommandMention("begin")).
					SetEphemeral(true).
					Build(),
			)
			return ctx, errors.New("auth error")
		}
	}

	// Add UUID to track the context
	ctx = context.WithValue(ctx, ContextIdKey, uuid.NewString())

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

func WithContext[K Event](app interfaces.App, handler func(logger *logrus.Entry, event *K) error, opts ...ContextOption) func(event *K) error {
	// Context builder
	var cb ContextBuilder
	cb.app = app

	for _, opt := range opts {
		opt(&cb)
	}

	funcName := getFuncClear(handler)

	return func(event *K) error {
		// Firstly, we extract the context
		// Try not to use to much reflection
		ctxVal := reflect.ValueOf(event).Elem().FieldByName("Ctx")
		if !ctxVal.IsValid() {
			return errors.New("invalid handler passed")
		}

		logger := logrus.
			WithField("func", funcName).
			WithField("user_id", (*event).User().ID)

		startTime := time.Now()
		defer func() {
			logger.
				WithField("duration", time.Since(startTime)).
				Info("user ran command")
		}()

		switch v := ctxVal.Interface().(type) {
		default:
			return errors.New("invalid handler passed")
		case context.Context:
			v, err := FillContextReply(&cb, *event, v)
			if err != nil {
				logger.
					WithError(err).
					Error("got error when filling context")
				return nil
			}

			ctxVal.Set(reflect.ValueOf(v))
		}

		return handler(logger, event)
	}
}

func WithContextD[D any, K Event](app interfaces.App, handler func(logger *logrus.Entry, data D, event *K) error, opts ...ContextOption) func(data D, event *K) error {
	// Context builder
	var cb ContextBuilder
	cb.app = app

	for _, opt := range opts {
		opt(&cb)
	}

	funcName := getFuncClear(handler)

	return func(data D, event *K) error {
		// Firstly, we extract the context
		// Try not to use to much reflection
		ctxVal := reflect.ValueOf(event).Elem().FieldByName("Ctx")
		if !ctxVal.IsValid() {
			return errors.New("invalid handler passed")
		}

		logger := logrus.
			WithField("func", funcName).
			WithField("user_id", (*event).User().ID)

		startTime := time.Now()
		defer func() {
			logger.
				WithField("duration", time.Since(startTime)).
				Info("command")
		}()

		switch v := ctxVal.Interface().(type) {
		default:
			return errors.New("invalid handler passed")
		case context.Context:
			v, err := FillContextReply(&cb, *event, v)
			if err != nil {
				logger.
					WithError(err).
					Error("got error when filling context")
				return nil
			}

			ctxVal.Set(reflect.ValueOf(v))
		}

		return handler(logger, data, event)
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

func WithPlayerAuctions() func(a *ContextBuilder) {
	return func(a *ContextBuilder) {
		a.withPlayerAuctions = true
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

func WithPlayerLimits() func(a *ContextBuilder) {
	return func(a *ContextBuilder) {
		a.withPlayerLimits = true
	}
}

func WithPlayerDungeons() func(a *ContextBuilder) {
	return func(a *ContextBuilder) {
		a.withPlayerDungeons = true
	}
}

func WithPlayerDaily() func(a *ContextBuilder) {
	return func(a *ContextBuilder) {
		a.withPlayerDaily = true
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
