package router

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
type ConfSub interface {
	GetTypes() []string
	GetPublihers() map[string][]string
	GetPub(namePub string) (bool, []string)
	GetAllEvent() bool
}

// Интерфейс сообщения от подписчика, сообщение подписки
type SubscriberMess interface {
	GetNameSub() string
	GetConfigSub() ConfSub
	GetReverseCh() chan interface{}
}

type ExternalSubMess interface {
	GetNameSub() string
	GetTypes() []string
	GetPublihers() map[string][]string
	GetAllEvent() bool
	GetReverseCh() chan interface{}
}
