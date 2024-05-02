package rbmq

import (
	"context"
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/models"
)

func (j *JobHandler) Listen() error {
	channel, err := j.ch.ConsumeWithContext(
		context.Background(),
		j.config.Rbmq.Jobs.Queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for msg := range channel {
		var job models.Job
		d := msg.Body
		err := json.Unmarshal(d, &job)
		if err != nil {
			logrus.Error(err)
			msg.Ack(false)
		} else {
			j.handleJob(&job)
			msg.Ack(false)
		}
	}

	return nil
}
