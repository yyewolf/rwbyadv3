package app

import (
	"encoding/json"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

func (app *App) SendDMJob(params map[string]interface{}) error {
	// get user id from param "user_id"
	id := params["user_id"].(string)

	var message discord.MessageCreate
	b, _ := json.Marshal(params["message"])
	json.Unmarshal(b, &message)

	c := app.Client()

	ch, err := c.Rest().CreateDMChannel(snowflake.MustParse(id))
	if err != nil {
		return err
	}

	_, err = c.Rest().CreateMessage(ch.ID(), message)
	return err
}
