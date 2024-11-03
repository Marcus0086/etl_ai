package orchestrator

import (
	"log"
	"maps"
	"sync"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/cron"
	"github.com/pocketbase/pocketbase/tools/types"

	"formdata/pkg/dockermanager"
	"formdata/pkg/utils"
)

var (
	sources = map[string]string{
		"file_extractor": "formdata-extractor:latest",
	}
	loaders = map[string]string{
		"json_loader": "formdata-loader:latest",
	}
)

func ConfigureOrchestrator(app *core.App, connection *models.Record) error {
	syncType := connection.GetString("sync_type")
	switch syncType {
	case "manual":
		return nil
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
	log.Println(
		"Starting ETL Workflow for connection",
		connection.Id,
		"Source ID:",
		sourceId,
		"Loader ID:",
		loaderId,
	)
	pbApp := *app
	sourceRecord, err := pbApp.Dao().FindRecordById("sources", sourceId)
	if err != nil {
		log.Fatal(err)
	}
	loaderRecord, err := pbApp.Dao().FindRecordById("loaders", loaderId)
	if err != nil {
		log.Fatal(err)
	}

	sourceType := sourceRecord.GetString("type")
	loaderType := loaderRecord.GetString("type")

	sourceConfigJsonRaw := sourceRecord.Get("config").(types.JsonRaw)
	sourceConfig, err := utils.ParseConfig(sourceType, sourceConfigJsonRaw)
	if err != nil {
		log.Fatal(err)
	}

	loaderConfigJsonRaw := loaderRecord.Get("config").(types.JsonRaw)
	loaderConfig, err := utils.ParseConfig(loaderType, loaderConfigJsonRaw)
	if err != nil {
		log.Fatal(err)
	}

	sourceImage := sources[sourceType]
	loaderImage := loaders[loaderType]

	log.Println("Source and Loader Image loaded:", sourceImage, loaderImage)
	sourceEnv := map[string]string{
		"CONNECTION_ID": connection.Id,
		"SOURCE_ID":     sourceId,
		"QUEUE_NAME":    sourceType + "_" + connection.Id,
	}
	maps.Copy(sourceEnv, utils.ConfigToEnv(sourceConfig))

	loaderEnv := map[string]string{
		"CONNECTION_ID": connection.Id,
		"LOADER_ID":     loaderId,
		"QUEUE_NAME":    sourceType + "_" + connection.Id,
	}
	maps.Copy(loaderEnv, utils.ConfigToEnv(loaderConfig))

	containerConfigs := []dockermanager.ContainerConfig{
		{
			Image:       sourceImage,
			Name:        sourceId,
			Env:         utils.BuildContainerEnv(sourceEnv),
			Cmd:         []string{},
			Network:     "formdata_network",
			MemoryLimit: 1024 * 1024 * 1024,
			CPUShares:   4,
			Mounts:      []string{"/root/assets"},
			AutoRemove:  true,
		},
		{
			Image:       loaderImage,
			Name:        loaderId,
			Env:         utils.BuildContainerEnv(loaderEnv),
			Cmd:         []string{},
			Network:     "formdata_network",
			MemoryLimit: 256 * 1024 * 1024,
			CPUShares:   1,
			Mounts:      []string{"/root/assets"},
			AutoRemove:  true,
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
}
