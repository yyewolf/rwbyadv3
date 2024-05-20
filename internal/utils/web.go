package utils

import (
	"errors"
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
