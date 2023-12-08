package models

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
)

// Сообщение содержащее информацию о подписчике
type SubscriberMess struct {
	Name          string
	ConfSubscribe ConfSub
	EvenCh        EventChan
}

type Publisher struct {
	Name     string   `json:"name"`
	TypeMess []string `json:"type_mess"`
}

// Конфигурация подписки потребителя
type ConfSub struct {
	Types      []string    `json:"types"`
	Publishers []Publisher `json:"publishers"`
	AllEvent   bool        `json:"all_event"`
}

func (conf *ConfSub) GenStringVue() string {
	confByte, err := json.Marshal(conf)
	if err != nil {
		log.Error(err)
		return ""
	}

	return string(confByte)
}

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
