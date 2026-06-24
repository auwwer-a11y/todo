package mongo

import (
	"context"
	"github.com/auwwer-a11y/todo/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationRepo struct {
	collection *mongo.Collection
}

func NewNotificationRepo(db *mongo.Database) *NotificationRepo {
	return &NotificationRepo{
		collection: db.Collection("notifications"),
	}
}

func (r *NotificationRepo) Create(ctx context.Context, notification *domain.Notification) error {
	_, err := r.collection.InsertOne(ctx, notification)
	return err
}

