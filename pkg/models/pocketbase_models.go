package models

type RequestBody struct {
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

type ConnectionBody struct {
	SourceID string                 `json:"source_id"`
	LoaderID string                 `json:"loader_id"`
	SyncType string                 `json:"sync_type"`
	Schedule string                 `json:"schedule"`
	Config   map[string]interface{} `json:"config"`
}

type SyncRequest struct {
	ConnectionId string `json:"connection_id"`
}
