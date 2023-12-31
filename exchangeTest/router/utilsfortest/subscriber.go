package utilsfortest

import (
	"context"
	"sync"
	"time"

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
func SubscribersWork(confSubscribers []Consumer, subscriberChan chan interface{}, numberMessages int) (chan ReportSubTest, context.Context) {
	wg := sync.WaitGroup{}
	linkWg := &wg
	ctxSubWork, cancelSubWork := context.WithCancel(context.Background())

	reportCh := make(chan ReportSubTest, len(confSubscribers)+20)
	go func() {
		defer cancelSubWork()
		defer log.Debug("все получатели завершили работу")
		for _, itemConsum := range confSubscribers {
			wg.Add(1)

			subMes := initSubMess(itemConsum)
			log.Debug(itemConsum, "Это сообщение от отправителя!!!")
			subscriberChan <- subMes

			ctx, cancel := context.WithCancel(context.Background())
			signalCh := UpdatableTimer(cancel, 10)
			initSubTest(ctx, subMes, numberMessages, reportCh, signalCh, linkWg)
		}
		wg.Wait()
	}()

	return reportCh, ctxSubWork
}

// обновляемый таймер
func UpdatableTimer(cancel func(), timeWait int) chan struct{} {
	signalCh := make(chan struct{})

	go func() {
		defer cancel()
		defer log.Warning("Работа таймера закончена")
		// defer wg.Done()

		duration := time.Duration(timeWait) * time.Second

		ctx, _ := context.WithTimeout(context.Background(), duration)

		for {
			select {
			case <-signalCh:
				log.Debug("Обновлен таймер")
				ctx, _ = context.WithTimeout(context.Background(), duration)
			case <-ctx.Done():
				return
			}
		}
	}()

	return signalCh

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
	signalCh                  chan struct{}
	reportCh                  chan ReportSubTest
	expectedNumberMess        int
	counterReceivedMess       int
	allEvent                  bool
	expectedlistCountEvents   map[keyEvent]int
	unexpectedListCountEvents map[keyEvent]int
	wg                        *sync.WaitGroup
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
	defer s.wg.Done()
	for {
		select {
		case <-s.eventCh:
			s.signalCh <- struct{}{}
			log.Infof("получатель %s получил событие", s.nameSub)
			s.counterReceivedMess++
		case <-ctx.Done():
			s.generatingReport()
			return
		}
	}

}

// процесс получения событий, запускаемый если флаг allEvent == false
func (s *subTest) processStandart(ctx context.Context) {
	defer s.wg.Done()
	for {
		select {
		case event := <-s.eventCh:
			s.signalCh <- struct{}{}
			log.Infof("получатель %s получил событие", s.nameSub)
			eventMsg, ok := event.(Event)
			s.counterReceivedMess++

			if !ok {
				log.Errorf("получателю %s не удалось привести к типу событие", s.nameSub)
				continue
			}

			s.checkEvent(eventMsg.GetPub(), eventMsg.GetTypeEvent())

		case <-ctx.Done():
			log.Debugf("Получатель %s готов отправить отчет", s.nameSub)
			s.generatingReport()
			return
		}
	}
}

// метод проверяет, если ли ключ в списке ожидаемых событий, или в списке неожидаемых событий
func (s *subTest) checkEvent(pub, typeEvent string) {
	keyPub := keyEvent{
		publisherName: pub,
	}

	_, ok := s.expectedlistCountEvents[keyPub]
	if ok {
		s.expectedlistCountEvents[keyPub]++
		return
	}

	keyFull := keyEvent{
		publisherName: pub,
		typeName:      typeEvent,
	}

	_, ok = s.expectedlistCountEvents[keyFull]
	if ok {
		s.expectedlistCountEvents[keyFull]++
		return
	}

	keyType := keyEvent{
		typeName: typeEvent,
	}

	_, ok = s.expectedlistCountEvents[keyType]
	if ok {
		s.expectedlistCountEvents[keyType]++
		return
	}

	_, ok = s.unexpectedListCountEvents[keyFull]
	if ok {
		s.unexpectedListCountEvents[keyFull]++
	} else {
		s.unexpectedListCountEvents[keyFull] = 1
	}
}

func initSubTest(ctx context.Context, subMess subMess, numberMess int, reportCh chan ReportSubTest, signalCh chan struct{}, wg *sync.WaitGroup) subTest {
	res := subTest{
		nameSub:  subMess.GetNameSub(),
		eventCh:  subMess.GetReverseCh(),
		allEvent: subMess.GetAllEvent(),
		reportCh: reportCh,
		signalCh: signalCh,
		wg:       wg,
	}

	if res.allEvent {
		return res
	}

	res.expectedNumberMess = countAllEvents(subMess, numberMess)
	res.expectedlistCountEvents = initExpectedListCountEvents(subMess)
	res.unexpectedListCountEvents = make(map[keyEvent]int)

	res.worker(ctx)

	return res
}

// функция генерирует плоский словарь типа <отправитель_тип: количество полученных сообщений>
func initExpectedListCountEvents(config subMess) map[keyEvent]int {
	result := make(map[keyEvent]int, 20)
	for publisher, types := range config.GetPublihers() {
		if len(types) == 0 {
			key := keyEvent{
				publisherName: publisher,
			}
			result[key] = 0
		} else {
			for _, t := range types {
				key := keyEvent{
					publisherName: publisher,
					typeName:      t,
				}
				result[key] = 0
			}
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
func countAllEvents(config subMess, numberMessages int) int {
	var result int

	if config.GetAllEvent() {
		return result
	}

	for _, types := range config.GetPublihers() {
		result = result + len(types)*numberMessages
	}

	result = +len(config.GetTypes())

	return result
}
