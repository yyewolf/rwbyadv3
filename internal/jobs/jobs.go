package jobs

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

func (j *JobHandler) OnEvent(key interfaces.JobKey, f func(params map[string]interface{}) error) {
	_, found := j.jobTypes[key]
	if found {
		logrus.Fatal("a job with this key already exists")
	}

	j.jobTypes[key] = f

	if j.ch != nil {
		j.ch.QueueBind(
			j.config.Rbmq.Jobs.Queue,
			string(key),
			j.config.Rbmq.Jobs.Exchange,
			false,
			nil,
		)
	}
}

func (j *JobHandler) CancelJob(key interfaces.JobKey, jobID string) error {
	return j.entClient.Job.DeleteOneID(uuid.MustParse(jobID)).Exec(context.Background())
}

func (j *JobHandler) handleJob(job *ent.Job) {
	logrus.WithField("job_key", job.Jobkey).WithField("params", job.Params).Info("job started")
	job.Errored = true

	savedJob, err := j.entClient.Job.Get(context.Background(), job.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			logrus.WithField("job_key", job.Jobkey).Error("job not found, deleting")
		} else {
			logrus.WithField("job_key", job.Jobkey).Error("error checking existance")
			j.reScheduleQueue = append(j.reScheduleQueue, job)
		}
		return
	}

	f, found := j.jobTypes[interfaces.JobKey(job.Jobkey)]
	if !found {
		logrus.WithField("job_key", job.Jobkey).Error("job type not found, deleting")
		j.entClient.Job.DeleteOne(job).Exec(context.Background())
		return
	}

	oldID := job.LastRunID

	ctx := context.TODO()

	err = ent.WithTx(ctx, j.entClient, func(tx *ent.Tx) error {
		if savedJob.DeltaTime != job.DeltaTime {
			return fmt.Errorf("job delta time changed")
		}

		if savedJob.LastRunID != oldID {
			return fmt.Errorf("job already ran")
		}

		err = tx.Job.UpdateOne(savedJob).
			SetLastRunID(jobRunID(savedJob)).
			Exec(ctx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logrus.WithField("job_key", job.Jobkey).Error(err)
		return
	}

	err = f(job.Params)
	if err != nil {
		logrus.WithField("job_key", job.Jobkey).WithField("job_id", job.ID).Error(err)
		j.reScheduleQueue = append(j.reScheduleQueue, job)
		return
	}

	// Job has not errored
	job.Errored = false

	if job.Recurring {
		j.reScheduleQueue = append(j.reScheduleQueue, job)
	} else {
		// delete job
		err = j.entClient.Job.DeleteOne(job).Exec(context.Background())
		if err != nil {
			logrus.Error(err)
		}
	}
}
