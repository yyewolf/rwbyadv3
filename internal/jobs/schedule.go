package jobs

import (
	"context"
	"database/sql"
	"encoding/json"
	"math"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/models"
)

func (j *JobHandler) ScheduleJob(key interfaces.JobKey, jobID string, runAt time.Time, params map[string]interface{}) (*models.Job, error) {
	var p types.JSON
	p.Marshal(params)

	job := &models.Job{
		ID:     jobID,
		Jobkey: string(key),
		RunAt:  runAt,
		Params: p,
	}

	tx, err := boil.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	err = job.Insert(context.Background(), tx, boil.Infer())
	if err != nil {
		return nil, err
	}

	bdy, err := json.Marshal(job)
	if err != nil {
		return nil, err
	}

	err = j.ch.PublishWithContext(
		context.Background(),
		j.config.Rbmq.Jobs.Exchange,
		"job",
		false,
		false,
		amqp091.Publishing{
			Headers: amqp091.Table{
				"x-delay": int(time.Until(runAt).Seconds() * 1000),
			},
			Body: bdy,
		},
	)
	if err != nil {
		j.reScheduleQueue = append(j.reScheduleQueue, job)
		tx.Rollback()
	}
	tx.Commit()

	return job, err
}

func setNextRun(job *models.Job) {
	deltaT := float64(time.Since(job.RunAt).Seconds())
	duration := float64(job.DeltaTime)
	amountOfTimeItShouldHaveRan := int64(math.Floor(deltaT/duration)) + 1

	job.RunAt = job.RunAt.Add(time.Duration(amountOfTimeItShouldHaveRan*job.DeltaTime) * time.Second)
}

func jobRunID(job *models.Job) int64 {
	deltaT := float64(time.Since(job.RunAt).Seconds())
	duration := float64(job.DeltaTime)
	amountOfTimeItShouldHaveRan := int64(math.Floor(deltaT/duration)) + 1

	return amountOfTimeItShouldHaveRan
}

func (j *JobHandler) reScheduleJob(job *models.Job) error {
	if !job.Recurring && !job.Errored {
		job.Retries++
		job.RunAt = job.RunAt.Add(10 * time.Duration(math.Pow(2, float64(job.Retries))) * time.Second)
	} else {
		job.Retries = 0

		setNextRun(job)
	}

	tx, err := boil.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = job.Update(context.Background(), tx, boil.Infer())
	if err != nil {
		return err
	}

	bdy, err := json.Marshal(job)
	if err != nil {
		return err
	}

	err = j.ch.PublishWithContext(
		context.Background(),
		j.config.Rbmq.Jobs.Exchange,
		"job",
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
		tx.Rollback()
	}
	tx.Commit()

	return err
}

func (j *JobHandler) ScheduleRecurringJob(key interfaces.JobKey, runAt time.Time, every time.Duration) (*models.Job, error) {
	var p types.JSON
	p.Marshal(make(map[string]interface{}))

	job := &models.Job{
		ID:        string(key),
		Jobkey:    string(key),
		RunAt:     runAt,
		DeltaTime: int64(every.Seconds()),
		Params:    p,
		Recurring: true,
	}

	tx, err := boil.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	exists, err := job.Exists(context.Background(), tx)
	if err != nil {
		return nil, err
	}

	setNextRun(job)
	if exists {
		_, err = job.Update(context.Background(), tx, boil.Infer())
		if err != nil {
			return nil, err
		}
	} else {
		err = job.Insert(context.Background(), tx, boil.Infer())
		if err != nil {
			return nil, err
		}
	}

	bdy, err := json.Marshal(job)
	if err != nil {
		return nil, err
	}

	if j.ch == nil {
		j.reScheduleQueue = append(j.reScheduleQueue, job)
		tx.Commit()
		return job, nil
	}

	err = j.ch.PublishWithContext(
		context.Background(),
		j.config.Rbmq.Jobs.Exchange,
		"job",
		false,
		false,
		amqp091.Publishing{
			Headers: amqp091.Table{
				"x-delay": int(time.Until(runAt).Seconds() * 1000),
			},
			Body: bdy,
		},
	)
	if err != nil {
		j.reScheduleQueue = append(j.reScheduleQueue, job)
		tx.Rollback()
	}
	tx.Commit()

	return job, err
}
