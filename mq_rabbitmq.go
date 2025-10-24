package sw

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("amqp dial: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("channel: %w", err)
	}
	return conn, ch, nil
}

func MQPing(ch *amqp.Channel) error {
	// Declare a passive queue to check connectivity; alternatively tx/select ok
	return ch.Qos(1, 0, false)
}
