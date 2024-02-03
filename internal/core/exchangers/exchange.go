package exchangers

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// Внутренний обменник, получает события и присваивает им offset.
// Отдает события очередям, которые подписаны на него.
// Отправляет offset события и его id на сохранение в БД.
type InternalExchange struct {
	Id               string                    // внутренний ID по которому будет работать маршрутизатор
	mx               sync.Mutex                // мьютекс, нужен для блокировки доступа
	ConfConsum       configConsum              // конфигурация подписки
	activeQueue      map[string]*queueConsumer // активные подписчики
	inactiveQueue    map[string]*queueConsumer // не активыне подписчики
	IncomingEventsCh chan interface{}          // канал для получения событий из маршрутизатора
	cancel           func()                    // функция для отмены контекста, внутренний атрибут
	lastOffset       int
}

// Инициализация exchange
func InitInternalExchange(conf recipConfig) *InternalExchange {
	ctx, cancel := context.WithCancel(context.Background())

	// инициализация структуры конфигурации подписчика
	config := configConsum{
		Types:      conf.GetTypes(),
		Publishers: conf.GetPublihers(),
		AllEvent:   conf.GetAllEvent(),
	}

	intCons := &InternalExchange{
		Id:               uuid.New().String(),
		ConfConsum:       config,
		mx:               sync.Mutex{},
		activeQueue:      make(map[string]*queueConsumer),
		inactiveQueue:    make(map[string]*queueConsumer),
		IncomingEventsCh: make(chan interface{}, 5),
		cancel:           cancel,
	}

	// запуск маршрутизации событий
	go intCons.eventRouter(ctx)

	log.Debugf("инициализирован внутренний потребитель %s", intCons.Id)

	return intCons
}

// Метод запускается как горутина.
// Получает события и отправляет его всем активным очередям и отправляет в канал на сохранение в БД
func (i *InternalExchange) eventRouter(ctx context.Context) {
	defer log.Infof("eventRouter has stopped working")
	for {
		select {
		case eventMess := <-i.IncomingEventsCh:
			event, err := i.eventConversion(eventMess)
			if err != nil {
				// если будет ошибка конвертации, то нужно прекратить работу программы
				return
			}
			i.sendingEventsQueue(event)

			// отправка события на запись в БД
			fmt.Println(eventMess)
		case <-ctx.Done():
			return
		}
	}

}

// отправка события всем активным очередям
func (i *InternalExchange) sendingEventsQueue(event event) {
	for name, queue := range i.activeQueue {
		res := queue.receivingEvent(event)
		if !res {
			i.sleepQueue(name)
		}
	}
}

// преобразование полученного интерфейса от маршрутизатора к внутренней структуре
func (i *InternalExchange) eventConversion(msg interface{}) (event, error) {
	var err error
	eventInterface, ok := msg.(eventFromRouter)
	if !ok {
		err = eventConversionError{}
		log.Error(err)
	}

	offset := i.getOffsetEvent()

	newEvent := event{
		offset:    offset,
		publisher: eventInterface.GetPub(),
		typeEvent: eventInterface.GetTypeEvent(),
		timePub:   eventInterface.GetTimePub(),
		msg:       eventInterface.GetMess(),
	}

	return newEvent, err
}

// метод возвращает offset который можно присвоить новому событию.
func (i *InternalExchange) getOffsetEvent() int {
	defer i.mx.Unlock()
	i.mx.Lock()
	offset := i.lastOffset + 1
	i.lastOffset = offset
	return offset
}

// TODO: Метод добавления подписчика
func (i *InternalExchange) AddConsumer(sub subInfo) error {
	// defer log.Infof("a subscriber: %s has been added to the queue", sub.name)

	i.mx.Lock()
	_, ok := i.activeQueue[sub.GetName()]

	// Проверка в активных подписчиках

	// Если есть в активных подписчиках проверить есть ли у него сейчас подключение
	// Если подключение есть то вернуть ошибку, что такой подписчик уже активен
	// Если подключения нет, то добавить его

	// Проверка в некативных подписчиках

	// Если есть в неактивных подписчиках то перевести его в активные подписчики (произвести запуск горутины)

	if ok {
		// _ := fmt.Errorf("")
	}
	// нужно проверить активность подписчика
	// если подписчик активен, то вернуть ошибку
	// если подписчик неактивен, то добавить ему новый канал и включить активность
	// если подписчика нет, то нужно создать его и добавить в подписчики
	//

	// cons.Subscribers[sub.name] = sub
	i.mx.Unlock()

	return nil
}

// инициализация очереди
func (i *InternalExchange) initNewQueueCons(sub subInfo) *queueConsumer {
	res := initQueueConsumer(sub, i.Id)
	i.activeQueue[sub.GetName()] = res

	return res
}

// Метод удаления подписчика
func (i *InternalExchange) delConsumer(sub subInfo) error {
	name := sub.GetName()

	defer i.mx.Unlock()

	i.mx.Lock()
	_, ok := i.inactiveQueue[name]
	if ok {
		delete(i.inactiveQueue, name)
		return nil
	}

	_, ok = i.activeQueue[name]
	if ok {
		err := fmt.Errorf("the subscriber %s is active", name)
		return err
	}

	return fmt.Errorf("the subscriber %s was not found", name)
}

// Метод перевода подписчика в спящее  состояние.
// В спящем состоянии промежуточный подписчик не перестает сущесвовать, но все активные процессы заканчиваются,
// прекращается отправка событий в цикле этому подписчику, т.е. он удаляется из списка отправки.
func (i *InternalExchange) sleepQueue(nameQueue string) {
	defer log.Infof("the queue: %s has been moved to inactiveQueue", nameQueue)
	queue := i.activeQueue[nameQueue]
	i.inactiveQueue[nameQueue] = queue
	delete(i.activeQueue, nameQueue)
}
