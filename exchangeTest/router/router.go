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
