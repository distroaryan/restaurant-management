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
	Capacity      int           `bson:"capacity" json:"capacity" validate:"required,min=1"`
	ReservedSeats int           `bson:"reserved_seats" json:"reserved_seats"`
	Status        TableStatus   `bson:"status" json:"status" validate:"required"`
	CreatedAt     time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time     `bson:"updated_at" json:"updated_at"`
}
