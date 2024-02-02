package utilsfortest

import (
	"time"
)

// Интерфейс события
type Event interface {
	GetId() int
	GetPub() string
	GetTypeEvent() string
	GetTimePub() time.Time
	GetMess() string
}

// Интерфейс конфигурации подписки потребителя
// type ConfSub interface {
// 	GetTypes() []string
// 	GetPublihers() map[string][]string
// 	GetAllEvent() bool
// }

// Интерфейс сообщения от подписчика, сообщение подписки
type SubscriberMess interface {
	GetNameSub() string
	GetTypes() []string
	GetPublihers() map[string][]string
	GetAllEvent() bool
	GetReverseCh() chan Event
}
