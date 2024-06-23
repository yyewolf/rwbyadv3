package jobs

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/models"
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
	logrus.WithField("job_key", job.Jobkey).Debug("handling job")
	job.Errored = true
	exists, err := models.JobExistsG(context.Background(), job.ID, job.Jobkey)
	if err != nil {
		logrus.WithField("job_key", job.Jobkey).Error("error checking existance")
		j.reScheduleQueue = append(j.reScheduleQueue, job)
		return
	}

	if !exists {
		logrus.WithField("job_key", job.Jobkey).Error("job canceled")
		return
	}

	f, found := j.jobTypes[interfaces.JobKey(job.Jobkey)]
	if !found {
		logrus.WithField("job_key", job.Jobkey).Error("job type not found, deleting")
		job.DeleteG(context.Background(), false)
		return
	}

	oldID := job.LastRunID
	tx, err := boil.BeginTx(context.Background(), nil)
	if err != nil {
		logrus.WithField("job_key", job.Jobkey).Error("couldn't start tx")
		j.reScheduleQueue = append(j.reScheduleQueue, job)
		tx.Rollback()
		return
	}

	savedJob, err := models.FindJob(context.Background(), tx, job.ID, job.Jobkey)
	if err != nil {
		logrus.WithField("job_key", job.Jobkey).Error("couldn't get job from db")
		j.reScheduleQueue = append(j.reScheduleQueue, job)
		tx.Rollback()
		return
	}

	if savedJob.DeltaTime != job.DeltaTime {
		logrus.WithField("job_key", job.Jobkey).Error("job delta time changed")
		tx.Rollback()
		return
	}

	if savedJob.LastRunID != oldID {
		logrus.WithField("job_key", job.Jobkey).Error("job already ran")
		tx.Rollback()
		return
	}

	savedJob.LastRunID = jobRunID(job)
	savedJob.Update(context.Background(), tx, boil.Infer())
	tx.Commit()

	params := make(map[string]interface{})
	err = job.Params.Unmarshal(&params)
	if err != nil {
		logrus.WithField("job_key", job.Jobkey).Error(err)
		j.reScheduleQueue = append(j.reScheduleQueue, job)
		return
	}

	err = f(params)
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
		_, err = job.DeleteG(context.Background(), false)
		if err != nil {
			logrus.Error(err)
		}
	}
}
