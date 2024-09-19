package outbox

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/higansama/xyz-multi-finance/internal/utils"
	"github.com/pkg/errors"
)

type Event struct {
	AggregateType string
	AggregateId   string
	Payload       any
}

type Table struct {
	Id            string    `bson:"_id,omitempty"`
	AggregateType string    `bson:"aggregatetype,omitempty"`
	AggregateId   string    `bson:"aggregateid,omitempty"`
	Type          string    `bson:"type,omitempty"`
	Payload       string    `bson:"payload,omitempty"`
	CreatedAt     time.Time `bson:"created_at"`
}

func NewOutbox(aggType string, aggId string, payload any) (*Table, error) {
	ps := utils.GetStructValue(payload)
	if ps.Kind() != reflect.Struct {
		return nil, errors.New("payload must be a struct")
	}

	bPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Table{
		AggregateType: aggType,
		AggregateId:   aggId,
		Type:          ps.Type().Name(),
		Payload:       string(bPayload),
		CreatedAt:     time.Now(),
	}, nil
}
