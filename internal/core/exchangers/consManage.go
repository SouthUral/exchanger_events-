package exchangers

import (
	"context"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
)

// TODO: нужно передать сюда состояние, которое будет выгружено из БД при запуске
// либо придет пустая структура (значит что состояния нет в БД), либо придет структура и по ней будут инициализированы внутренние получатели

// TODO: Нужно сделать три канала которые будут переданы сюда от модуля состояния:
// Каналы на отправку:
// 1-й канал - нужен для запросов на восстановление очереди // по каналу отправляется запрос со структурой в которой будет канал для обратного возврата событий
// 2-й канал - нужен для отправки измененного состояния получателей (внутренний и промежуточный) и для отправки данных мониторинга очередей
// 3-й канал - нужен для записи событий в БД

// Каналы на получение:
// TODO: Нужен канал от HTTP API, по которому будут поступать команды для получателей

// Инициализирует менеджера получателей, который управляет внутренними получателями.
// Объект реализует логику, и получает команды от API.
// apiCh: канал, по которому приходят запросы от API;
// subCh: канал, в который отправляются сообщения в роутер, для подписки;
// bdCh: канал для связи с БД, для сохранения сообщений из очереди, сохранения состояния, запроса пропущеных сообщений
func InitConsManager(apiCh <-chan interface{}, subCh chan interface{}, bdCh chan interface{}) func() {
	// здесь происходит общение с API
	// TODO: при загрузке нужно получить конфигурацию для восстановления состояния
	// TODO: либо отправить в модуль работы с БД запрос через канал, на восстановление состояния

	ctx, cancel := context.WithCancel(context.Background()) // зачем этот контекст?

	internalConsumers := initStorageInternalCons(subCh, bdCh)

	go func() {
		for {
			select {
			case mess := <-apiCh:
				// на каждый тип запроса создается своя горутина
				msgApi, ok := mess.(messApi)
				if !ok {
					log.Error("сообщение от модуля API невозможно привести к внутреннему типу модуля consumermanagement")
				}
				switch msgApi.GetTypeRequest() {
				case AddConsumer:
					go internalConsumers.addConsumer(msgApi)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return cancel
}

// // функция добавления подписчика (работает в горутине).
// // функция сама формирует отрицательный ответ и отправляет его API
// func addConsumer(apiCh <-chan interface{}, internalConsumers *storageInternalCons) {
// 	internalConsumers.addConsumer()

// }

// инициализатор storageInternalCons
func initStorageInternalCons(subCh, bdChan chan interface{}) *storageInternalCons {
	storage := make(map[string]*InternalExchange, 20)

	res := &storageInternalCons{
		SubCh:          subCh,
		BdCh:           bdChan,
		storageIntCons: storage,
		mx:             sync.Mutex{},
	}
	return res
}

// Хранилище всех внутренних потребителей
type storageInternalCons struct {
	SubCh          chan interface{}
	BdCh           chan interface{}
	storageIntCons map[string]*InternalExchange
	mx             sync.Mutex
}

// метод для добавления подписчика
func (s *storageInternalCons) addConsumer(msg messApi) {
	key := s.getStorageKey(msg)
	s.mx.Lock()
	internalCons, ok := s.storageIntCons[key]
	s.mx.Unlock()
	if ok {
		err := internalCons.AddConsumer(msg)
		if err != nil {
			subCh := msg.GetReverseCh()
			subCh <- msgForSub{
				err: err,
			}
		}
	} else {
		intCons := s.addNewInternalCons(msg)
		intCons.AddConsumer(msg)
	}
}

// метод добавления нового внутреннего потребителя
func (s *storageInternalCons) addNewInternalCons(conf recipConfig) *InternalExchange {
	defer s.mx.Unlock()
	internalCons := InitInternalExchange(conf)
	s.subscribingRecipientToEvents(internalCons)
	key := s.getStorageKey(conf)
	s.mx.Lock()
	s.storageIntCons[key] = internalCons
	return internalCons
}

// метод для подписки внутреннего получателя на события в маршрутизаторе
func (s *storageInternalCons) subscribingRecipientToEvents(cons *InternalExchange) {
	mess := subscriptionMess{
		nameInternalSub: cons.Id,
		types:           cons.ConfConsum.Types,
		publishers:      cons.ConfConsum.Publishers,
		allEvent:        cons.ConfConsum.AllEvent,
		reverseCh:       cons.IncomingEventsCh,
	}

	s.SubCh <- mess
	log.Debug("сообщение о подписке отправлено в маршрутизатор")
}

// метод формирует ключ для хранилища из сообщения API.
// ключ это конфигурация подписки переведенная в строку
func (s *storageInternalCons) getStorageKey(msg recipConfig) string {
	res := fmt.Sprintf("types: %s, publishers: %s, allEvent: %t",
		msg.GetTypes(),
		msg.GetPublihers(),
		msg.GetAllEvent(),
	)
	return res
}
