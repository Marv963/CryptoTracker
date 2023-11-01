package messaging

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/Marv963/CryptoTracker/app/pkg/config"
)

type RabbitMQ struct {
	Channel  *amqp.Channel
	Exchange string
}

func NewRabbitMQ(cfg config.RabbitMQ) (*RabbitMQ, error) {
	connString := fmt.Sprintf("amqp://%s:%s@%s/", cfg.Username, cfg.Password, cfg.URL)
	conn, err := amqp.Dial(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %s", err)
	}

	err = ch.Qos(
		20,    // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, fmt.Errorf("failed to set QoS: %s", err)
	}

	return &RabbitMQ{Channel: ch, Exchange: cfg.Exchange}, nil
}

func (r *RabbitMQ) Receive(queue string) (<-chan amqp.Delivery, error) {
	msgs, err := r.Channel.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer: %s", err)
	}

	return msgs, nil
}

func (r *RabbitMQ) BindQueue(exchange string, queue string, routingKey string) error {
	err := r.Channel.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare an exchange: %s", err)
	}

	_, err = r.Channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %s", err)
	}

	err = r.Channel.QueueBind(
		queue,      // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind a queue: %s", err)
	}

	return nil
}
