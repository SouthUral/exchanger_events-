package router

import (
	models "github.com/SouthUral/exchangeTest/models"
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
func InitRouter() (models.EventChan, chan models.SubscriberMess, func()) {
	// TODO: нужно понять какой буфер делать у каналов

	typeRoutData := initTypeRouter()
	publishRoutData := initPublishRouter()

	eventCh, eventRoutCancel := initEventRouter(typeRoutData.eventCh, publishRoutData.eventCh)
	subCh, subRoutCancel := initSubscribeRouter(typeRoutData.subscrCh, publishRoutData.subscrCh)

	cancel := func() {
		defer log.Debug("все маршрутизаторы выключены")
		typeRoutData.cancel()
		publishRoutData.cancel()
		eventRoutCancel()
		subRoutCancel()
	}

	return eventCh, subCh, cancel
}

// Инициализатор маршрутизатора событий
func initEventRouter(typeEventRoutCh, publishEventRoutCh models.EventChan) (models.EventChan, func()) {
	eventCh := make(models.EventChan, 100)
	done := make(chan struct{})
	cancel := func() {
		close(done)
	}

	go func() {
		defer log.Debug("Работа маршрутизатора событий завершена")
		for {
			select {
			case event := <-eventCh:
				typeEventRoutCh <- event
				publishEventRoutCh <- event
			case <-done:
				return
			}
		}
	}()

	log.Debug("маршрутизатор событий запущен")
	return eventCh, cancel
}

// Инициализатор маршрутизатора сообщений получателей
func initSubscribeRouter(typeSubscRoutCh, publishSubscRoutCh chan models.SubscriberMess) (chan models.SubscriberMess, func()) {
	subCh := make(chan models.SubscriberMess, 100)
	done := make(chan struct{})
	cancel := func() {
		close(done)
	}

	go func() {
		defer log.Debug("работа маршутизатора сообщений получателей завершена")
		for {
			select {
			case subMess := <-subCh:
				typeSubscRoutCh <- subMess
				publishSubscRoutCh <- subMess
			case <-done:
				return
			}
		}
	}()

	log.Debug("маршрутизатор сообщений получателей запущен")
	return subCh, cancel
}
