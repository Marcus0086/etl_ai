package models

import (
	"errors"
)

type FileExtractorConfig struct {
	URL string `json:"url,omitempty"`
}

func (c *FileExtractorConfig) Validate() error {
	if c.URL == "" {
		return errors.New("url is required")
	}
	// _, err := url.ParseRequestURI(c.URL)
	// if err != nil {
	// 	return errors.New("url is not valid")
	// }
	return nil
}
