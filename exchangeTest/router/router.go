package router

import (
	log "github.com/sirupsen/logrus"
)

// Модуль в котором производится маршрутизация событий по получателям
// Задачи:

// TODO: Создание канала для получения событий
// TODO: Создание канала для получения подписчиков

// TODO: Функукция для маршрутизации типов
// TODO: Функция для маршрутизации по отправителям
// TODO: Функция для маршрутизации всех событий (для подписчиков, которым нужны все события)
// TODO: Маршрутизатор событий (принимает входящие каналы(для событий) от маршрутизаторов: типа, отправителя, по всем)
// TODO: Маршрутизатор подписчиков (принимает входящие каналы подписчиков от маршрутизаторов)

// Функция для запуска маршрутизатора
func InitRouter() (EventChan, SubscriberChan) {
	var routersCh routersChans
	eventCh, _ := initEventRouter(routersCh)
	SubscrCh := make(SubscriberChan, 100)

	return eventCh, SubscrCh
}

// Инициализация всех маршрутизаторов (TypeRouter, PublishRouter, EventRouter)
func initRouters() {

}

// Функция для инициализации маршрутизатора событий по типам
func initTypeRouter() (EventChan, SubscriberChan, func()) {
	eventCh := make(EventChan, 100)
	subscrCh := make(SubscriberChan, 100)
	done := make(chan struct{})
	cancel := func() {
		close(done)
	}

	return eventCh, subscrCh, cancel
}

// Маршрутизатор сообщений по типу
func typeRouter(eventCh EventChan, subscrCh SubscriberChan, done chan struct{}) {
	types := make(map[string]eventRoutKIT)

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

		case <-done:
			return
		}
	}
}

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
			if !ok {
				subscribers[sub.name] = sub.evenCh
			}
		case <-done:
			return
		}
	}
}

// Функция для инициализации маршрутизатора событий по отправителям
func initPublishRouter() (EventChan, SubscriberChan, func()) {
	eventCh := make(EventChan, 100)
	subscrCh := make(SubscriberChan, 100)
	done := make(chan struct{})
	cancel := func() {
		close(done)
	}

	return eventCh, subscrCh, cancel
}

// Функция для инициализации маршрутизатора всех событий
func initAllEventRouter() (EventChan, SubscriberChan, func()) {
	eventCh := make(EventChan, 100)
	subscrCh := make(SubscriberChan, 100)
	done := make(chan struct{})
	cancel := func() {
		close(done)
	}

	return eventCh, subscrCh, cancel
}

// Функция инициализации маршрутизатора событий
func initEventRouter(routersCh routersChans) (EventChan, func()) {
	eventCh := make(EventChan, 100)
	done := make(chan struct{})
	cancel := func() {
		close(done)
	}
	go eventRouter(routersCh, eventCh, done)
	log.Debug("eventRouter запущен")
	return eventCh, cancel
}

// Функция маршрутизации событий
func eventRouter(routersCh routersChans, eventCh EventChan, done chan struct{}) {
	defer log.Debug("Работа маршрутизатора событий завершена")

	for {
		select {
		case event := <-eventCh:
			routersCh.typeCh <- event
			routersCh.publishCh <- event
			routersCh.allCh <- event
		case <-done:
			return
		}
	}
}
