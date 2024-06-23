package jobs

import (
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/values"
	"github.com/yyewolf/rwbyadv3/models"
)

type JobHandler struct {
	config *env.Config
	conn   *amqp.Connection
	ch     *amqp.Channel
	close  chan bool
	closed bool

	jobTypes        map[interfaces.JobKey]func(params map[string]interface{}) error
	reScheduleQueue []*models.Job
}

func New(options ...Option) interfaces.JobHandler {
	var jobHandler = &JobHandler{}

	for _, opt := range options {
		opt(jobHandler)
	}

	jobHandler.close = make(chan bool)
	jobHandler.jobTypes = make(map[interfaces.JobKey]func(params map[string]interface{}) error)

	return jobHandler
}

func (j *JobHandler) Start() error {
	conn, err := amqp.DialConfig(fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		j.config.Rbmq.User,
		j.config.Rbmq.Pass,
		j.config.Rbmq.Host,
		j.config.Rbmq.Port,
	), amqp.Config{
		Heartbeat: 10 * time.Second,
	})
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	j.conn = conn
	j.ch = ch

	err = j.Init()
	if err != nil {
		return err
	}

	go func() {
		for {
			j.Listen()
			if j.closed {
				return
			}
			time.Sleep(2 * time.Second)
		}
	}()

	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-j.close:
			j.closed = true
			return nil
		case <-ticker.C:
			if conn.IsClosed() {
				conn, err := amqp.DialConfig(
					fmt.Sprintf("amqp://%s:%s@%s:%s/",
						j.config.Rbmq.User,
						j.config.Rbmq.Pass,
						j.config.Rbmq.Host,
						j.config.Rbmq.Port,
					),
					amqp.Config{
						Heartbeat: 10 * time.Second,
					})
				if err != nil {
					return err
				}

				ch, err := j.conn.Channel()
				if err != nil {
					return err
				}

				j.conn = conn
				j.ch = ch
			} else {
				for i := 0; i < 5; i++ {
					if len(j.reScheduleQueue) == 0 {
						break
					}
					job := j.reScheduleQueue[0]
					err := j.reScheduleJob(job)
					if err == nil {
						j.reScheduleQueue = j.reScheduleQueue[1:]
						break
					}
					logrus.WithError(err).Error("failed to reschedule job")
					time.Sleep(200 * time.Millisecond)
				}
			}
		}
	}
}

func (j *JobHandler) WaitAvailable() {
	for {
		if j.conn != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func (j *JobHandler) Shutdown() error {
	select {
	case j.close <- true:
		return nil
	default:
		return values.ErrAppAlreadyClosed
	}
}
