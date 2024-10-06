package messagequeues

import "time"

type ETLMessage struct {
	Data            []byte    `json:"data"`
	MetaData        MetaData  `json:"metadata"`
	CreatedAt       time.Time `json:"created_at"`
	UpdateAt        time.Time `json:"updated_at"`
	StringifiedData string    `json:"stringified_data"`
	IsEnd           bool
}

type MetaData struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}
