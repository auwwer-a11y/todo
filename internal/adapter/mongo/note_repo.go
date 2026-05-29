package mongo

import (
	"context"
	"github.com/auwwer-a11y/todo/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type NoteRepo struct {
	collection *mongo.Collection
}

func NewNoteRepo(db *mongo.Database) *NoteRepo {
	return &NoteRepo{
		collection: db.Collection("notes"),
	}
}

func (r *NoteRepo) Create(ctx context.Context, note *domain.Note) error {
	_, err := r.collection.InsertOne(ctx, note)
	return err
}

func (r *NoteRepo) GetAllByTaskID(ctx context.Context, taskID string) ([]*domain.Note, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"task_id": taskID})
	if err != nil {
		return nil, err
	}
	var notes []*domain.Note
	if err := cursor.All(ctx, &notes); err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *NoteRepo) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}