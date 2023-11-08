package router

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Функция для инициализации маршрутизатора событий по отправителям
func initPublishRouter() eventRoutData {
	eventCh := make(EventChan, 100)
	subscrCh := make(SubscriberChan, 100)
	done := make(chan struct{})
	cancel := func() {
		close(done)
	}

	go publishRouter(eventCh, subscrCh, done)
	log.Debug("Запущен маршрутизатор событий по отправителю")

	return eventRoutData{
		eventCh:  eventCh,
		subscrCh: subscrCh,
		cancel:   cancel,
	}
}

// Маршрутизатор событий по отправителю
func publishRouter(eventCh EventChan, subscrCh SubscriberChan, done chan struct{}) {
	defer log.Debugf("Работа маршрутизатора типов событий завершена")

	publishers := make(map[string]eventRoutData)

	for {
		select {
		case event := <-eventCh:
			publisherData, ok := publishers[event.Publisher]
			if ok {
				publisherData.eventCh <- event
			} else {
				publisherData := initPublisherEventRouter(event.Publisher)
				publishers[event.Publisher] = publisherData
				publisherData.eventCh <- event
			}
		case sub := <-subscrCh:
			for _, pub := range sub.publishers {
				publisherData, ok := publishers[pub.name]
				if ok {
					// в маршрутизатор отправителя попадет информация только по текущему отправителю
					publisherData.subscrCh <- SubscriberMess{
						name:       sub.name,
						types:      sub.types,
						publishers: []publisher{pub},
						allEvent:   sub.allEvent,
						evenCh:     sub.evenCh,
					}
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

func initPublisherEventRouter(namePublisher string) eventRoutData {
	done := make(chan struct{})

	publisherData := eventRoutData{
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
	defer log.Debugf("работа маршрутизатора событий по отправителю %s завершена", namePublisher)

	typeRoutData := initTypeRouter()

	for {
		select {
		case event := <-eventCh:
			typeRoutData.eventCh <- event
		case sub := <-subscrCh:
			// если publisher.types ничего не содержит, тогда получатель подписывается на все типы событий
			if len(sub.publishers[0].types) == 0 {
				sub.allEvent = true
			}
			typeRoutData.subscrCh <- sub
		case <-done:
			typeRoutData.cancel()
			return
		}
	}
}
