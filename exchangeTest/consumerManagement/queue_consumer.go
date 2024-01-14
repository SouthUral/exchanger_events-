package consumermanagement

import (
	"context"
	"sync"

	deq "github.com/gammazero/deque"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// Промежуточный потребитель или очередь конкретного подписчика
// Может иметь два состояние, активное/неактивное/сон, в неактичном состоянии просто сохраняет события в очереди

// TODO: общая структура в котой содержится вся текущая информация
// TODO: горутина очереди, в которой происходит загрузка и выгрузка данных из очереди
// TODO: горутина отправки, горутина отправляет события внешнему получателю

// структура очереди потребителя
type queueConsumer struct {
	idQueue             string           // уникальный индентификатор, необходим для восстановления состояния программы
	idExchange          string           // id InternalExchange
	name                string           // имя, соответствует имени внешнего потребителя
	limitSizeQueue      int              // предельное количество событий в очереди, по достижению которого в режиме inactive очередь сбросится и перейдет в спящий режим
	currentStatus       string           // текущее состояние (active/inactive)
	inputCh             chan event       // канал, по которому приходят события
	mx                  sync.RWMutex     // мьютекс для конкурентного доступа к атрибутам
	outputCh            chan interface{} // канал, по которому события отправляются внешнему потребителю
	lastMsg             event            // последнее отпрвленное (или неотправленное событие, пока непонятно)
	lastdispatchedEvent string
	storageModuleCh     chan interface{} // канал для связи с модулей хранения
}

func initQueueConsumer(subInfo subInfo, idExchange string) *queueConsumer {

	res := &queueConsumer{
		idQueue:        uuid.New().String(),
		idExchange:     idExchange,
		name:           subInfo.GetName(),
		currentStatus:  statusActive,
		inputCh:        make(chan event, 100),
		outputCh:       subInfo.GetReverseCh(),
		mx:             sync.RWMutex{},
		limitSizeQueue: subInfo.GetMaxSize(),
	}

	log.Infof("Created queue for consumer %s", subInfo.GetName())

	go res.Queue(subInfo.GetCtx())

	return res
}

func (q *queueConsumer) Queue(ctx context.Context) {
	var que deq.Deque[event]
	context, cancel := context.WithCancel(context.Background())

	maxSizeQueue := q.limitSizeQueue
	nameQueue := q.name
	idExchange := q.idExchange

	// запуск горутины чтения из очереди и отправки события внешнему потребителю
	go func() {
		defer log.Infof("loss of connection with an external recipient queue %s:%s", idExchange, nameQueue)

		for {
			select {
			case <-ctx.Done():
				// отмена контекста внешнего получателя, переход в ждущий режим
				q.mx.Lock()
				// перевод очереди в неактивный режим, т.е. очередь все еще работает, но теперь по достижении лимита очередь уснет
				q.currentStatus = statusInActive
				q.mx.Unlock()
				return
			default:
				if que.Len() > 0 {
					msg := que.PopFront()
					q.outputCh <- msg

					q.mx.Lock()
					q.lastMsg = msg
					q.mx.Unlock()
				}

			}
		}
	}()

	// горутина чтения событий из exhange, записывает полученные события в очередь
	// в активном состоянии запись в очередь ведется без ограничений
	// в неактивном состоянии запись ведется только до достижения лимита
	go func() {
		defer log.Infof("the queue %s:%s has stopped recording", idExchange, nameQueue)

		for {
			select {
			case <-context.Done():
				return
			case msg := <-q.inputCh:
				q.mx.RLock()
				status := q.currentStatus
				q.mx.RUnlock()

				switch status {
				case statusActive:
					que.PushBack(msg)
				case statusInActive:
					if que.Len() < maxSizeQueue {
						que.PushBack(msg)
					} else {
						// прекращение работы очереди, переход в спящий режим
						q.switchingToSleep()
						cancel()
					}
				}
			}
		}

	}()

}

// метод: перевод intermediateConsumer в спящий режим.
// Отключение всех активных горутин, после перевода по inputCh не должны поступать события
func (q *queueConsumer) switchingToSleep() {
}
