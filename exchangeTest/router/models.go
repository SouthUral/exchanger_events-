package router

// Структура которую возвращает инициализатор маршрутизатора (любого типа)
type routData struct {
	eventCh  chan Event          // Канал для приема/передачи событий
	subscrCh chan SubscriberMess // Канал для примема/передачи сведений о подписчике
}

// реализация интерфейса SubscriberMess
type subMess struct {
	name      string
	config    ConfSub
	reverseCh chan Event
}

func (s *subMess) GetNameSub() string {
	return s.name
}

func (s *subMess) GetConfigSub() ConfSub {
	return s.config
}

func (s *subMess) GetReverseCh() chan Event {
	return s.reverseCh
}

// реализация интерфейса ConfSub
type confSub struct {
	types      []string
	publishers map[string][]string
	allEvent   bool
}

func (c *confSub) GetTypes() []string {
	return c.types
}

func (c *confSub) GetPublihers() map[string][]string {
	return c.publishers
}

func (c *confSub) GetPub(namePub string) (bool, []string) {
	publisher, ok := c.publishers[namePub]
	return ok, publisher
}

func (c *confSub) GetAllEvent() bool {
	return c.allEvent
}
