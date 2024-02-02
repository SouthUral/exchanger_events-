package utilsfortest

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

// обновляемый таймер, вызывает функцию cancel если долго не получает сообщение по каналу
func UpdatableTimer(cancel func(), timeWait int) chan struct{} {
	signalCh := make(chan struct{})

	go func() {
		defer cancel()
		defer log.Warning("Работа таймера закончена")

		duration := time.Duration(timeWait) * time.Second

		ctx, _ := context.WithTimeout(context.Background(), duration)

		for {
			select {
			case <-signalCh:
				log.Debug("Обновлен таймер")
				ctx, _ = context.WithTimeout(context.Background(), duration)
			case <-ctx.Done():
				return
			}
		}
	}()

	return signalCh

}
