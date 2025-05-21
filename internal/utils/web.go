package utils

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/ent"
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

func GetSessionFromContext(c echo.Context) *ent.Cookie {
	return c.Get("session").(*ent.Cookie)
}
