package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"

	"formdata/pkg/messagequeues"
	"formdata/pkg/models"
	"formdata/pkg/utils"
	"formdata/pkg/wokrerpool"
)

const (
	chunkSize  = 1024 * 1024 // 1MB
	numWorkers = 4           // Number of worker goroutines
)

func main() {
	log.Println("Extractor started")

	etlProcess, err := utils.NewETLProcess("file_extractor")
	if err != nil {
		log.Fatalf("Failed to create ETL process: %v", err)
	}
	defer etlProcess.Cleanup()
	extractFile(etlProcess.Config.(*models.FileExtractorConfig).URL, etlProcess.MqClient, etlProcess.QueueName)

	ch, err := etlProcess.MqClient.NewChannel()
	if err != nil {
		log.Printf("Failed to create channel: %v", err)
		return
	}
	defer ch.Close()

	if err := etlProcess.SendEndMessage(); err != nil {
		log.Printf("Failed to send end message: %v", err)
	}

	log.Println("Extraction complete.")
}

func extractFile(filePath string, mqClient *messagequeues.RabbitMQClient, queueName string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	wp := wokrerpool.New[[]byte](numWorkers)
	workerFunc := utils.CreateWorkerFunc(mqClient, queueName, "assets/data/big.txt", "extractor")
	wp.Start(ctx, workerFunc)

	reader := bufio.NewReader(file)
	buffer := make([]byte, chunkSize)

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			log.Printf("Error reading file: %v", err)
			continue
		}

		wp.Submit(buffer[:n])
	}
	wp.Stop()
}
