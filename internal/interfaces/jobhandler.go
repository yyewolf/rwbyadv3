package interfaces

import "github.com/yyewolf/rwbyadv3/ent"

type JobKey string

type JobHandler interface {
	OnEvent(key JobKey, f func(params map[string]interface{}) error)
	SendEvent(key JobKey, params map[string]interface{}) (*ent.Job, error)
	CancelJob(key JobKey, jobID string) error

	Init() error
	Start() error
	Shutdown() error

	WaitAvailable()
}
