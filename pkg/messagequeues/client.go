package messagequeues

type MessageQueueClient interface {
	Publish(queueName string, msg ETLMessage) error
	Consume(queueName string) (<-chan ETLMessage, error)
	Close() error
}
