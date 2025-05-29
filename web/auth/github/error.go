package github

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/web/landing"
)

// ErrorPage redirects to the landing error page with the appropriate error type
func ErrorPage(c echo.Context, code int) error {
	var errorType string

	switch code {
	case 0: // Special case for GitHub star error
		errorType = landing.ErrorGithubStar
	case http.StatusNotFound:
		errorType = landing.ErrorNotFound
	case http.StatusForbidden:
		errorType = landing.ErrorForbidden
	case http.StatusUnauthorized:
		errorType = landing.ErrorUnauthorized
	case http.StatusInternalServerError:
		errorType = landing.ErrorInternal
	default:
		// GitHub-specific errors default to GitHub star error
		errorType = landing.ErrorGithubStar
	}

	return landing.ErrorPage(c, errorType)
}
