package main

import (
	"formdata/pkg/messagequeues"
	"formdata/pkg/utils"
	"log"
)

func main() {
	log.Println("Web scraper started")

	etlProcess, err := utils.NewETLProcess("web_scraper")
	if err != nil {
		log.Fatalf("Failed to create ETL process: %v", err)
	}
	defer etlProcess.Cleanup()
	// TODO: Implement web scraping

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

func scrapeWebsite(url string, mqClient *messagequeues.RabbitMQClient, queueName string) {

}
