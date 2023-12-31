package utilsfortest

import (
	"time"
)

// Структура для парсинга json
type Config struct {
	Publishers []Publisher `json:"publishers"`
	Consumers  []Consumer  `json:"consumers"`
}

type Consumer struct {
	Name       string      `json:"name"`
	Publishers []Publisher `json:"publishers"`
	Types      []string    `json:"type_mess"`
	AllEvent   bool        `json:"all_event"`
}

type Publisher struct {
	Name     string   `json:"name"`
	TypeMess []string `json:"type_mess"`
}

// структура события
type event struct {
	id          int
	publisher   string
	typeEvent   string
	timePublish time.Time
	message     string
}

func (e event) GetId() int {
	return e.id
}

func (e event) GetPub() string {
	return e.publisher
}

func (e event) GetTypeEvent() string {
	return e.typeEvent
}

func (e event) GetTimePub() time.Time {
	return e.timePublish
}

func (e event) GetMess() string {
	return e.message
}

// реализация интерфейса SubscriberMess
type subMess struct {
	Name       string
	Types      []string
	Publishers map[string][]string
	AllEvent   bool
	ReverseCh  chan interface{}
}

func (s subMess) GetNameSub() string {
	return s.Name
}

func (s subMess) GetTypes() []string {
	return s.Types
}

func (s subMess) GetReverseCh() chan interface{} {
	return s.ReverseCh
}

func (s subMess) GetPublihers() map[string][]string {
	return s.Publishers
}

func (s subMess) GetAllEvent() bool {
	return s.AllEvent
}

func initSubMess(subData Consumer) subMess {
	selfEventCh := make(chan interface{}, 500)

	publishers := make(map[string][]string, 0)
	for _, item := range subData.Publishers {
		publishers[item.Name] = item.TypeMess
	}

	res := subMess{
		Name:       subData.Name,
		Publishers: publishers,
		Types:      subData.Types,
		AllEvent:   subData.AllEvent,
		ReverseCh:  selfEventCh,
	}

	return res
}
