package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Menu struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string        `bson:"name" json:"name" validate:"required"`
	Description string        `bson:"description" json:"description"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at" json:"updated_at"`
}
