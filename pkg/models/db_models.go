package models

type Config interface {
	Validate() error
}
type RequestBody struct {
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config,omitempty"`
}

type ConnectionBody struct {
	SourceID string                 `json:"source_id"`
	LoaderID string                 `json:"loader_id"`
	SyncType string                 `json:"sync_type"`
	Schedule string                 `json:"schedule,omitempty"`
	Config   map[string]interface{} `json:"config,omitempty"`
}

type SyncRequest struct {
	ConnectionId string `json:"connection_id"`
}
