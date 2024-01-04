package models

import (
	"time"
)

const (
	StateActive  = "Active"
	StateDeleted = "Deleted"
)

type Entity struct {
	CreatedAt  time.Time `bson:"created_at"`
	ModifiedAt time.Time `bson:"modified_at"`
	State      string    `bson:"state"`
}
