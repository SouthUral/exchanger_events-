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
