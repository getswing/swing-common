package sw

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch *amqp.Channel
}

func NewPublisher(ch *amqp.Channel) *Publisher {
	return &Publisher{ch: ch}
}

// PublishJSON publishes a persistent JSON message to a named queue
func (p *Publisher) PublishJSON(ctx context.Context, queue string, payload any) error {
	_, err := p.ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	pubCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return p.ch.PublishWithContext(pubCtx, "", queue, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
		DeliveryMode: amqp.Persistent,
	})
}
