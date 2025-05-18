package metrics

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
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
		amount, err := a.Db().Card.Query().Count(context.Background())
		if err != nil {
			return 0
		}
		return float64(amount)
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
}
