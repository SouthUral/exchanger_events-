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
	Types      []string    `json:"types"`
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
	Name      string
	Config    confSub
	ReverseCh chan interface{}
}

func (s subMess) GetNameSub() string {
	return s.Name
}

func (s subMess) GetConfigSub() ConfSub {
	return s.Config
}

func (s subMess) GetReverseCh() chan interface{} {
	return s.ReverseCh
}

func initSubMess(subData Consumer) subMess {
	selfEventCh := make(chan interface{}, 100)

	config := initConSub(subData)

	res := subMess{
		Name:      subData.Name,
		Config:    config,
		ReverseCh: selfEventCh,
	}

	return res
}

// структура конфигурации подписчика имплементирующая интерфейс ConfSub
type confSub struct {
	Types      []string
	Publishers map[string][]string
	AllEvent   bool
}

func (c confSub) GetTypes() []string {
	return c.Types
}

func (c confSub) GetPublihers() map[string][]string {
	return c.Publishers
}

func (c confSub) GetPub(namePub string) (bool, []string) {
	pub, ok := c.Publishers[namePub]
	return ok, pub
}

func (c confSub) GetAllEvent() bool {
	return c.AllEvent
}

func initConSub(cons Consumer) confSub {
	publishers := make(map[string][]string, 0)
	for _, item := range cons.Publishers {
		publishers[item.Name] = item.TypeMess
	}

	result := confSub{
		Types:      cons.Types,
		Publishers: publishers,
		AllEvent:   cons.AllEvent,
	}

	return result
}
