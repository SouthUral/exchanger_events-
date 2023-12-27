package utilsfortest

import (
	"context"

	log "github.com/sirupsen/logrus"
)

// функция запускает всех получателей, возвращает cancel для закрытия горутин получателей и канал, по которому прийдут отчеты
// Принимаемые параметры:
// confSubscribers - список конфигурации подписчиков
// subscriberChan - канал, который нужно получить от маршрутизатора, в него передается сообщение для подписки
// numberMessages - количество сообщений которое будет сгенерировано по каждому типу событий
// Возвращаемые параметры:
// func() - функция cancel() контекста для прекращния работы всех получателей
// chan reportSubTest - канал, по которому прийдут отчеты от горутин получателей
func SubscribersWork(confSubscribers []Consumer, subscriberChan chan<- interface{}, numberMessages int) (func(), chan ReportSubTest) {
	reportCh := make(chan ReportSubTest, len(confSubscribers))
	ctx, cancel := context.WithCancel(context.Background())
	for _, itemConsum := range confSubscribers {
		subMes := initSubMess(itemConsum)
		subscriberChan <- subMes
		initSubTest(ctx, subMes, numberMessages, reportCh)

	}
	return cancel, reportCh
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
	eventCh                   chan interface{}
	reportCh                  chan ReportSubTest
	expectedNumberMess        int
	counterReceivedMess       int
	allEvent                  bool
	expectedlistCountEvents   map[keyEvent]int
	unexpectedListCountEvents map[keyEvent]int
}

// стуктура для отчета
type ReportSubTest struct {
	NameSub                   string
	ExpectedNumberMess        int
	CounterReceivedMess       int
	AllEvent                  bool
	ExpectedlistCountEvents   map[keyEvent]int
	UnexpectedListCountEvents map[keyEvent]int
}

func (s *subTest) worker(ctx context.Context) {
	if s.allEvent {
		go s.processWithAllEvent(ctx)
	} else {
		go s.processStandart(ctx)
	}
}

// метод для генерации и отправки отчета
func (s *subTest) generatingReport() {
	report := ReportSubTest{
		NameSub:                   s.nameSub,
		ExpectedNumberMess:        s.expectedNumberMess,
		CounterReceivedMess:       s.counterReceivedMess,
		AllEvent:                  s.allEvent,
		ExpectedlistCountEvents:   s.expectedlistCountEvents,
		UnexpectedListCountEvents: s.unexpectedListCountEvents,
	}

	s.reportCh <- report
}

// процесс получения событий, запускаемый если allEvent == true
func (s *subTest) processWithAllEvent(ctx context.Context) {
	for {
		select {
		case <-s.eventCh:
			s.counterReceivedMess++
		case <-ctx.Done():
			s.generatingReport()
			return
		}
	}

}

// процесс получения событий, запускаемый если флаг allEvent == false
func (s *subTest) processStandart(ctx context.Context) {
	for {
		select {
		case event := <-s.eventCh:
			eventMsg, ok := event.(Event)
			s.counterReceivedMess++

			if !ok {
				log.Errorf("получателю %s не удалось привести к типу событие", s.nameSub)
				continue
			}

			// формируется ключ, в котором присутствет и отправитель и тип события
			keyWithPublisher := keyEvent{
				publisherName: eventMsg.GetPub(),
				typeName:      eventMsg.GetTypeEvent(),
			}
			s.checkEvent(keyWithPublisher)

			// формируется ключ, в котором есть только тип события
			keyWithoutPublisher := keyEvent{
				typeName: eventMsg.GetTypeEvent(),
			}
			s.checkEvent(keyWithoutPublisher)

		case <-ctx.Done():
			s.generatingReport()
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

func initSubTest(ctx context.Context, subMess subMess, numberMess int, reportCh chan ReportSubTest) subTest {
	res := subTest{
		nameSub:  subMess.GetNameSub(),
		eventCh:  subMess.GetReverseCh(),
		allEvent: subMess.GetConfigSub().GetAllEvent(),
		reportCh: reportCh,
	}

	if res.allEvent {
		return res
	}

	res.expectedNumberMess = countAllEvents(subMess.Config, numberMess)
	res.expectedlistCountEvents = initExpectedListCountEvents(subMess.Config)
	res.unexpectedListCountEvents = make(map[keyEvent]int)

	res.worker(ctx)

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
