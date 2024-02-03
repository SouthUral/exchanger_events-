package exchangers

import (
	"context"
	"time"
)

// Интерфейс входящего сообщения от API
type messApi interface {
	GetTypeRequest() string            // возвращает тип запроса API
	GetReverseCh() chan interface{}    // возвращает канал, по которому нужно будет отдать ответ API
	GetName() string                   // возвращает имя сервиса, который сделал запрос
	GetTypes() []string                // возвращает список типов
	GetPublihers() map[string][]string // возвращает словарь ключ - имя отправителя, значение - список типов
	GetAllEvent() bool                 // флаг, показывает, нужно ли подписать на все события
	GetCtx() context.Context           // контекст API
	GetMaxSize() int                   // предельный размер очереди в неактивном состоянии
	GetTimeLastEvent() time.Time       // время последнего полученного события, игнорируется при первом запуске очереди
	GetOffsetLastEvent() int           // offset последнего полученного события, игнорируется при первом запуске очереди
	GetRestoringQueue() bool           // флаг, восстанавливать ли очередь после переподключения, игнорируется при первом запуске очереди
}

// интерфейс предоставляющий информацию о подписчике
type subInfo interface {
	GetReverseCh() chan interface{} // возвращает канал, по которому нужно будет отдать ответ API
	GetName() string                // возвращает имя сервиса, который сделал запрос
	GetCtx() context.Context        // контекст API
	GetMaxSize() int                // предельный размер очереди в неактивном состоянии
	GetTimeLastEvent() time.Time    // время последнего полученного события, игнорируется при первом запуске очереди
	GetOffsetLastEvent() int        // offset последнего полученного события, игнорируется при первом запуске очереди
	GetRestoringQueue() bool        // флаг, восстанавливать ли очередь после переподключения, игнорируется при первом запуске очереди
}

// интерфейс для получения конфигурации
type recipConfig interface {
	GetTypes() []string                // возвращает список типов
	GetPublihers() map[string][]string // возвращает словарь ключ - имя отправителя, значение - список типов
	GetAllEvent() bool                 // флаг, показывает, нужно ли подписать на все события
	GetMaxSize() int                   // предельный размер очереди в неактивном состоянии
}

// Интерфейс события от маршрутизатора
type eventFromRouter interface {
	// GetId() int
	GetPub() string        // возвращает имя сервиса отправителя
	GetTypeEvent() string  // возвращает имя типа события
	GetTimePub() time.Time // время отправки события (пока непонятно откуда событие было отправлено)
	GetMess() string       // Возвращает сообщение события
}
