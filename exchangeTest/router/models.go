package router

import (
	"time"
)

// Канал для отправки событий
type EventChan chan Event

// Событие, которое отправляет publisher
type Event struct {
	id          int
	publisher   string
	typeEvent   string
	timePublish time.Time
	message     string
}

// Канал для отправки подписчиков
type SubscriberChan chan SubscriberMess

// Сообщение содержащее информацию о подписчике
type SubscriberMess struct {
	name       string
	types      []string
	publishers []publisher
	allEvent   bool
	evenCh     EventChan
}

// Структура содержащая каналы маршрутизаторов
type routersChans struct {
	typeCh    chan Event
	publishCh chan Event
	allCh     chan Event
}

// Структура которую возвращает инициализатор маршрутизатора (любого типа)
type eventRoutData struct {
	eventCh  EventChan
	subscrCh SubscriberChan
	cancel   func()
}

type publisher struct {
	name  string
	types []string
}

type publisherRoutData struct {
	eventCh  EventChan
	subscrCh SubscriberChan
	cancel   func()
}
