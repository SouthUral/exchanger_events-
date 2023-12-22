package utilsfortest

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

func StartSubscribers(confSubscribers []Consumer, subscriberChan chan SubscriberMess, numberMessages int) func() {
	ctx, cancel := context.WithCancel(context.Background())
	for _, item := range confSubscribers {
		Subscriber(ctx, item, subscriberChan, numberMessages)
	}
	return cancel
}

// Функция для запуска подписчика.
// Отправляет информацию о подписчике в router.SubscriberChan
// TODO: нужна проверка конфига подписчика (типы событий не должны пересекаться с типами событий в Отправителе) иначе будут дубли
// TODO: запустить цикл прослушки (получение) событий
// TODO: нужно сразу освобождать канал, для того чтобы там не накапливались события
// TODO: нужна проверка уникальности события, оно не должно приходить два раза, при проверке выдавать ошибку, если уже есть похожее событие
func Subscriber(ctx context.Context, subscribeConf Consumer, subscriberChan chan SubscriberMess, resCh chan resultWork, numberMessages int) {
	// TODO: тут нужно где-то проверить конфиг, чтобы типы в отправителях не пересекались с общими типами

	subMes := initSubMess(subscribeConf)
	listTypes := genListTypes(subMes.Config)
	expectedNumberMess := countAllEvents(subMes.Config, numberMessages)

	go func() {
		defer log.Warningf("Подписчик %s прекратил работу", subscribeConf.Name)

		var counterReceivedMess int

		res := resultWork{
			expectedNumberMess: expectedNumberMess,
		}

		for {
			select {
			case event := <-subMes.GetReverseCh():
				counterReceivedMess++

				timeDifference := time.Since(event.GetTimePub())
				log.Infof("time difference = %v, subscriber: %s\n", timeDifference, subMes.Name)
			case <-ctx.Done():
				break
			}
		}

		res.counterReceivedMess = counterReceivedMess

		resCh <- res
		return

	}()

	subscriberChan <- subMes

}

type keyEvent struct {
	publisherName string
	typeName      string
}

func (k keyEvent) GetPublisherName() string {
	return k.publisherName
}

func (k keyEvent) GetTypeName() string {
	return k.typeName
}

type subTest struct {
	nameSub                   string
	eventCh                   chan Event
	expectedNumberMess        int
	counterReceivedMess       int
	allEvent                  bool
	expectedlistCountEvents   map[keyEvent]int
	unexpectedListCountEvents map[keyEvent]int
}

func (s *subTest) worker() {

}

func (s *subTest) generatingReport() {
	// TODO: нужно создать структура для отчета и заполнить его
	// TODO: нужно создать канал для отчетов, определить место его инициализации и передать туда отчет
	// TODO: нужно отловить все отчеты горутин из каналов и вернуть их списком или структурой
}

// процесс, запускаемый если allEvent == true
func (s *subTest) processWithAllEvent(ctx context.Context) {
	for {
		select {
		case <-s.eventCh:
			s.counterReceivedMess++
		case <-ctx.Done():
			return
		}
	}

}

func (s *subTest) processStandart(ctx context.Context) {
	for {
		select {
		case event := <-s.eventCh:
			s.counterReceivedMess++
			keyWithPublisher := keyEvent{
				publisherName: event.GetPub(),
				typeName:      event.GetTypeEvent(),
			}
			s.checkEvent(keyWithPublisher)

			keyWithoutPublisher := keyEvent{
				typeName: event.GetTypeEvent(),
			}
			s.checkEvent(keyWithoutPublisher)

		case <-ctx.Done():
			return
		}
	}
}

// метод проверяет, если ли ключ в списке ожидаемых событий, или в списке неожидаемых событий
func (s *subTest) checkEvent(key keyEvent) {
	_, ok := s.expectedlistCountEvents[key]
	if ok {
		s.expectedlistCountEvents[key]++
	} else {
		_, ok := s.unexpectedListCountEvents[key]
		if ok {
			s.unexpectedListCountEvents[key]++
		} else {
			s.unexpectedListCountEvents[key] = 1
		}
	}
}

func initSubTest(subMess subMess, numberMess int) subTest {
	res := subTest{
		nameSub:  subMess.GetNameSub(),
		eventCh:  subMess.GetReverseCh(),
		allEvent: subMess.GetConfigSub().GetAllEvent(),
	}

	if res.allEvent {
		return res
	}

	res.expectedNumberMess = countAllEvents(subMess.Config, numberMess)
	res.expectedlistCountEvents = initExpectedListCountEvents(subMess.Config)
	res.unexpectedListCountEvents = make(map[keyEvent]int)

	return res
}

// функция генерирует плоский словарь типа <отправитель_тип: количество полученных сообщений>
func initExpectedListCountEvents(config confSub) map[keyEvent]int {
	result := make(map[keyEvent]int, 20)
	for publisher, types := range config.GetPublihers() {
		for _, t := range types {
			key := keyEvent{
				publisherName: publisher,
				typeName:      t,
			}
			result[key] = 0
		}
	}

	for _, t := range config.GetTypes() {
		key := keyEvent{
			typeName: t,
		}
		result[key] = 0
	}

	return result
}

// Функция считает количество всех сообщений которые должен получить подписчик
func countAllEvents(config confSub, numberMessages int) int {
	var result int

	if config.GetAllEvent() {
		return result
	}

	for _, types := range config.GetPublihers() {
		counter := len(types) * numberMessages
		result = +counter
	}

	result = +len(config.GetTypes())

	return result
}
