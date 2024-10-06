package main

import (
	"etl/pkg/dockermanager"
	"log"
	"sync"
	"time"
)

func main() {
	containerConfigs := []dockermanager.ContainerConfig{
		{
			Image:       "etl-extractor:latest",
			Name:        "etl_extractor",
			Env:         []string{},
			Cmd:         []string{},
			Network:     "etl_network",
			MemoryLimit: 512 * 1024 * 1024,
			CPUShares:   2,
			Mounts:      []string{"/root/assets"},
		},
		{
			Image:       "etl-loader:latest",
			Name:        "etl_loader",
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
