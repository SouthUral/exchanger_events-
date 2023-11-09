package router

import (
	"time"

	conf "github.com/SouthUral/exchangeTest/confreader"
)

// Канал для отправки событий
type EventChan chan Event

// Событие, которое отправляет publisher
type Event struct {
	Id          int
	Publisher   string
	TypeEvent   string
	TimePublish time.Time
	Message     string
}

// Канал для отправки подписчиков
type SubscriberChan chan SubscriberMess

// Сообщение содержащее информацию о подписчике
type SubscriberMess struct {
	Name       string
	Types      []string
	Publishers []conf.Publisher
	AllEvent   bool
	EvenCh     EventChan
}

// Структура которую возвращает инициализатор маршрутизатора (любого типа)
type eventRoutData struct {
	eventCh  EventChan
	subscrCh SubscriberChan
	cancel   func()
}

// type Publisher struct {
// 	name  string
// 	types []string
// }
