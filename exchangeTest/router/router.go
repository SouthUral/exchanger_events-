package router

import (
	"context"

	log "github.com/sirupsen/logrus"
)

// Модуль в котором производится маршрутизация событий по получателям
// Задачи:

// TODO: Создание канала для получения событий
// TODO: Создание канала для получения подписчиков

// TODO: Функукция для маршрутизации типов
// TODO: Функция для маршрутизации по отправителям
// TODO: Маршрутизатор событий (принимает входящие каналы(для событий) от маршрутизаторов: типа, отправителя, по всем)
// TODO: Маршрутизатор подписчиков (принимает входящие каналы подписчиков от маршрутизаторов)

// Функция для запуска маршрутизатора
func InitRouter() (chan Event, chan SubscriberMess, func()) {
	// TODO: нужно понять какой буфер делать у каналов

	ctx, cancel := context.WithCancel(context.Background())

	typeRoutData := initTypeRouter(ctx, 100)
	publishRoutData := initPublishRouter(ctx, 100)

	eventCh := initEventRouter(ctx, typeRoutData.eventCh, publishRoutData.eventCh)
	subCh := initSubscribeRouter(ctx, typeRoutData.subscrCh, publishRoutData.subscrCh)

	return eventCh, subCh, cancel
}

// Инициализатор маршрутизатора событий
func initEventRouter(ctx context.Context, typeEventRoutCh, publishEventRoutCh chan Event) chan Event {
	eventCh := make(chan Event, 100)

	go func() {
		defer log.Debug("работа маршрутизатора событий завершена")
		for {
			select {
			case event := <-eventCh:
				typeEventRoutCh <- event
				publishEventRoutCh <- event
			case <-ctx.Done():
				return
			}
		}
	}()

	log.Debug("маршрутизатор событий запущен")
	return eventCh
}

// Инициализатор маршрутизатора сообщений получателей
func initSubscribeRouter(ctx context.Context, typeSubscRoutCh, publishSubscRoutCh chan SubscriberMess) chan SubscriberMess {
	subCh := make(chan SubscriberMess, 100)

	go func() {
		defer log.Debug("работа маршутизатора сообщений получателей завершена")
		for {
			select {
			case subMess := <-subCh:
				typeSubscRoutCh <- subMess
				publishSubscRoutCh <- subMess
			case <-ctx.Done():
				return
			}
		}
	}()

	log.Debug("маршрутизатор сообщений получателей запущен")
	return subCh
}
