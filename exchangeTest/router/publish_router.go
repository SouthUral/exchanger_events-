package router

import (
	"context"
	"time"

	models "github.com/SouthUral/exchangeTest/models"
	log "github.com/sirupsen/logrus"
)

// Функция для инициализации маршрутизатора событий по отправителям
func initPublishRouter(ctx context.Context) eventRoutData {
	eventCh := make(models.EventChan, 100)
	subscrCh := make(chan models.SubscriberMess, 100)

	go publishRouter(ctx, eventCh, subscrCh)
	log.Debug("Запущен маршрутизатор событий по отправителю")

	return eventRoutData{
		eventCh:  eventCh,
		subscrCh: subscrCh,
	}
}

// Маршрутизатор событий по отправителю
func publishRouter(ctx context.Context, eventCh models.EventChan, subscrCh chan models.SubscriberMess) {
	defer log.Debugf("Работа маршрутизатора типов событий завершена")

	publishers := make(map[string]eventRoutData)

	for {
		select {
		case event := <-eventCh:
			publisherData, ok := publishers[event.Publisher]
			if ok {
				publisherData.eventCh <- event
			} else {
				publisherData := initPublisherEventRouter(ctx, event.Publisher)
				publishers[event.Publisher] = publisherData
				publisherData.eventCh <- event
			}
		case sub := <-subscrCh:
			subConf := sub.ConfSubscribe
			for _, pub := range subConf.Publishers {
				publisherData, ok := publishers[pub.Name]
				if ok {
					// в маршрутизатор отправителя попадет информация только по текущему отправителю
					publisherData.subscrCh <- models.SubscriberMess{
						Name: sub.Name,
						ConfSubscribe: models.ConfSub{
							Types:      subConf.Types,
							Publishers: []models.Publisher{pub},
							AllEvent:   subConf.AllEvent,
						},
						EvenCh: sub.EvenCh,
					}
				} else {
					// сообщение будет отсылаться до тех пор, пока не появится нужный отправитель
					log.Warningf("Подписчик %s не может подписаться на отправителя %s, этот отправитель не существует", sub.Name, pub.Name)
					go func(subMess models.SubscriberMess) {
						time.Sleep(5 * time.Second)
						log.Debugf("Повторная попытка подписать %s на отправителя %s", subMess.Name, subMess.ConfSubscribe.Publishers[0].Name)
						subscrCh <- subMess
						return
					}(models.SubscriberMess{
						Name: sub.Name,
						ConfSubscribe: models.ConfSub{
							Types:      subConf.Types,
							Publishers: []models.Publisher{pub},
							AllEvent:   subConf.AllEvent,
						},
						EvenCh: sub.EvenCh,
					})
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func initPublisherEventRouter(ctx context.Context, namePublisher string) eventRoutData {

	publisherData := eventRoutData{
		eventCh:  make(models.EventChan, 100),
		subscrCh: make(chan models.SubscriberMess, 100),
	}

	go publisherEventRouter(ctx, publisherData.eventCh, publisherData.subscrCh, namePublisher)
	log.Debugf("publisherEventRouter для отправителя %s запущен", namePublisher)

	return publisherData
}

func publisherEventRouter(ctx context.Context, eventCh models.EventChan, subscrCh chan models.SubscriberMess, namePublisher string) {
	defer log.Debugf("работа маршрутизатора событий по отправителю %s завершена", namePublisher)

	typeRoutData := initTypeRouter(ctx)

	for {
		select {
		case event := <-eventCh:
			typeRoutData.eventCh <- event
		case subMess := <-subscrCh:
			subConf := subMess.ConfSubscribe
			// если publisher.types ничего не содержит, тогда получатель подписывается на все типы событий
			if len(subConf.Publishers[0].TypeMess) == 0 {
				subConf.AllEvent = true
			}
			typeRoutData.subscrCh <- subMess
		case <-ctx.Done():
			return
		}
	}
}
