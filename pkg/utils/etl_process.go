package utils

import (
	"fmt"
	"formdata/pkg/messagequeues"
	"formdata/pkg/models"
	"os"
)

func NewETLProcess(name string) (*models.ETLProcess, error) {
	connectionID := os.Getenv("CONNECTION_ID")
	if connectionID == "" {
		return nil, fmt.Errorf("CONNECTION_ID environment variable not set")
	}
	queueName := os.Getenv("QUEUE_NAME")
	if queueName == "" {
		return nil, fmt.Errorf("QUEUE_NAME environment variable not set")
	}

	config, err := ConfigFromEnv(name)
	if err != nil {
		return nil, err
	}

	mqClient, err := messagequeues.New()
	if err != nil {
		return nil, err
	}

	return &models.ETLProcess{
		ConnectionID: connectionID,
		Name:         name,
		MqClient:     mqClient,
		QueueName:    queueName,
		Config:       config,
	}, nil
}
