package consumermanagement

// структура содержит информацию о внутреннем потребителе
type internalConsum struct {
	// conf           models.ConfSub
	internalConsCh chan messIntenalCons
	consumers      map[string]infoQueueConsumer
}

// сообщение для подписки в маршрутизатор
type subscriptionMess struct {
	nameInternalSub string
	types           []string
	publishers      map[string][]string
	allEvent        bool
	reverseCh       chan interface{}
}

func (s subscriptionMess) GetNameSub() string {
	return s.nameInternalSub
}

func (s subscriptionMess) GetTypes() []string {
	return s.types
}

func (s subscriptionMess) GetPublihers() map[string][]string {
	return s.publishers
}

func (s subscriptionMess) GetAllEvent() bool {
	return s.allEvent
}

func (s subscriptionMess) GetReverseCh() chan interface{} {
	return s.reverseCh
}

// конфигурация подписки
type configConsum struct {
	Types      []string            // типы событий
	Publishers map[string][]string // словарь с отправителями и их типами
	AllEvent   bool                // подписка на все события
}

// сообщение для внутренних потребителей (нужно для подписки пром.потр. на внутреннего потребителя)
type messIntenalCons struct {
	nameCons string
	ConsCh   chan string
}

// Информация о пром.потр. (очереди)
type infoQueueConsumer struct {
	isActive  bool
	queueChan chan string
}

// структура содержащая информацию о подписчике и канале связи с ним
type subscriber struct {
	name string
	// subCh models.RevSubCh
}

// сообщение для потребителя
type msgForSub struct {
	err error
}

func (m msgForSub) GetError() error {
	return m.err
}
