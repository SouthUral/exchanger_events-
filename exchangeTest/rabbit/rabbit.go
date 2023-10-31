package rabbit

import (
	amqp "github.com/rabbitmq/amqp091-go"

	log "github.com/sirupsen/logrus"
)

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
