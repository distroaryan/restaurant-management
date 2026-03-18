package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// TableStatus represents whether a table is available or fully booked.
type TableStatus string

const (
	TableStatusAvailable TableStatus = "AVAILABLE"
	TableStatusFull      TableStatus = "FULL"
)

type Table struct {
	ID            bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string        `bson:"name" json:"name" validate:"required"`
	UserID        string        `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Status        TableStatus   `bson:"status" json:"status" validate:"required"`
	CreatedAt     time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time     `bson:"updated_at" json:"updated_at"`
}
