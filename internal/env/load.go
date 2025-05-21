package env

import (
	logrusloki "github.com/schoentoon/logrus-loki"
	"github.com/yyewolf/rwbyadv3/internal/values"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var cfg Config

func Load() {
	godotenv.Load()

	if err := env.Parse(&cfg); err != nil {
		logrus.Fatalf("failed to load env: %v", err)
	}

	switch cfg.Mode {
	case values.Dev:
		logrus.SetLevel(logrus.DebugLevel)
	case values.Preprod:
		logrus.SetLevel(logrus.DebugLevel)
	case values.Prod:
		logrus.SetLevel(logrus.InfoLevel)
	case values.Unset:
		logrus.Fatalf("MODE is not set, be sure to have a .env file or set the environment variables")
	default:
		logrus.Fatalf("MODE is not set, be sure to have a .env file or set the environment variables")
	}

	if cfg.Loki.Enabled {
		hook, err := logrusloki.NewLokiDefaults(cfg.Loki.URI)
		if err != nil {
			logrus.Fatalf("failed to create loki hook: %v", err)
		}

		logrus.AddHook(hook)
	}

	logrus.Infof("Environment loaded: %s", cfg.Mode)
	logrus.Debugf("Environment: %+v", cfg)
}
