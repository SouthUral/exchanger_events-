package utilsfortest

// TODO: получение конфига, из которого будут генерироваться отправители;
// TODO: получение канала, куда нужно отправлять события;
// TODO: получение даннх для подключения к Rabbit;

// TODO: инициализация отправителей, основная функция, которая принимает все параметры и запускает в цикле горутины;
// TODO: горутина отправителя: генерация события, отправка события во внутренний роут (канал маршрутизатора),
// отправка события в rabbit
// TODO: приостановка отправки событий если Rabbit не работает, переподключение Rabbit

// TODO: попробуем сделать отправителя без отправки сообщения в RabbitMQ!!!

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// Функция запуска отправителей;
// eventCh - канал для отправки событий во внутренний маршрутизатор;
// confPublishers - структура с кофигурациями отправителей;
// rmqURL - URL для подключения к RabbitMQ
// numberMessages - количество сообщений, которые должны отправить отправители
func StartPublishers(confPublishers []Publisher, eventChan chan Event, numberMessages int) func() {
	ctx, cancel := context.WithCancel(context.Background())

	for _, pubslishConf := range confPublishers {
		startPublisher(ctx, pubslishConf, eventChan, numberMessages)
	}

	return cancel
}

// генерирует и отправляет события во внутренний маршрутизатор и в exchange Rabbit
func startPublisher(ctx context.Context, confPublisher Publisher, eventCh chan Event, numberMessages int) {

	go func() {
		defer log.Debugf("Отправитель %s прекратил работу", confPublisher.Name)

		for _, typeEvent := range confPublisher.TypeMess {
			startGenEvent(ctx, typeEvent, confPublisher.Name, eventCh, numberMessages)
		}

	}()
}

// Генератор сообщения для типа события
func startGenEvent(ctx context.Context, nameType, namePublisher string, eventCh chan Event, numberMessages int) {

	go func() {
		defer log.Debugf("Генератор событий %s.%s прекратил работу", namePublisher, nameType)
		for i := 0; i < numberMessages; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				event := event{
					id:          i,
					publisher:   namePublisher,
					typeEvent:   nameType,
					timePublish: time.Now(),
					message:     fmt.Sprintf("Событие %s:%d", nameType, i),
				}
				eventCh <- event
				time.Sleep(5 * time.Millisecond)
			}
		}
	}()
}
