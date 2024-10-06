package orchestrator

import (
	"etl/pkg/dockermanager"
	"log"
	"sync"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/cron"
)

var (
	sources = map[string]string{
		"file_extractor": "etl-extractor:latest",
	}
	loaders = map[string]string{
		"json_loader": "etl-loader:latest",
	}
)

func ConfigureOrchestrator(app *core.App, connection *models.Record) error {
	syncType := connection.GetString("sync_type")
	switch syncType {
	case "manual":
		// Do something
		// return nil
	case "scheduled":
		schedule := connection.GetString("schedule")
		go scheduler(app, connection, schedule)
	default:
		log.Fatal("Invalid sync type")
	}
	return nil
}

func scheduler(app *core.App, connection *models.Record, schedule string) {
	c := cron.New()
	if err := c.Add(connection.Id, schedule, func() {
		log.Printf("Running job for connection %s", connection.Id)
		StartEtlWorkflow(app, connection)
	}); err != nil {
		log.Fatal(err)
	}
	c.Start()
}

func StartEtlWorkflow(app *core.App, connection *models.Record) {
	sourceId := connection.GetStringSlice("source_id")[0]
	loaderId := connection.GetStringSlice("loader_id")[0]
	log.Println("Starting ETL Workflow for connection", connection.Id, "Source ID:", sourceId, "Loader ID:", loaderId)
	pbApp := *app
	sourceRecord, err := pbApp.Dao().FindRecordById("sources", sourceId)
	if err != nil {
		log.Fatal(err)
	}
	loaderRecord, err := pbApp.Dao().FindRecordById("loaders", loaderId)
	if err != nil {
		log.Fatal(err)
	}

	sourceImage := sources[sourceRecord.GetString("type")]
	loaderImage := loaders[loaderRecord.GetString("type")]

	log.Println("Source and Loader Image loaded:", sourceImage, loaderImage)

	containerConfigs := []dockermanager.ContainerConfig{
		{
			Image: sourceImage,
			Name:  sourceId,
			Env: []string{
				"CONNECTION_ID=" + connection.Id,
				"SOURCE_ID=" + sourceId,
			},
			Cmd:         []string{},
			Network:     "etl_network",
			MemoryLimit: 512 * 1024 * 1024,
			CPUShares:   2,
			Mounts:      []string{"/root/assets"},
		},
		{
			Image: loaderImage,
			Name:  loaderId,
			Env: []string{
				"CONNECTION_ID=" + connection.Id,
				"LOADER_ID=" + loaderId,
			},
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
