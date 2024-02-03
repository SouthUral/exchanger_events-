package exchangers

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// Внутренний потребитель, по факту это внутрення очередь, получает события и присваивает им offset.
// Отдает события промежуточному потребителю.
// Отправляет события на сохранения в БД.
// структура внутреннего подписчика
type InternalExchange struct {
	Id               string                    // внутренний ID по которому будет работать маршрутизатор
	mx               sync.Mutex                // мьютекс, нужен для блокировки доступа
	ConfConsum       configConsum              // конфигурация подписки
	activeQueue      map[string]*queueConsumer // активные подписчики
	inactiveQueue    map[string]*queueConsumer // не активыне подписчики
	IncomingEventsCh chan interface{}          // канал для получения событий из маршрутизатора
	lastOffset       int
}

// Инициализация внутреннего подписчика
func InitInternalExchange(conf recipConfig) *InternalExchange {
	ctx, _ := context.WithCancel(context.Background())

	config := configConsum{
		Types:      conf.GetTypes(),
		Publishers: conf.GetPublihers(),
		AllEvent:   conf.GetAllEvent(),
	}

	intCons := InternalExchange{
		Id:               uuid.New().String(),
		ConfConsum:       config,
		mx:               sync.Mutex{},
		activeQueue:      make(map[string]*queueConsumer),
		inactiveQueue:    make(map[string]*queueConsumer),
		IncomingEventsCh: make(chan interface{}, 100),
	}

	// запуск маршрутизации событий
	go intCons.eventRouter(ctx)

	log.Debug("инициализирован внутренний потребитель %s", intCons.Id)

	return &intCons
}

// Метод для получения событий из маршрутизатора.
// Получает события и отправляет его всем активным подписчикам и отправляет в канал на сохранение в БД
func (i *InternalExchange) eventRouter(ctx context.Context) {
	defer log.Infof("the event router has stopped working")
	for {
		select {
		case eventMess := <-i.IncomingEventsCh:
			event := i.eventConversion(eventMess)
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
func (i *InternalExchange) eventConversion(msg interface{}) event {
	eventInterface, ok := msg.(eventFromRouter)
	if !ok {
		log.Fatal("it is not possible to convert an event to an interface")
	}

	offset := i.lastOffset + 1

	newEvent := event{
		offset:    offset,
		publisher: eventInterface.GetPub(),
		typeEvent: eventInterface.GetTypeEvent(),
		timePub:   eventInterface.GetTimePub(),
		msg:       eventInterface.GetMess(),
	}

	i.lastOffset = offset

	return newEvent
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
