package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "PENDING"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusCompleted  OrderStatus = "COMPLETED"
	OrderStatusCancelled  OrderStatus = "CANCELLED"
)

type OrderItem struct {
	FoodID    bson.ObjectID `bson:"food_id" json:"food_id" validate:"required"`
	Quantity  int           `bson:"quantity" json:"quantity" validate:"required,min=1"`
	UnitPrice float64       `bson:"unit_price" json:"unit_price" validate:"required,min=0"`
}

type Order struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	TableID   bson.ObjectID `bson:"table_id,omitempty" json:"table_id,omitempty"`
	UserID    string        `bson:"user_id,omitempty" json:"user_id,omitempty"` // Example: From JWT, could be empty if not required
	Status    OrderStatus   `bson:"status" json:"status"`
	Items       []OrderItem   `bson:"items" json:"items" validate:"required,min=1"`
	TotalAmount float64       `bson:"total_amount" json:"total_amount"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}
