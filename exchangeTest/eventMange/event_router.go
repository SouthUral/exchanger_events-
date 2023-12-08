package eventmange

import (
	apiV1 "github.com/SouthUral/exchangeTest/api/api_v1"
	models "github.com/SouthUral/exchangeTest/models"

	log "github.com/sirupsen/logrus"
)

// Внешняя функция, должна принять все каналы и запустить в горутине eventManager
func InitEventManager(
	eventRtCh models.EventChan, // канал для отправки событий в роутер
	eventAPICh apiV1.EventAPICh, // канал для приема событий от API gRPC
	lastID int, // id из БД
) func() {
	done := make(chan struct{})
	cancel := func() {
		close(done)
	}

	go eventManager(lastID, eventRtCh, eventAPICh, done)

	log.Debug("eventManager is starting")

	return cancel
}

// Распределитель событий (на будущее, если нужно будет сохранять плоский список событий).
// Принимает событие от API gRPC выдет ему Id и отправляет обратно в API Id события.
// Отправляет событие с Id в маршрутизатор
// Отправляет событие с Id на сохранение в БД (пока непонятно в какую БД)
func eventManager(lastID int, eventRtCh models.EventChan, eventAPICh apiV1.EventAPICh, done <-chan struct{}) {
	// TODO: в будущем при запуске сервиса, нужно будет забрать последнее Id
	defer log.Warning("eventManager is finished")

	currentId := lastID

	for {
		select {
		case mess := <-eventAPICh:
			currentId++
			event := mess.Event
			event.Id = currentId

			// отправка события в маршрутизатор
			eventRtCh <- event

			mess.RevСh <- apiV1.ReverseMess{
				IdEvent: currentId,
				// пока Err nil, в будущем можно убрать
				Err: nil,
			}

		case <-done:
			return
		}
	}
}
