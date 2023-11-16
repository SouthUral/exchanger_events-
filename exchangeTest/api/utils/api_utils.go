package utils

import (
	"fmt"
	"time"

	pb "github.com/SouthUral/exchangeTest/grpc"
	rt "github.com/SouthUral/exchangeTest/router"
)

// Преобразует событие gRPC во внутренне событие
func ConvertEvent(event *pb.Event) (rt.Event, error) {
	var newEvent rt.Event

	err := checkEvent(event)
	if err != nil {
		return newEvent, err
	}

	newEvent = rt.Event{
		Publisher:   event.Publisher,
		TypeEvent:   event.TypeEvent,
		TimePublish: time.Now(),
		Message:     event.Message,
	}
	return newEvent, err
}

// Функция для проверки заполенности pb.Event
func checkEvent(event *pb.Event) error {
	var err error
	if event.Publisher == "" {
		err = fmt.Errorf("the sender of the event is not defined")
		return err
	}

	if event.TypeEvent == "" {
		err = fmt.Errorf("the event type is not defined")
		return err
	}

	return err
}
