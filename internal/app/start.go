package app

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/internal/values"
)

func (a *App) Start() {
	go func() {
		err := a.state.Connect(context.TODO())

		select {
		case a.errorChannel <- err:
		default:
		}
	}()

	select {
	case <-a.shutdown:
		if a.state != nil {
			a.state.Close()
		}

		if a.db != nil {
			a.db.Disconnect()
		}
	case err := <-a.errorChannel:
		logrus.Error(err)
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
