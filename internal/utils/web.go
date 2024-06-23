package utils

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/models"
)

func GetRedirectForW(w string) (redirectUri string, err error) {
	switch w {
	case "main":
		redirectUri = "/"
	default:
		return "", errors.New("invalid redirect")
	}

	return redirectUri, nil
}

func GetSessionFromContext(c echo.Context) *models.AuthCookie {
	return c.Get("session").(*models.AuthCookie)
}
