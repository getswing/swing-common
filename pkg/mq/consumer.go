package common_mq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func Consume(ch *amqp.Channel, queue string) (<-chan amqp.Delivery, error) {
	_, err := ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return ch.Consume(queue, "", false, false, false, false, nil)
}
