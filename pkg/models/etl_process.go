package models

import (
	"fmt"
	"formdata/pkg/messagequeues"
)

type ETLProcess struct {
	ConnectionID string
	Name         string
	MqClient     *messagequeues.RabbitMQClient
	QueueName    string
	Config       Config
}

func (e *ETLProcess) SendEndMessage() error {
	ch, err := e.MqClient.NewChannel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %v", err)
	}
	defer ch.Close()

	endMessage := messagequeues.ETLMessage{
		IsEnd: true,
	}

	if err := e.MqClient.Publish(ch, e.QueueName, endMessage); err != nil {
		return fmt.Errorf("failed to publish end message: %v", err)
	}

	return nil
}

func (e *ETLProcess) Cleanup() {
	if e.MqClient != nil {
		e.MqClient.Close()
	}
}
