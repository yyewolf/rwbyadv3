package metrics

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/models"
)

var (
	CommandsProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rwby_processed_commands_total",
			Help: "The total number of processed commands",
		},
		[]string{},
	)
)

func AliveCardsGauge(a interfaces.App) func() float64 {
	return func() float64 {
		amount, err := models.Cards().CountG(context.Background())
		if err != nil {
			return 0
		}
		return float64(amount)
	}
}

func DeadCardsGauge(a interfaces.App) func() float64 {
	return func() float64 {
		query := models.NewQuery(
			qm.From("cards"),
			qm.Select("COUNT(*)"),
			qm.Where("deleted_at is not null"),
		)

		var r struct {
			Count int64
		}

		err := query.BindG(context.Background(), &r)
		if err != nil {
			return -1.0
		}
		return float64(r.Count)
	}
}

func NewMetricsHandler(app interfaces.App, g *echo.Group) {
	g.Any("/", echo.WrapHandler(promhttp.Handler()))

	// Register funcs
	promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "rwby_cards_alive_total",
			Help: "The total number of alive cards",
		},
		AliveCardsGauge(app),
	)

	promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "rwby_cards_dead_total",
			Help: "The total number of dead cards",
		},
		DeadCardsGauge(app),
	)
}
