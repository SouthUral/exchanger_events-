package apiv1

import (
	rt "github.com/SouthUral/exchangeTest/router"
)

// Канал по которому API gRPC отправляют события
type EventAPICh chan EventApi

// Cтруктура содержит сообщение от API.
type EventApi struct {
	RevСh ReverseCh
	Event rt.Event
}

// Канал для ответного сообщения в API gRPC
type ReverseCh chan ReverseMess

// Ответное сообщение для API
type ReverseMess struct {
	IdEvent int
	Err     error
}

// Сообщения от SubRouter менеджеру потребителей
type SubMess struct {
	TypeMess    string            // тип сообщения (может быть сообщение об остановке, либо просто регистрация нового потребителя)
	Conf        rt.SubscriberMess // конфигурация по которой будет производится подписка в маршрутизаторе
	ChooseEvent string            // либо послденее непрочитанное либо текущее
	RevCh       RevSubCh          // канал по которому возвращаются события
}

// Канал для обратной связи, по которому будут возвращаться события от маршрутизатора
type RevSubCh chan RevMessSub

// Сообщение которое получит SubRouter в ответ на свой запрос
type RevMessSub struct {
	Event rt.Event
	Err   error
}
