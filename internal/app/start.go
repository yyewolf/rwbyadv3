package app

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/internal/values"
	"go.temporal.io/sdk/worker"
)

func (a *App) Start() {
	// Begin job handler here
	go func() {
		for {
			a.jobHandler.Start()

			err := a.jobHandler.Init()
			if err != nil {
				logrus.WithField("error", err).Error("Failed to initialize job handler")
				return
			}
		}
	}()

	go func() {
		err := a.client.OpenGateway(context.TODO())
		if err == nil {
			return
		}

		select {
		case a.errorChannel <- err:
		default:
		}
	}()

	go func() {
		err := a.temporalWorker.Run(worker.InterruptCh())
		if err == nil {
			return
		}

		select {
		case a.errorChannel <- err:
		default:
		}
	}()

	if a.enableWeb {
		go func() {
			err := a.webApp.Start()
			if err == nil {
				return
			}

			select {
			case a.errorChannel <- err:
			default:
			}
		}()
	}

	select {
	case <-a.shutdown:
		if a.client != nil {
			a.client.Close(context.TODO())
		}
		if a.enableWeb && a.webApp != nil {
			a.webApp.Stop()
		}

		a.temporalWorker.Stop()
		a.temporalClient.Close()
		a.jobHandler.Shutdown()
	case err := <-a.errorChannel:
		logrus.WithField("error", err).Error("An error stopped execution")
	}

	time.Sleep(2 * time.Second)
}

func (a *App) Shutdown() error {
	select {
	case a.shutdown <- struct{}{}:
		return nil
	default:
		return values.ErrAppAlreadyClosed
	}
}
