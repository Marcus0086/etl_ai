package models

import "errors"

type JsonLoaderConfig struct {
	Path string `json:"path,omitempty"`
}

func (c *JsonLoaderConfig) Validate() error {
	if c.Path == "" {
		return errors.New("path is required")
	}
	return nil
}
