package orchestrator

import (
	"etl/pkg/dockermanager"
	"log"
	"sync"
	"time"

	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/cron"
)

func ConfigureOrchestrator(connection *models.Record) error {
	syncType := connection.GetString("sync_type")
	switch syncType {
	case "manual":
		// Do something
		// return nil
	case "scheduled":
		schedule := connection.GetString("schedule")
		go scheduler(connection, schedule)
	default:
		log.Fatal("Invalid sync type")
	}
	return nil
}

func scheduler(connection *models.Record, schedule string) {
	c := cron.New()
	if err := c.Add(connection.Id, schedule, func() {
		log.Printf("Running job for connection %s", connection.Id)
		StartEtlWorkflow(connection)
	}); err != nil {
		log.Fatal(err)
	}
	c.Start()
}

func StartEtlWorkflow(connection *models.Record) {
	sourceId := connection.GetString("source_id")
	loaderId := connection.GetString("loader_id")
	containerConfigs := []dockermanager.ContainerConfig{
		{
			Image:       "etl-extractor:latest",
			Name:        sourceId,
			Env:         []string{},
			Cmd:         []string{},
			Network:     "etl_network",
			MemoryLimit: 512 * 1024 * 1024,
			CPUShares:   2,
			Mounts:      []string{"/root/assets"},
		},
		{
			Image:       "etl-loader:latest",
			Name:        loaderId,
			Env:         []string{},
			Cmd:         []string{},
			Network:     "etl_network",
			MemoryLimit: 256 * 1024 * 1024,
			CPUShares:   1,
			Mounts:      []string{"/root/assets"},
		},
	}

	var wg sync.WaitGroup
	containerIDs := make([]string, 0, len(containerConfigs))
	for _, config := range containerConfigs {
		wg.Add(1)
		go func(cfg dockermanager.ContainerConfig) {
			defer wg.Done()
			containerID, err := dockermanager.StartContainer(&cfg)
			if err != nil {
				panic(err)
			}
			containerIDs = append(containerIDs, containerID)
		}(config)
	}
	wg.Wait()

	time.Sleep(2 * time.Second)

	for _, containerID := range containerIDs {
		err := dockermanager.StopContainer(containerID)
		if err != nil {
			log.Printf("Failed to stop container %s: %v", containerID, err)
		} else {
			log.Printf("Stopped container %s", containerID)
		}
	}
}
