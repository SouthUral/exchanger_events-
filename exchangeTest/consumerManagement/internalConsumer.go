package consumermanagement

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
type InternalConsumer struct {
	Id                  string                // внутренний ID по которому будет работать маршрутизатор
	mx                  sync.Mutex            // мьютекс, нужен для блокировки доступа
	ConfConsum          configConsum          // конфигурация подписки
	activeSubscribers   map[string]subscriber // активные подписчики
	inactiveSubscribers map[string]subscriber // не активыне подписчики
	IncomingEventsCh    chan interface{}      // канал для получения событий из маршрутизатора
}

// Инициализация внутреннего подписчика
func InitInternalConsumer(conf recipConfig) *InternalConsumer {
	ctx, _ := context.WithCancel(context.Background())

	config := configConsum{
		Types:      conf.GetTypes(),
		Publishers: conf.GetPublihers(),
		AllEvent:   conf.GetAllEvent(),
	}

	intCons := InternalConsumer{
		Id:                  uuid.New().String(),
		ConfConsum:          config,
		mx:                  sync.Mutex{},
		activeSubscribers:   make(map[string]subscriber),
		inactiveSubscribers: make(map[string]subscriber),
		IncomingEventsCh:    make(chan interface{}, 100),
	}

	// запуск маршрутизации событий
	go intCons.eventRouter(ctx)

	log.Debug("инициализирован внутренний потребитель %s", intCons.Id)

	return &intCons
}

// Метод для получения событий из маршрутизатора.
// Получает события и отправляет его всем активным подписчикам и отправляет в канал на сохранение в БД
func (i *InternalConsumer) eventRouter(ctx context.Context) {
	for {
		select {
		case eventMess := <-i.IncomingEventsCh:
			// заглушка
			// здесь должно быть присвоение оффсета входящим событиям, этот оффсет обязательно нужно сохранять куда-то
			// отправка событий всем подписчикам
			// отправка события на запись в БД
			fmt.Println(eventMess)
		case <-ctx.Done():
			return
		}
	}

}

// TODO: Метод добавления подписчика
func (i *InternalConsumer) AddConsumer(sub subInfo) error {
	// defer log.Infof("a subscriber: %s has been added to the queue", sub.name)

	i.mx.Lock()
	_, ok := i.activeSubscribers[sub.GetName()]

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

// Метод удаления подписчика
func (i *InternalConsumer) delConsumer(sub subInfo) error {
	name := sub.GetName()

	defer i.mx.Unlock()

	i.mx.Lock()
	_, ok := i.inactiveSubscribers[name]
	if ok {
		delete(i.inactiveSubscribers, name)
		return nil
	}

	_, ok = i.activeSubscribers[name]
	if ok {
		err := fmt.Errorf("the subscriber %s is active", name)
		return err
	}

	return fmt.Errorf("the subscriber %s was not found", name)
}

// Метод перевода подписчика в спящее  состояние.
// В спящем состоянии промежуточный подписчик не перестает сущесвовать, но все активные процессы заканчиваются,
// прекращается отправка событий в цикле этому подписчику, т.е. он удаляется из списка отправки.
func (i *InternalConsumer) sleepConsumer(sub subscriber) {

}
