package main

import (
	"bufio"
	"context"
	"etl/pkg/messagequeues"
	"etl/pkg/wokrerpool"
	"io"
	"log"
	"os"
	"time"
)

const (
	chunkSize  = 1024 * 1024 // 1MB
	numWorkers = 4           // Number of worker goroutines
	queueName  = "extractor_queue"
)

func main() {
	log.Println("Extractor started")

	mqClient, err := messagequeues.New()
	if err != nil {
		log.Fatalf("Failed to connect to message queue: %v", err)
	}
	defer mqClient.Close()

	extractFile("assets/data/big.txt", mqClient)

	ch, err := mqClient.NewChannel()
	if err != nil {
		log.Printf("Failed to create channel: %v", err)
		return
	}
	defer ch.Close()

	endMessage := messagequeues.ETLMessage{
		IsEnd: true,
	}

	err = mqClient.Publish(ch, queueName, endMessage)
	if err != nil {
		log.Printf("Failed to publish end message: %v", err)
	}

	log.Println("Extraction complete.")
}

func extractFile(filePath string, mqClient *messagequeues.RabbitMQClient) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	wp := wokrerpool.New[[]byte](numWorkers)
	workerFunc := createWorkerFunc(mqClient)
	wp.Start(ctx, workerFunc)

	reader := bufio.NewReader(file)
	buffer := make([]byte, chunkSize)

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error reading file: %v", err)
			continue
		}
		chunk := make([]byte, n)
		copy(chunk, buffer[:n])

		wp.Submit(chunk)
	}
	wp.Stop()
}

func createWorkerFunc(mqClient *messagequeues.RabbitMQClient) wokrerpool.WorkerFunc[[]byte] {
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
				Source:      "assets/data/big.txt",
				Destination: "extractor",
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
