package repository

import (
	"context"

	"github.com/distroaryan/restaurant-management/internal/database"
	"github.com/distroaryan/restaurant-management/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OrderRepository struct {
	collection *mongo.Collection
}

func NewOrderRepository(db *database.Database) *OrderRepository {
	return &OrderRepository{
		collection: db.GetCollection("orders"),
	}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	order.Status = models.OrderStatusPending
	resp, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return err
	}
	if oid, ok := resp.InsertedID.(bson.ObjectID); ok {
		order.ID = oid
	}
	return nil
}

func (r *OrderRepository) GetOrderById(ctx context.Context, id string) (*models.Order, error) {
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectId}
	var order models.Order

	err = r.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetOrdersByTable(ctx context.Context, tableId string) ([]*models.Order, error) {
	objectId, err := bson.ObjectIDFromHex(tableId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"table_id": objectId}
	var orders []*models.Order

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, orderId string, orderStatus models.OrderStatus) error {
	objectId, err := bson.ObjectIDFromHex(orderId)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objectId}
	update := bson.M{
		"$set": bson.M{"status": orderStatus},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}
