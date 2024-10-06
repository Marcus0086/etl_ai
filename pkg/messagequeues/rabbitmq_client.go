package messagequeues

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct {
	connection *amqp.Connection
}

func (r *RabbitMQClient) NewChannel() (*amqp.Channel, error) {
	return r.connection.Channel()
}

func New() (*RabbitMQClient, error) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		return nil, err
	}

	return &RabbitMQClient{
		connection: conn,
	}, nil
}

func (r *RabbitMQClient) Publish(channel *amqp.Channel, queueName string, msg ETLMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	return channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (r *RabbitMQClient) Consume(channel *amqp.Channel, queueName string) (<-chan ETLMessage, error) {
	_, err := channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}
	deliveries, err := channel.Consume(
		queueName, // queue
		"",        // consumer tag
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}

	msgs := make(chan ETLMessage, 1)

	go func() {
		defer close(msgs)
		for delivery := range deliveries {
			var msg ETLMessage
			if err := json.Unmarshal(delivery.Body, &msg); err != nil {
				continue
			}
			msgs <- msg
		}
	}()

	return msgs, nil

}

func (c *RabbitMQClient) Close() error {
	return c.connection.Close()
}
