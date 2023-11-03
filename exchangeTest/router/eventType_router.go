package router

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Функция для инициализации маршрутизатора событий по типам
func initTypeRouter() (EventChan, SubscriberChan, func()) {
	eventCh := make(EventChan, 100)
	subscrCh := make(SubscriberChan, 100)
	done := make(chan struct{})
	cancel := func() {
		close(done)
	}

	go typeRouter(eventCh, subscrCh, done)
	log.Debug("Запущен маршрутизатор типов событий")

	return eventCh, subscrCh, cancel
}

// Маршрутизатор сообщений по типу
func typeRouter(eventCh EventChan, subscrCh SubscriberChan, done chan struct{}) {
	types := make(map[string]eventRoutKIT)
	defer log.Debugf("Работа маршрутизатора типов событий завершена")

	for {
		select {
		case event := <-eventCh:
			routKIT, ok := types[event.typeEvent]
			if !ok {
				routKIT := initTypeEventRouter(event.typeEvent)
				types[event.typeEvent] = routKIT
				routKIT.eventCh <- event
			}
			routKIT.eventCh <- event

		case sub := <-subscrCh:
			// TODO: нужно сделать отправку подписчика определенным типам (которые указаны в сообщении)
			for _, eventType := range sub.types {
				eventKIT, ok := types[eventType]
				if ok {
					eventKIT.subscrCh <- sub
				} else {
					log.Warningf("Подписчик %s не может подписаться на тип события %s, этот тип события не существует", sub.name, eventType)
					go func(subMess SubscriberMess) {
						time.Sleep(5 * time.Second)
						log.Debugf("Повторная попытка подписать %s на событие %s", subMess.name, subMess.types[0])
						subscrCh <- subMess
						return
					}(SubscriberMess{
						name:   sub.name,
						types:  []string{eventType},
						evenCh: sub.evenCh,
					})
				}

			}
		case <-done:
			return
		}
	}
}

// Запускает в отдельной горутине typeEventRouter (маршрутизатор типа)
func initTypeEventRouter(eventType string) eventRoutKIT {
	done := make(chan struct{})

	routKIT := eventRoutKIT{
		eventCh:  make(EventChan, 100),
		subscrCh: make(SubscriberChan, 100),
		cancel: func() {
			close(done)
		},
	}

	go typeEventRouter(routKIT.eventCh, routKIT.subscrCh, done, eventType)

	log.Debugf("typeEventRouter для типа %s запущен", eventType)

	return routKIT
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
			if event.typeEvent != eventType {
				log.Errorf("Ожидается тип события %s, получен %s", eventType, event.typeEvent)
				continue
			}

			// TODO: можно добавить отписку подписчика от данного типа событий
			for subscr, ch := range subscribers {
				ch <- event
				log.Debugf("Событие типа %s отправлено подписчику %s", event.typeEvent, subscr)
			}
		case sub := <-subscrCh:
			_, ok := subscribers[sub.name]
			if ok {
				log.Warningf("%s уже подписан на тип событий %s", sub.name, eventType)
			} else {
				subscribers[sub.name] = sub.evenCh
			}
		case <-done:
			return
		}
	}
}
