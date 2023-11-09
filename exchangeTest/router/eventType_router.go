package router

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Функция для инициализации маршрутизатора событий по типам
func initTypeRouter() eventRoutData {
	eventCh := make(EventChan, 100)
	subscrCh := make(SubscriberChan, 100)
	done := make(chan struct{})
	cancel := func() {
		close(done)
	}

	go typeRouter(eventCh, subscrCh, done)
	log.Debug("Запущен маршрутизатор типов событий")

	return eventRoutData{
		eventCh:  eventCh,
		subscrCh: subscrCh,
		cancel:   cancel,
	}
}

// Маршрутизатор сообщений по типу
func typeRouter(eventCh EventChan, subscrCh SubscriberChan, done chan struct{}) {
	defer log.Debugf("Работа маршрутизатора типов событий завершена")

	types := make(map[string]eventRoutData)

	// добавление маршрутизатора на все события
	allEventRouter := initTypeEventRouter("allEvent")
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
				routData := initTypeEventRouter(event.TypeEvent)
				types[event.TypeEvent] = routData
				routData.eventCh <- event
			}
		case sub := <-subscrCh:
			// отправка подписчика маршрутизатору всех типов
			if sub.AllEvent {
				allEventRouter.subscrCh <- sub
			}

			for _, eventType := range sub.Types {
				eventData, ok := types[eventType]
				if ok {
					eventData.subscrCh <- sub
				} else {
					log.Warningf("Подписчик %s не может подписаться на тип события %s, этот тип события не существует", sub.Name, eventType)
					go func(subMess SubscriberMess) {
						time.Sleep(5 * time.Second)
						log.Debugf("Повторная попытка подписать %s на событие %s", subMess.Name, subMess.Types[0])
						subscrCh <- subMess
						return
					}(SubscriberMess{
						Name:   sub.Name,
						Types:  []string{eventType},
						EvenCh: sub.EvenCh,
					})
				}

			}
		case <-done:
			for _, routData := range types {
				routData.cancel()
			}
		}
	}
}

// Запускает в отдельной горутине typeEventRouter (маршрутизатор типа)
func initTypeEventRouter(eventType string) eventRoutData {
	done := make(chan struct{})

	routData := eventRoutData{
		eventCh:  make(EventChan, 100),
		subscrCh: make(SubscriberChan, 100),
		cancel: func() {
			close(done)
		},
	}

	go typeEventRouter(routData.eventCh, routData.subscrCh, done, eventType)

	log.Debugf("typeEventRouter для типа %s запущен", eventType)

	return routData
}

// Маршрутизатор конкретного типа события.
// Содержит словарь со всеми подписчиками.
// При получении события отправляет его всем подписчикам.
// При получении подписчика, сохраняет его в свой словарь.
func typeEventRouter(eventCh EventChan, subscrCh SubscriberChan, done chan struct{}, eventType string) {
	defer log.Debugf("Работа маршрутизатора типа %s завершена", eventType)

	subscribers := make(map[string]EventChan)

	for {
		select {
		case event := <-eventCh:
			// проверка типа события (ошибки быть не дожно, проверка на всякий случай)
			// if event.typeEvent != eventType {
			// 	log.Errorf("Ожидается тип события %s, получен %s", eventType, event.typeEvent)
			// 	continue
			// }

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
		case <-done:
			return
		}
	}
}
