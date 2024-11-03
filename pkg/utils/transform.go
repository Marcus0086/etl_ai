package utils

import (
	"encoding/json"
	"fmt"
	"formdata/pkg/models"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase/tools/types"
)

func ConfigToEnv(config models.Config) map[string]string {
	env := make(map[string]string)

	switch config := config.(type) {
	case *models.FileExtractorConfig:
		env["URL"] = config.URL
	case *models.JsonLoaderConfig:
		env["PATH"] = config.Path
	case *models.WebScraperConfig:
		env["URL"] = config.URL
	}

	return env
}

func ConfigFromEnv(configType string) (models.Config, error) {
	switch configType {
	case "file_extractor":
		url := os.Getenv("URL")
		if url == "" {
			return nil, fmt.Errorf("URL environment variable not set")
		}
		return &models.FileExtractorConfig{
			URL: url,
		}, nil
	case "json_loader":
		path := os.Getenv("PATH")
		if path == "" {
			return nil, fmt.Errorf("PATH environment variable not set")
		}
		return &models.JsonLoaderConfig{
			Path: path,
		}, nil
	default:
		return nil, fmt.Errorf("unknown config type: %s", configType)
	}
}

func ParseConfig(configType string, rawConfig types.JsonRaw) (models.Config, error) {
	configBytes, err := rawConfig.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	var config models.Config
	switch configType {
	case "file_extractor":
		c := &models.FileExtractorConfig{}
		if err := json.Unmarshal(configBytes, c); err != nil {
			return nil, fmt.Errorf("failed to unmarshal file extractor config: %w", err)
		}
		config = c
	case "json_loader":
		c := &models.JsonLoaderConfig{}
		if err := json.Unmarshal(configBytes, c); err != nil {
			return nil, fmt.Errorf("failed to unmarshal json loader config: %w", err)
		}
		config = c
	default:
		return nil, fmt.Errorf("unknown config type: %s", configType)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return config, nil
}

func BuildContainerEnv(config map[string]string) []string {
	env := []string{}
	for key, value := range config {
		env = append(env, fmt.Sprintf("%s=%s", strings.ToUpper(key), value))
	}
	return env
}
