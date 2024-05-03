package jobs

import (
	"github.com/yyewolf/rwbyadv3/internal/env"
)

type Option func(j *JobHandler)

func WithConfig(config *env.Config) Option {
	return func(j *JobHandler) {
		j.config = config
	}
}
