package router

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Функция для инициализации маршрутизатора событий по отправителям
func initPublishRouter() (EventChan, SubscriberChan, func()) {
	eventCh := make(EventChan, 100)
	subscrCh := make(SubscriberChan, 100)
	done := make(chan struct{})
	cancel := func() {
		close(done)
	}

	go publishRouter(eventCh, subscrCh, done)
	log.Debug("Запущен маршрутизатор событий по отправителю")

	return eventCh, subscrCh, cancel
}

// Маршрутизатор событий по отправителю
func publishRouter(eventCh EventChan, subscrCh SubscriberChan, done chan struct{}) {
	defer log.Debugf("Работа маршрутизатора типов событий завершена")

	publishers := make(map[string]publisherRoutData)

	for {
		select {
		case event := <-eventCh:
			publisherData, ok := publishers[event.publisher]
			if ok {
				publisherData.eventCh <- event
			} else {
				publisherData := initPublisherEventRouter(event.publisher)
				publishers[event.publisher] = publisherData
				publisherData.eventCh <- event
			}
		case sub := <-subscrCh:
			for _, pub := range sub.publishers {
				publisherData, ok := publishers[pub.name]
				if ok {
					publisherData.subscrCh <- sub
				} else {
					// сообщение будет отсылаться до тех пор, пока не появится нужный отправитель
					log.Warningf("Подписчик %s не может подписаться на отправителя %s, этот отправитель не существует", sub.name, pub.name)
					go func(subMess SubscriberMess) {
						time.Sleep(5 * time.Second)
						log.Debugf("Повторная попытка подписать %s на отправителя %s", subMess.name, subMess.publishers[0].name)
						subscrCh <- subMess
						return
					}(SubscriberMess{
						name:       sub.name,
						types:      sub.types,
						publishers: []publisher{pub},
						allEvent:   sub.allEvent,
						evenCh:     sub.evenCh,
					})
				}
			}
		case <-done:
			for _, publisherData := range publishers {
				publisherData.cancel()
			}
			return
		}
	}
}

func initPublisherEventRouter(namePublisher string) publisherRoutData {
	done := make(chan struct{})

	publisherData := publisherRoutData{
		eventCh:  make(EventChan, 100),
		subscrCh: make(SubscriberChan, 100),
		cancel: func() {
			close(done)
		},
	}

	go publisherEventRouter(publisherData.eventCh, publisherData.subscrCh, done, namePublisher)
	log.Debugf("publisherEventRouter для отправителя %s запущен", namePublisher)

	return publisherData
}

func publisherEventRouter(eventCh EventChan, subscrCh SubscriberChan, done chan struct{}, namePublisher string) {
	eventTypeCh, subscrTypeCh, typeClose := initTypeRouter()

	for {
		select {
		case event := <-eventCh:
			eventTypeCh <- event
		case sub := <-subscrCh:
			subscrTypeCh <- sub
		case <-done:
			typeClose()
			return
		}
	}
}
