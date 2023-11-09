package subscriber

import (
	"time"

	conf "github.com/SouthUral/exchangeTest/confreader"
	router "github.com/SouthUral/exchangeTest/router"

	log "github.com/sirupsen/logrus"
)

func StartSubscribers(confSubscribers []conf.Consumer, subscriberChan router.SubscriberChan) func() {
	consumers := make(map[string]func())
	cancelAllSubscriber := func() {
		for _, itemFunc := range consumers {
			itemFunc()
		}
	}

	for _, item := range confSubscribers {
		consumers[item.Name] = Subscriber(item, subscriberChan)
	}

	return cancelAllSubscriber
}

// Функция для запуска подписчика.
// Отправляет информацию о подписчике в router.SubscriberChan
// TODO: нужна проверка конфига подписчика (типы событий не должны пересекаться с типами событий в Отправителе) иначе будут дубли
// TODO: запустить цикл прослушки (получение) событий
// TODO: нужно сразу освобождать канал, для того чтобы там не накапливались события
// TODO: нужна проверка уникальности события, оно не должно приходить два раза, при проверке выдавать ошибку, если уже есть похожее событие
func Subscriber(subscribeConf conf.Consumer, subscriberChan router.SubscriberChan) func() {
	selfEventCh := make(router.EventChan, 100)

	done := make(chan struct{})
	cancel := func() {
		close(done)
	}
	// TODO: тут нужно где-то проверить конфиг, чтобы типы в отправителях не пересекались с общими типами
	subMes := router.SubscriberMess{
		Name:       subscribeConf.Name,
		Types:      subscribeConf.Types,
		Publishers: subscribeConf.Publishers,
		AllEvent:   subscribeConf.AllEvent,
		EvenCh:     selfEventCh,
	}

	go func() {
		defer log.Warningf("Подписчик %s прекратил работу", subscribeConf.Name)
		for {
			select {
			case event := <-selfEventCh:
				timeDifference := time.Since(event.TimePublish)
				log.Infof("Время от генерации события %s.%s до получения: %s", event.Publisher, event.TypeEvent, timeDifference)
				log.Infof("Событие %d.%s.%s получено получателем %s", event.Id, event.Publisher, event.TypeEvent, subscribeConf.Name)
				// TODO: далее события нужно записать в кэш
			case <-done:
				return
			}
		}
	}()

	subscriberChan <- subMes

	return cancel
}
