package main

import (
	"encoding/json"
	"formdata/pkg/models"
	"formdata/pkg/utils"
	"log"
	"os"
)

func main() {
	log.Println("Loader started")

	etlProcess, err := utils.NewETLProcess("json_loader")
	if err != nil {
		log.Fatalf("Failed to create ETL process: %v", err)
	}
	defer etlProcess.Cleanup()

	channel, err := etlProcess.MqClient.NewChannel()
	if err != nil {
		log.Fatalf("Failed to create channel: %v", err)
	}
	outputFileName := etlProcess.Config.(*models.JsonLoaderConfig).Path + etlProcess.ConnectionID + ".json"
	outputFile, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open output file: %v", err)
	}
	defer outputFile.Close()

	msgs, err := etlProcess.MqClient.Consume(channel, etlProcess.QueueName)
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
