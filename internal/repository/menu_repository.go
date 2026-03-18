package repository

import (
	"context"

	"github.com/distroaryan/restaurant-management/internal/database"
	"github.com/distroaryan/restaurant-management/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MenuRepository struct {
	collection *mongo.Collection
}

func NewMenuRepository(db *database.Database) *MenuRepository {
	return &MenuRepository{
		collection: db.GetCollection("menus"),
	}
}

func (r *MenuRepository) CreateMenu(ctx context.Context, menu *models.Menu) error {
	resp, err := r.collection.InsertOne(ctx, menu)
	if err != nil {
		return err
	}
	if objectId, ok := resp.InsertedID.(bson.ObjectID); ok {
		menu.ID = objectId
	}
	return nil
}

func (r *MenuRepository) GetMenuById(ctx context.Context, id string) (*models.Menu, error) {
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectId}
	var menu models.Menu

	err = r.collection.FindOne(ctx, filter).Decode(&menu)
	if err != nil {
		return nil, err
	}

	return &menu, nil
}

func (r *MenuRepository) GetAllMenu(ctx context.Context) ([]*models.Menu, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var menus []*models.Menu
	if err := cursor.All(ctx, &menus); err != nil {
		return nil, err
	}
	return menus, nil
}
