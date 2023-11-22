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
	Name          string
	ConfSubscribe ConfSub
	EvenCh        EventChan
}

// Конфигурация подписки
type ConfSub struct {
	Types      []string         `json:"types"`
	Publishers []conf.Publisher `json:"publishers"`
	AllEvent   bool             `json:"all_event"`
}

// Структура которую возвращает инициализатор маршрутизатора (любого типа)
type eventRoutData struct {
	eventCh  EventChan
	subscrCh SubscriberChan
	cancel   func()
}
