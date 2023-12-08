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

// Канал по которому API gRPC отправляют события
type EventAPICh chan EventApi

// Канал для ответного сообщения в API gRPC
type ReverseCh chan ReverseMess

// Канал для обратной связи, по которому будут возвращаться события от маршрутизатора
type RevSubCh chan RevMessSub

// Сообщение которое получит SubRouter в ответ на свой запрос
type RevMessSub struct {
	Event Event
	Err   error
}

// Cтруктура содержит сообщение от API.
type EventApi struct {
	RevСh ReverseCh
	Event Event
}

// Ответное сообщение для API
type ReverseMess struct {
	IdEvent int
	Err     error
}

// Сообщения от SubRouter менеджеру потребителей
type SubMess struct {
	TypeMess    string         // тип сообщения (может быть сообщение об остановке, либо просто регистрация нового потребителя)
	Conf        SubscriberMess // конфигурация по которой будет производится подписка в маршрутизаторе
	ChooseEvent string         // либо послденее непрочитанное либо текущее
	RevCh       RevSubCh       // канал по которому возвращаются события
}
