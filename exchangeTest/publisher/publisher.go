package publisher

// TODO: получение конфига, из которого будут генерироваться отправители;
// TODO: получение канала, куда нужно отправлять события;
// TODO: получение даннх для подключения к Rabbit;

// TODO: инициализация отправителей, основная функция, которая принимает все параметры и запускает в цикле горутины;
// TODO: горутина отправителя: генерация события, отправка события во внутренний роут (канал маршрутизатора),
// отправка события в rabbit
// TODO: приостановка отправки событий если Rabbit не работает, переподключение Rabbit

// TODO: попробуем сделать отправителя без отправки сообщения в RabbitMQ!!!

import (
	"fmt"
	"time"

	conf "github.com/SouthUral/exchangeTest/confreader"
	router "github.com/SouthUral/exchangeTest/router"

	log "github.com/sirupsen/logrus"
)

// Функция запуска отправителей;
// eventCh - канал для отправки событий во внутренний маршрутизатор;
// confPublishers - структура с кофигурациями отправителей;
// rmqURL - URL для подключения к RabbitMQ
func StartPublishers(eventCh router.EventChan, confPublishers []conf.Publisher, eventChan router.EventChan) func() {
	publishersStorage := make(map[string]func())

	closeAllPublishers := func() {
		for _, item := range publishersStorage {
			item()
		}
	}

	for _, pubslishConf := range confPublishers {
		publishersStorage[pubslishConf.Name] = publisher(pubslishConf, eventChan)
	}

	return closeAllPublishers
}

// генерирует и отправляет события во внутренний маршрутизатор и в exchange Rabbit
func publisher(confPublisher conf.Publisher, eventCh router.EventChan) func() {
	// TODO: функция генерации события (горутина, отправляет событие в Rabbit горутину и в маршрутизатор)
	// TODO: отправка события в внутренний маршрутизатор
	done := make(chan struct{})
	cancel := func() {
		close(done)
	}
	selfEventsCh := make(router.EventChan, 100)

	go func() {
		defer log.Debugf("Отправитель %s прекратил работу", confPublisher.Name)
		genEventStorage := make(map[string]func())
		for _, typeEvent := range confPublisher.TypeMess {
			genEventStorage[typeEvent] = genEvent(typeEvent, confPublisher.Name, selfEventsCh)
		}

		for {
			select {
			case event := <-selfEventsCh:
				eventCh <- event
			case <-done:
				for _, item := range genEventStorage {
					item()
				}
				return
			}
		}
	}()

	return cancel
}

// Генератор событий
func genEvent(nameType, namePublisher string, eventCh router.EventChan) func() {
	done := make(chan struct{})
	cancel := func() {
		close(done)
	}
	go func() {
		defer log.Debugf("Генератор событий %s.%s прекратил работу", namePublisher, nameType)
		for i := 0; i < 1000; i++ {
			select {
			case <-done:
				return
			default:
				event := router.Event{
					Id:          i,
					Publisher:   namePublisher,
					TypeEvent:   nameType,
					TimePublish: time.Now(),
					Message:     fmt.Sprintf("Событие %s:%d", nameType, i),
				}
				eventCh <- event
				time.Sleep(1 * time.Second)
			}
		}
	}()
	return cancel
}
