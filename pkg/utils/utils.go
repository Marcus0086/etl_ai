package utils

import (
	"formdata/pkg/messagequeues"
	"formdata/pkg/wokrerpool"
	"log"
	"time"
)

func CreateWorkerFunc(
	mqClient *messagequeues.RabbitMQClient,
	queueName,
	source,
	destination string,
) wokrerpool.WorkerFunc[[]byte] {
	return func(workerID int, chunk []byte) error {
		ch, err := mqClient.NewChannel()
		if err != nil {
			log.Printf("Worker %d: Failed to create channel: %v", workerID, err)
			return err
		}
		defer ch.Close()

		msg := messagequeues.ETLMessage{
			Data: chunk,
			MetaData: messagequeues.MetaData{
				Source:      source,
				Destination: destination,
			},
			CreatedAt: time.Now(),
			UpdateAt:  time.Now(),
		}
		err = mqClient.Publish(ch, queueName, msg)
		if err != nil {
			log.Printf("Worker %d: Failed to publish message: %v", workerID, err)
			return err
		}
		return nil
	}
}
