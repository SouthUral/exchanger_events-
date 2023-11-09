package router

import (
	"time"

	conf "github.com/SouthUral/exchangeTest/confreader"

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
			for _, pub := range sub.Publishers {
				publisherData, ok := publishers[pub.Name]
				if ok {
					// в маршрутизатор отправителя попадет информация только по текущему отправителю
					publisherData.subscrCh <- SubscriberMess{
						Name:       sub.Name,
						Types:      sub.Types,
						Publishers: []conf.Publisher{pub},
						AllEvent:   sub.AllEvent,
						EvenCh:     sub.EvenCh,
					}
				} else {
					// сообщение будет отсылаться до тех пор, пока не появится нужный отправитель
					log.Warningf("Подписчик %s не может подписаться на отправителя %s, этот отправитель не существует", sub.Name, pub.Name)
					go func(subMess SubscriberMess) {
						time.Sleep(5 * time.Second)
						log.Debugf("Повторная попытка подписать %s на отправителя %s", subMess.Name, subMess.Publishers[0].Name)
						subscrCh <- subMess
						return
					}(SubscriberMess{
						Name:       sub.Name,
						Types:      sub.Types,
						Publishers: []conf.Publisher{pub},
						AllEvent:   sub.AllEvent,
						EvenCh:     sub.EvenCh,
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
			if len(sub.Publishers[0].TypeMess) == 0 {
				sub.AllEvent = true
			}
			typeRoutData.subscrCh <- sub
		case <-done:
			typeRoutData.cancel()
			return
		}
	}
}
