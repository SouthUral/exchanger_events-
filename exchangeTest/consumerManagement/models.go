package consumermanagement

import "time"

// сообщение для подписки в маршрутизатор
type subscriptionMess struct {
	nameInternalSub string
	types           []string
	publishers      map[string][]string
	allEvent        bool
	reverseCh       chan interface{}
}

func (s subscriptionMess) GetNameSub() string {
	return s.nameInternalSub
}

func (s subscriptionMess) GetTypes() []string {
	return s.types
}

func (s subscriptionMess) GetPublihers() map[string][]string {
	return s.publishers
}

func (s subscriptionMess) GetAllEvent() bool {
	return s.allEvent
}

func (s subscriptionMess) GetReverseCh() chan interface{} {
	return s.reverseCh
}

// конфигурация подписки
type configConsum struct {
	Types      []string            // типы событий
	Publishers map[string][]string // словарь с отправителями и их типами
	AllEvent   bool                // подписка на все события
}

// сообщение для внутренних потребителей (нужно для подписки пром.потр. на внутреннего потребителя)
type messIntenalCons struct {
	nameCons string
	ConsCh   chan string
}

// сообщение для потребителя
type msgForSub struct {
	err error
}

func (m msgForSub) GetError() error {
	return m.err
}

// событие
type event struct {
	id        string    // id события
	offset    int       // порядковый номер события, присваивается внутренним потребителем
	publisher string    // имя отрпавителя
	typeEvent string    // название типа события
	timePub   time.Time // время отправления (публикации)
	msg       string    // сообщение
}

func (e event) GetOffset() int {
	return e.offset
}

func (e event) GetPub() string {
	return e.publisher
}

func (e event) GetTypeEvent() string {
	return e.typeEvent
}

func (e event) GetTimePub() time.Time {
	return e.timePub
}

func (e event) GetMess() string {
	return e.msg
}

// структура сообщения для модуля работы с хранением информации
type msgForStorageModule struct {
	typeMsg string // тип сообщения (что нужно от модуля)

	reverseCh chan interface{}
}

func (m msgForStorageModule) GetTypeMsg() string {
	return m.typeMsg
}

func (m msgForStorageModule) GetReverseCh() chan interface{} {
	return m.reverseCh
}
