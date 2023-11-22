package consumermanagement

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	rt "github.com/SouthUral/exchangeTest/router"
)

// Внутренний потребитель, по факту это внутрення очередь, получает события и присваивает им offset.
// Отдает события промежуточному потребителю.
// Отправляет события на сохранения в БД.
// структура внутреннего подписчика
type InternalConsumer struct {
	ConfConsum       rt.ConfSub            // конфигурация подписчика
	Id               string                // внутренний ID по которому будет работать маршрутизатор
	mx               sync.Mutex            // мьютекс, нужен для блокировки доступа
	Subscribers      map[string]subscriber // map с подписчиками
	SubChan          rt.SubscriberChan     // канал, по которому идет отправка информации о подписке на события
	IncomingEventsCh rt.EventChan
}

// Инициализация внутреннего подписчика
// TODO: для запуска нужен конфиг и канал для передачи сообщения подписки
func InitInternalConsumer(conf rt.ConfSub, subChan rt.SubscriberChan) *InternalConsumer {
	intCons := InternalConsumer{
		Id:               uuid.New().String(),
		ConfConsum:       conf,
		mx:               sync.Mutex{},
		Subscribers:      make(map[string]subscriber),
		SubChan:          subChan,
		IncomingEventsCh: make(rt.EventChan, 100),
	}

	intCons.SubscribingEvents()

	return &intCons
}

// метод для подписки на события в маршрутизаторе
func (cons *InternalConsumer) SubscribingEvents() {
	defer log.Info("the subscription message has been sent to the router")
	cons.SubChan <- rt.SubscriberMess{
		Name:          cons.Id,
		ConfSubscribe: cons.ConfConsum,
		EvenCh:        cons.IncomingEventsCh,
	}
}

// Метод запускаемы как горутина, прослушивает входящие событие от роута.
// Прослушивает канал для получения/удаления новых подписчиков
// Должен иметь доступ к каналу событий и каналу получения новых подписчиков
func (cons *InternalConsumer) mainInternalCons() {
	for {
		select {
		case eventMess := <-cons.IncomingEventsCh:
			// заглушка
			// здесь должно быть присвоение оффсета входящим событиям, этот оффсет обязательно нужно сохранять куда-то
			// отправка событий всем подписчикам
			// отправка события на запись в БД
			fmt.Println(eventMess)
		}
	}

}

// TODO: Метод добавления подписчика
func (cons *InternalConsumer) AddConsumer(sub subscriber) error {
	defer log.Infof("a subscriber: %s has been added to the queue", sub.name)

	var err error

	cons.mx.Lock()
	_, ok := cons.Subscribers[sub.name]
	if ok {
		err = fmt.Errorf("")
	}
	// нужно проверить активность подписчика
	// если подписчик активен, то вернуть ошибку
	// если подписчик неактивен, то добавить ему новый канал и включить активность
	// если подписчика нет, то нужно создать его и добавить в подписчики
	//

	cons.Subscribers[sub.name] = sub
	cons.mx.Unlock()
}

// TODO: Метод удаления подписчика
func (cons *InternalConsumer) delConsumer(sub subscriber) {
	defer log.Infof("a subscriber: %s has been deleted from the queue", sub.name)

	cons.mx.Lock()
	delete(cons.Subscribers, sub.name)
	cons.mx.Unlock()
}

// TODO: горутна прослушивания событий

// TODO: нужен метод подписки на события в роутере

// TODO: Метод отписки от всех событий в роутере
