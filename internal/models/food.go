package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Food struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string        `bson:"name" json:"name" validate:"required"`
	Price     float64       `bson:"price" json:"price" validate:"required"`
	MenuID    bson.ObjectID `bson:"menu_id" json:"menu_id" validate:"required"`
	Image     string        `bson:"image" json:"image"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}
