package queue

import (
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitMQProducer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func NewRabbitMQProducer(url string, queueName string) (*RabbitMQProducer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName, // nombre de la cola
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	return &RabbitMQProducer{conn: conn, channel: ch, queue: queueName}, nil
}

func (p *RabbitMQProducer) Publish(body []byte) error {
	// Verificación de seguridad
	if p == nil || p.channel == nil {
		return fmt.Errorf("el productor o el canal de RabbitMQ no están inicializados")
	}

	return p.channel.Publish(
		"",      // exchange
		p.queue, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
