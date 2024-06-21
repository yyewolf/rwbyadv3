package interfaces

import (
	"time"

	"github.com/yyewolf/rwbyadv3/models"
)

type JobKey string

type JobHandler interface {
	OnEvent(key JobKey, f func(params map[string]interface{}) error)
	SendEvent(key JobKey, jobID string, params map[string]interface{}) (*models.Job, error)
	ScheduleJob(key JobKey, jobID string, runAt time.Time, params map[string]interface{}) (*models.Job, error)
	ScheduleRecurringJob(key JobKey, runAt time.Time, every time.Duration) (*models.Job, error)
	CancelJob(key JobKey, jobID string) error

	Init() error
	Start() error
	Shutdown() error
}
