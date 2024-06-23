package jobs

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/models"
)

func (j *JobHandler) SendEvent(key interfaces.JobKey, jobID string, params map[string]interface{}) (*models.Job, error) {
	var p types.JSON
	p.Marshal(params)

	job := &models.Job{
		ID:     jobID,
		Jobkey: string(key),
		RunAt:  time.Now().Add(100 * time.Millisecond),
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
		string(key),
		false,
		false,
		amqp091.Publishing{
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
