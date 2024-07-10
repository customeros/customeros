package events

import (
	"time"
)

type BaseEvent struct {
	CreatedAt time.Time `json:"createdAt"`
	AppSource string    `json:"appSource"`
	Source    string    `json:"source"`
}
