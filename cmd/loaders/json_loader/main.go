package main

import (
	"encoding/json"
	"etl/pkg/messagequeues"
	"log"
	"os"
)

func main() {
	log.Println("Loader started")

	mqClient, err := messagequeues.New()
	if err != nil {
		log.Fatalf("Failed to connect to message queue: %v", err)
	}
	defer mqClient.Close()

	channel, err := mqClient.NewChannel()
	if err != nil {
		log.Fatalf("Failed to create channel: %v", err)
	}

	outputFile, err := os.OpenFile("assets/data/output.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open output file: %v", err)
	}
	defer outputFile.Close()

	msgs, err := mqClient.Consume(channel, "extractor_queue")
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	for msg := range msgs {
		if msg.IsEnd {
			break
		}
		msg.StringifiedData = string(msg.Data)
		msg.Data = nil
		jsonEtlMessage, jerr := json.Marshal(msg)
		if jerr != nil {
			log.Printf("Failed to unmarshal message: %v", jerr)
			continue
		}
		_, err := outputFile.Write(jsonEtlMessage)
		if err != nil {
			log.Printf("Failed to write message to file: %v", err)
			continue
		}
	}

	log.Println("Loader finished")
}
