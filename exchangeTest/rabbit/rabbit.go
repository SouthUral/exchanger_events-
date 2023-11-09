package rabbit

import (
	"context"

	// router "github.com/SouthUral/exchangeTest/router"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

// TODO: Работа консюмера (подписка на очередь)
// TODO: Работа паблишера (отправка сообщения в очередь)

// func Publisher(URL string, exchangeName string, namePub string, eventChan router.EventChan) {

// 	go func() {
// 		connPub, err := createConnRabbit(URL)

// 		if err != nil {
// 			log.Errorf("Отправитель Rabbit %s прекратил свою работу, ошибка коннекта", namePub)
// 			return
// 		}

// 		chanPub, err := createChannRabbit(connPub)
// 		if err != nil {
// 			log.Errorf("Отправитель Rabbit %s прекратил свою работу, ошибка создания канала", namePub)
// 			return
// 		}

// 		for {
// 			select {
// 			case event := <-eventChan:
// 				// TODO: какой routing_key нужно передать в Rabbit, есть два основных параметра: typeEvent, Publisher
// 				// из них надо собрать routing_key
// 				sendingMess(chanPub, context.Background(), exchangeName)
// 			}
// 		}
// 		sendingMess()
// 	}()

// }

// Функция создания подключения к RabbitMQ
func createConnRabbit(URL string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(URL)
	if err != nil {
		log.Error(err.Error())
	}
	return conn, err
}

// Функция создания канала для работы с RabbitMQ
func createChannRabbit(conn *amqp.Connection) (*amqp.Channel, error) {
	rabbitChan, err := conn.Channel()
	if err != nil {
		log.Error(err.Error())
	}
	return rabbitChan, err
}

// Функция декларирования Exchange
func declaringEchange(rabbitChan *amqp.Channel, nameExchange string, typeExchange string) error {
	err := rabbitChan.ExchangeDeclare(
		nameExchange,
		typeExchange,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

// Функция декларирования очереди
func declaringQueue(rabbitChan *amqp.Channel, nameQueue string) (amqp.Queue, error) {
	queue, err := rabbitChan.QueueDeclare(
		nameQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Errorf("Ошибка создания очереди: %s", err.Error())
	}
	return queue, err
}

// Функция binding очереди к Exchange
func bindingQueue(rabbitChan *amqp.Channel, nameQueue, nameExchange, routingKey string) error {
	err := rabbitChan.QueueBind(
		nameQueue,
		routingKey,
		nameExchange,
		false,
		nil,
	)
	if err != nil {
		log.Errorf("Ошибка binding очереди %s, к exchange %s : %s", nameQueue, nameExchange, err.Error())
	}
	return err
}

// Функция отправки сообщения в Exchange
func sendingMess(chanRabbit *amqp.Channel, ctx context.Context, exchangeName, routingKey, mess string) error {
	err := chanRabbit.PublishWithContext(
		ctx,
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(mess),
		},
	)

	if err != nil {
		log.Errorf("Ошибка отправки сообщения: %s", err.Error())
	}
	log.Debugf("Сообщение отправлено: %s", mess)
	return err
}

// Функция закрытия канала и коннекта RabbitMQ
func exitAndClose(chanRb *amqp.Channel, connRb *amqp.Connection, whoIsIt string) {
	errChan := chanRb.Close()
	if errChan != nil {
		log.Errorf("Ошибка закрытия канала RabbitMQ: %s", whoIsIt)
	}

	errConn := connRb.Close()
	if errConn != nil {
		log.Errorf("Ошибка закрытия коннекта RabbitMQ: %s", whoIsIt)
	}

	log.Warningf("Закрыты канал и коннект к RabbitMQ: %s", whoIsIt)
}
