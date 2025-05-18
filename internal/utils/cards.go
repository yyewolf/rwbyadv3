package utils

import (
	"net/url"

	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/internal/cards"
)

type Card struct{}

var Cards Card

func (card Card) IconURI(c *ent.Card) string {
	uri, _ := url.JoinPath(Players.c.App.BaseURI, cards.MustGetImageURI(c.CardType, "icon", "webp"))
	return uri
}
