package inbox

import (
	"time"
)

type Inbox struct {
	Id        string    `bson:"_id,omitempty"`
	MessageId string    `bson:"message_id,omitempty"`
	Consumer  string    `bson:"consumer,omitempty"`
	CreatedAt time.Time `bson:"created_at"`
}

func New(msgId string, consumer string) *Inbox {
	return &Inbox{
		MessageId: msgId,
		Consumer:  consumer,
		CreatedAt: time.Now(),
	}
}
