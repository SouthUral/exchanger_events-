package router

import (
	models "github.com/SouthUral/exchangeTest/models"
)

// Структура которую возвращает инициализатор маршрутизатора (любого типа)
type eventRoutData struct {
	eventCh  models.EventChan
	subscrCh models.SubscriberChan
	cancel   func()
}
