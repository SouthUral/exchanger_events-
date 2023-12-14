package consumermanagement

import (
	models "github.com/SouthUral/exchangeTest/models"
)

// структура содержит информацию о внутреннем потребителе
type internalConsum struct {
	conf           models.ConfSub
	internalConsCh chan messIntenalCons
	consumers      map[string]infoQueueConsumer
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
	name  string
	subCh models.RevSubCh
}
