package router

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type publishRouter struct {
	routData
	pubEventRouters map[string]typeRouter
}

// Функция для инициализации маршрутизатора событий по отправителям
func initPublishRouter(ctx context.Context, сapacityChanals int) publishRouter {
	router := &publishRouter{}

	router.eventCh = make(chan Event, сapacityChanals)
	router.subscrCh = make(chan SubscriberMess, сapacityChanals)
	router.pubEventRouters = make(map[string]typeRouter, 100)

	go router.routing(ctx)

	return *router
}

// Маршрутизатор событий по отправителю
func (p *publishRouter) routing(ctx context.Context) {
	defer log.Debugf("Работа маршрутизатора типов событий завершена")

	log.Debug("Запущен маршрутизатор событий по отправителю")

	for {
		select {
		case event := <-p.eventCh:
			eventPub := event.GetPub()
			publisherData, ok := p.pubEventRouters[eventPub]
			if ok {
				publisherData.eventCh <- event
			} else {
				publisherData := initTypeRouter(ctx, 100)
				p.pubEventRouters[eventPub] = publisherData
				publisherData.eventCh <- event
			}
		case subMess := <-p.subscrCh:
			subPublishers := subMess.GetConfigSub().GetPublihers()
			for pubName, eventTypes := range subPublishers {
				typeRouter, ok := p.pubEventRouters[pubName]
				copyMess := copySubMess(subMess, eventTypes)
				if ok {
					typeRouter.subscrCh <- &copyMess
				} else {
					publisherData := initTypeRouter(ctx, 100)
					p.pubEventRouters[pubName] = publisherData
					publisherData.subscrCh <- &copyMess
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

// функция копирует сообщение от подписчика и берет типы событий одного отправителя из конфига и ставит их на место всех типов в конфиге.
// Нужно для работы роутера по отправителям
func copySubMess(mess SubscriberMess, evenTypes []string) subMess {
	newMess := &subMess{}
	newMess.name = mess.GetNameSub()
	newMess.reverseCh = mess.GetReverseCh()

	newConfig := &confSub{}

	newConfig.types = evenTypes
	newConfig.allEvent = mess.GetConfigSub().GetAllEvent()

	newMess.config = newConfig

	return *newMess
}
