package router

import (
	"context"
	"time"

	models "github.com/SouthUral/exchangeTest/models"
	log "github.com/sirupsen/logrus"
)

// Функция для инициализации маршрутизатора событий по типам
func initTypeRouter(ctx context.Context) eventRoutData {
	eventCh := make(models.EventChan, 100)
	subscrCh := make(chan models.SubscriberMess, 100)

	go typeRouter(ctx, eventCh, subscrCh)
	log.Debug("Запущен маршрутизатор типов событий")

	return eventRoutData{
		eventCh:  eventCh,
		subscrCh: subscrCh,
	}
}

// Маршрутизатор сообщений по типу
func typeRouter(ctx context.Context, eventCh models.EventChan, subscrCh chan models.SubscriberMess) {
	defer log.Debugf("Работа маршрутизатора типов событий завершена")

	types := make(map[string]eventRoutData)

	// добавление маршрутизатора на все события
	allEventRouter := initTypeEventRouter(ctx, "allEvent")
	types["allEvent"] = allEventRouter

	for {
		select {
		case event := <-eventCh:
			// отправка события маршрутизатору всех типов событий
			allEventRouter.eventCh <- event

			routData, ok := types[event.TypeEvent]
			if ok {
				routData.eventCh <- event
			} else {
				routData := initTypeEventRouter(ctx, event.TypeEvent)
				types[event.TypeEvent] = routData
				routData.eventCh <- event
			}
		case subMess := <-subscrCh:
			// отправка подписчика маршрутизатору всех типов
			subConf := subMess.ConfSubscribe
			if subConf.AllEvent {
				allEventRouter.subscrCh <- subMess
			}

			for _, eventType := range subConf.Types {
				eventData, ok := types[eventType]
				if ok {
					eventData.subscrCh <- subMess
				} else {
					log.Warningf("Подписчик %s не может подписаться на тип события %s, этот тип события не существует", subMess.Name, eventType)
					go func(subMess models.SubscriberMess) {
						time.Sleep(5 * time.Second)
						log.Debugf("Повторная попытка подписать %s на событие %s", subMess.Name, subConf.Types[0])
						subscrCh <- subMess
						return
					}(models.SubscriberMess{
						Name: subMess.Name,
						ConfSubscribe: models.ConfSub{
							Types: []string{eventType},
						},
						EvenCh: subMess.EvenCh,
					})
				}

			}
		case <-ctx.Done():
			return
		}
	}
}

// Запускает в отдельной горутине typeEventRouter (маршрутизатор типа)
func initTypeEventRouter(ctx context.Context, eventType string) eventRoutData {

	routData := eventRoutData{
		eventCh:  make(models.EventChan, 100),
		subscrCh: make(chan models.SubscriberMess, 100),
	}

	go typeEventRouter(ctx, routData.eventCh, routData.subscrCh, eventType)

	log.Debugf("typeEventRouter для типа %s запущен", eventType)

	return routData
}

// Маршрутизатор конкретного типа события.
// Содержит словарь со всеми подписчиками.
// При получении события отправляет его всем подписчикам.
// При получении подписчика, сохраняет его в свой словарь.
func typeEventRouter(ctx context.Context, eventCh models.EventChan, subscrCh chan models.SubscriberMess, eventType string) {
	defer log.Debugf("Работа маршрутизатора типа %s завершена", eventType)

	subscribers := make(map[string]models.EventChan)

	for {
		select {
		case event := <-eventCh:
			// TODO: можно добавить отписку подписчика от данного типа событий
			for subscr, ch := range subscribers {
				ch <- event
				log.Debugf("Событие типа %s отправлено подписчику %s", event.TypeEvent, subscr)
			}
		case sub := <-subscrCh:
			_, ok := subscribers[sub.Name]
			if ok {
				log.Warningf("%s уже подписан на тип событий %s", sub.Name, eventType)
			} else {
				subscribers[sub.Name] = sub.EvenCh
			}
		case <-ctx.Done():
			return
		}
	}
}
