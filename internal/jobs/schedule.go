package jobs

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/yyewolf/rwbyadv3/ent"
)

func setNextRun(job *ent.Job) {
	deltaT := float64(time.Since(job.RunAt).Seconds())
	duration := float64(job.DeltaTime)
	amountOfTimeItShouldHaveRan := int64(math.Floor(deltaT/duration)) + 1

	job.RunAt = job.RunAt.Add(time.Duration(amountOfTimeItShouldHaveRan*job.DeltaTime) * time.Second)
}

func jobRunID(job *ent.Job) int64 {
	deltaT := float64(time.Since(job.RunAt).Seconds())
	duration := float64(job.DeltaTime)
	amountOfTimeItShouldHaveRan := int64(math.Floor(deltaT/duration)) + 1

	if amountOfTimeItShouldHaveRan < 0 {
		amountOfTimeItShouldHaveRan = 0
	}

	return amountOfTimeItShouldHaveRan
}

func (j *JobHandler) reScheduleJob(job *ent.Job) error {
	if !job.Recurring && !job.Errored {
		job.Retries++
		job.RunAt = job.RunAt.Add(10 * time.Duration(math.Pow(2, float64(job.Retries))) * time.Second)
	} else {
		job.Retries = 0

		setNextRun(job)
	}

	ctx := context.Background()

	return ent.WithTx(ctx, j.entClient, func(tx *ent.Tx) error {
		job, err := tx.Job.UpdateOne(job).
			SetRetries(job.Retries).
			SetRunAt(job.RunAt).
			Save(ctx)

		bdy, err := json.Marshal(job)
		if err != nil {
			return err
		}

		err = j.ch.PublishWithContext(
			context.Background(),
			j.config.Rbmq.Jobs.Exchange,
			job.Jobkey,
			false,
			false,
			amqp091.Publishing{
				Headers: amqp091.Table{
					"x-delay": int(time.Until(job.RunAt).Seconds() * 1000),
				},
				Body: bdy,
			},
		)
		if err != nil {
			return err
		}

		return nil
	})
}
