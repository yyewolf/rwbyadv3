package app

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/internal/values"
)

func (a *App) Start() {
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
	case err := <-a.errorChannel:
		logrus.WithField("error", err).Error("An error stopped execution")
	}
}

func (a *App) Shutdown() error {
	select {
	case a.shutdown <- struct{}{}:
		return nil
	default:
		return values.ErrAppAlreadyClosed
	}
}
