package apiv1

import (
	"context"
	"time"

	ut "github.com/SouthUral/exchangeTest/api/utils"
	pb "github.com/SouthUral/exchangeTest/grpc"

	models "github.com/SouthUral/exchangeTest/models"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func initServer() server {
	evtCh := make(models.EventAPICh, 100)

	srv := server{
		eventCh: evtCh,
	}

	return srv
}

type server struct {
	eventCh models.EventAPICh
}

// Логика приема события от отправителя
func (s *server) SendingEvent(ctx context.Context, in *pb.Event, opts ...grpc.CallOption) (*pb.EventID, error) {
	defer log.Debug("a reply has been sent to the sender")

	event, err := ut.ConvertEvent(in)

	answer := pb.EventID{}

	if err != nil {
		return &answer, err
	}

	revCh := make(models.ReverseCh)
	mess := models.EventApi{
		RevСh: revCh,
		Event: event,
	}

	s.eventCh <- mess

	// ожидание ответа
	for {
		select {
		case answMess := <-revCh:
			if answMess.Err != nil {
				return &answer, answMess.Err
			}
			answer.Id = answer.Id
			return &answer, err
		default:
			time.Sleep(5 * time.Millisecond)
			continue
		}
	}
}

func (s *server) SubscribingEvents(ctx context.Context, in *pb.ConsumerData, opts ...grpc.CallOption) (pb.ExchangeEvents_SubscribingEventsClient, error) {

}
