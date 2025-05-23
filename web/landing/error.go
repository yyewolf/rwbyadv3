package landing

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

// Predefined error types
const (
	ErrorNotFound        = "404"
	ErrorForbidden       = "403"
	ErrorUnauthorized    = "401"
	ErrorInternal        = "500"
	ErrorDungeonNotFound = "dungeon_not_found"
	ErrorAuthentication  = "auth_error"
	ErrorGithubStar      = "github_star_error"
	ErrorTrade           = "trade_error"
	ErrorMarket          = "market_error"
)

// ErrorPage redirects to the landing error page with the specified error code
func ErrorPage(c echo.Context, errorType string) error {
	return c.Redirect(http.StatusFound, fmt.Sprintf("/landing/error/?error=%s", errorType))
}

// RegisterErrorHandler registers the error page routes
func RegisterErrorHandler(app interfaces.App, g *echo.Group) {
	// No specific routes needed for the error page as it will be handled by SvelteKit
}
