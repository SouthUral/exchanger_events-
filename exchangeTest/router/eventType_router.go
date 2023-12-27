package router

import (
	"context"

	log "github.com/sirupsen/logrus"
)

// маршрутизатор сообщений по типу
type typeRouter struct {
	routData
	typeRouters map[string]*typeEventRouter
}

// инициализация маршрутизатора событий по типам
func initTypeRouter(ctx context.Context, сapacityChanals int) typeRouter {
	router := &typeRouter{}

	router.typeRouters = make(map[string]*typeEventRouter, 100)
	router.eventCh = make(chan Event, сapacityChanals)
	router.subscrCh = make(chan SubscriberMess, сapacityChanals)

	go router.routing(ctx)

	return *router
}

// метод добавления нового маршрутизатора по названию типа события
func (t *typeRouter) addNewTypeRouter(ctx context.Context, typeEvent string) typeEventRouter {
	typeRouter := initTypeEventRouter(ctx, typeEvent)
	t.typeRouters[typeEvent] = &typeRouter
	return typeRouter
}

func (t *typeRouter) routing(ctx context.Context) {
	defer log.Debugf("Работа маршрутизатора типов событий завершена")

	log.Debug("Запущен маршрутизатор типов событий")

	// добавление маршрутизатора на все события
	allEventRouter := t.addNewTypeRouter(ctx, "allEvent")

	for {
		select {
		case event := <-t.eventCh:
			// отправка события маршрутизатору всех типов событий
			allEventRouter.eventCh <- event

			// отправка события роутеру этого события
			typeEvent := event.GetTypeEvent()
			typeRouter, ok := t.typeRouters[typeEvent]
			if ok {
				typeRouter.eventCh <- event
			} else {
				typeRouter := t.addNewTypeRouter(ctx, typeEvent)
				typeRouter.eventCh <- event
			}
		case subMess := <-t.subscrCh:
			// отправка подписчика маршрутизатору всех типов
			subConf := subMess.GetConfigSub()
			subName := subMess.GetNameSub()
			if subConf.GetAllEvent() {
				allEventRouter.subscrCh <- subMess
			} else {
				typesEvent := subMess.GetConfigSub().GetTypes()
				for _, typeEvent := range typesEvent {
					typeRouter, ok := t.typeRouters[typeEvent]
					if ok {
						typeRouter.subscrCh <- subMess
					} else {
						log.Warningf("для подписчика <%s> создан пустой маршрутизатор типа <%s>", subName, typeEvent)
						// если маршрутизатора типа нет, то он создается, и ему отправляется информация об подписчике
						typeRouter := t.addNewTypeRouter(ctx, typeEvent)
						typeRouter.subscrCh <- subMess
					}
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

// Запускает в отдельной горутине typeEventRouter (маршрутизатор типа)
func initTypeEventRouter(ctx context.Context, eventType string) typeEventRouter {
	router := &typeEventRouter{}

	router.eventCh = make(chan Event, 100)
	router.subscrCh = make(chan SubscriberMess, 100)
	router.typeEvent = eventType
	router.subscribers = make(map[string]chan interface{})

	go router.routing(ctx)

	return *router
}

// Маршрутизатор конкретного типа события.
// Содержит словарь со всеми подписчиками.
// При получении события отправляет его всем подписчикам.
// При получении подписчика, сохраняет его в свой словарь.
type typeEventRouter struct {
	routData
	typeEvent   string
	subscribers map[string]chan interface{}
}

// метод распределения событий по подписчикам
func (t *typeEventRouter) routing(ctx context.Context) {
	defer log.Debugf("Работа маршрутизатора типа %s завершена", t.typeEvent)

	log.Debugf("typeEventRouter для типа %s запущен", t.typeEvent)

	for {
		select {
		case event := <-t.eventCh:
			// TODO: можно добавить отписку подписчика от данного типа событий
			for subscr, ch := range t.subscribers {
				ch <- event
				log.Debugf("Событие типа %s отправлено подписчику %s", event.GetTypeEvent(), subscr)
			}
		case subMess := <-t.subscrCh:
			subName := subMess.GetNameSub()
			_, ok := t.subscribers[subName]
			if ok {
				log.Warningf("%s уже подписан на тип событий %s", subName, t.typeEvent)
			} else {
				t.subscribers[subName] = subMess.GetReverseCh()
			}
		case <-ctx.Done():
			return
		}
	}
}
