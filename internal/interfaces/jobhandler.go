package interfaces

import (
	"github.com/yyewolf/rwbyadv3/models"
)

type JobKey string

type JobHandler interface {
	OnEvent(key JobKey, f func(params map[string]interface{}) error)
	SendEvent(key JobKey, jobID string, params map[string]interface{}) (*models.Job, error)
	CancelJob(key JobKey, jobID string) error

	Init() error
	Start() error
	Shutdown() error

	WaitAvailable()
}
