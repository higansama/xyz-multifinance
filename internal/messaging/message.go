package messaging

import (
	"time"
)

type Message struct {
	Id            string    `json:"id,omitempty"`
	AggregateType string    `json:"aggregatetype,omitempty"`
	AggregateId   string    `json:"aggregateid,omitempty"`
	Type          string    `json:"type,omitempty"`
	Payload       string    `json:"payload,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}
