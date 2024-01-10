package consumermanagement

// Интерфейс входящего сообщения от API
type messApi interface {
	GetTypeRequest() string            // возвращает тип запроса API
	GetReverseCh() chan interface{}    // возвращает канал, по которому нужно будет отдать ответ API
	GetName() string                   // возвращает имя сервиса, который сделал запрос
	GetTypes() []string                // возвращает список типов
	GetPublihers() map[string][]string // возвращает словарь ключ - имя отправителя, значение - список типов
	GetAllEvent() bool                 // флаг, показывает, нужно ли подписать на все события
}

// интерфейс предоставляющий информацию о подписчике
type subInfo interface {
	GetReverseCh() chan interface{} // возвращает канал, по которому нужно будет отдать ответ API
	GetName() string                // возвращает имя сервиса, который сделал запрос
}

// интерфейс для получения конфигурации
type recipConfig interface {
	GetTypes() []string                // возвращает список типов
	GetPublihers() map[string][]string // возвращает словарь ключ - имя отправителя, значение - список типов
	GetAllEvent() bool                 // флаг, показывает, нужно ли подписать на все события
}
