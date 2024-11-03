package models

import "fmt"

type WebScraperConfig struct {
	URL string `json:"url"`
}

func (c *WebScraperConfig) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("URL is required")
	}
	return nil
}
