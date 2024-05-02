package rbmq

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/models"
)

func (j *JobHandler) RegisterJobKey(key interfaces.JobKey, f func(params map[string]interface{}) error) {
	_, found := j.jobTypes[key]
	if found {
		logrus.Fatal("a job with this key already exists")
	}

	j.jobTypes[key] = f
}

func (j *JobHandler) ScheduleJob(key interfaces.JobKey, jobID string, runAt time.Time, params map[string]interface{}) (*models.Job, error) {
	var p types.JSON
	p.Marshal(params)

	job := &models.Job{
		ID:     jobID,
		Jobkey: string(key),
		RunAt:  runAt,
		Params: p,
	}

	err := job.InsertG(context.Background(), boil.Infer())
	if err != nil {
		return nil, err
	}

	bdy, err := json.Marshal(job)
	if err != nil {
		return nil, err
	}

	j.ch.PublishWithContext(
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

	return job, err
}

func (j *JobHandler) reScheduleJob(job *models.Job) error {
	job.Retries++
	job.RunAt = job.RunAt.Add(10 * time.Duration(math.Pow(2, float64(job.Retries))) * time.Second)

	_, err := job.UpdateG(context.Background(), boil.Infer())
	if err != nil {
		return err
	}

	bdy, err := json.Marshal(job)
	if err != nil {
		return err
	}

	j.ch.PublishWithContext(
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

	return err
}

func (j *JobHandler) CancelJob(key interfaces.JobKey, jobID string) error {
	job, err := models.FindJobG(context.Background(), jobID, string(key))
	if err != nil {
		return err
	}

	_, err = job.DeleteG(context.Background(), false)
	if err != nil {
		return err
	}

	return nil
}

func (j *JobHandler) handleJob(job *models.Job) {
	exists, err := models.JobExistsG(context.Background(), job.ID, job.Jobkey)
	if err != nil {
		logrus.Error("error checking existance")
		j.reScheduleJob(job)
		return
	}

	if !exists {
		logrus.Error("job canceled")
		return
	}

	f, found := j.jobTypes[interfaces.JobKey(job.Jobkey)]
	if !found {
		logrus.Error("job type not found")
		j.reScheduleJob(job)
		return
	}

	params := make(map[string]interface{})
	err = job.Params.Unmarshal(&params)
	if err != nil {
		logrus.Error(err)
		j.reScheduleJob(job)
		return
	}

	err = f(params)
	if err != nil {
		logrus.Error(err)
		j.reScheduleJob(job)
		return
	}

	// delete job
	_, err = job.DeleteG(context.Background(), false)
	if err != nil {
		logrus.Error(err)
	}
}
