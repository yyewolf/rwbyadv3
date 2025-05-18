package jobs

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

func (j *JobHandler) SendEvent(key interfaces.JobKey, params map[string]interface{}) (*ent.Job, error) {
	var job *ent.Job
	var err error

	ctx := context.Background()

	err = ent.WithTx(ctx, j.entClient, func(tx *ent.Tx) error {
		job, err = tx.Job.Create().
			SetJobkey(string(key)).
			SetRunAt(time.Now().Add(100 * time.Millisecond)).
			SetParams(params).
			Save(ctx)

		bdy, err := json.Marshal(job)
		if err != nil {
			return err
		}

		err = j.ch.PublishWithContext(
			context.Background(),
			j.config.Rbmq.Jobs.Exchange,
			string(key),
			false,
			false,
			amqp091.Publishing{
				Body: bdy,
			},
		)
		if err != nil {
			j.reScheduleQueue = append(j.reScheduleQueue, job)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return job, err
}
