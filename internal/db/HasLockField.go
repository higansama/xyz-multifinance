package db

import "github.com/google/uuid"

type HasLockField struct {
	LockId string `json:"lock_id,omitempty" bson:"lock_id,omitempty"`
}

func (l *HasLockField) UpdateLockField() {
	l.LockId = uuid.New().String()
}
