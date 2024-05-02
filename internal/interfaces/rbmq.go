package interfaces

import (
	"time"

	"github.com/yyewolf/rwbyadv3/models"
)

type JobKey string

type JobHandler interface {
	RegisterJobKey(key JobKey, f func(params map[string]interface{}) error)
	ScheduleJob(key JobKey, jobID string, runAt time.Time, params map[string]interface{}) (*models.Job, error)
	CancelJob(key JobKey, jobID string) error

	Init() error
	Start() error
	Shutdown() error
}
