package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"pigeomail/database"
)

type mongoRepo struct {
	client *mongo.Client
}

func NewMongoRepository(cfg *database.Config) (_ IEmailRepository, err error) {
	var client *mongo.Client
	if client, err = database.NewMongoDBClient(cfg); err != nil {
		return nil, err
	}

	return &mongoRepo{client: client}, nil
}

func (m *mongoRepo) GetEmailByChatID(ctx context.Context, chatID int64) (email EMail, err error) {
	collection := m.client.Database("pigeomail").Collection("email")
	if err = collection.FindOne(ctx, bson.D{{"chat_id", chatID}}).Decode(&email); err != nil {
		if err == mongo.ErrNoDocuments {
			return email, database.ErrNotFound
		}

		return email, err
	}

	return email, nil
}

func (m *mongoRepo) GetEmailByName(ctx context.Context, name string) (email EMail, err error) {
	collection := m.client.Database("pigeomail").Collection("email")
	if err = collection.FindOne(ctx, bson.D{{"name", name}}).Decode(&email); err != nil {
		if err == mongo.ErrNoDocuments {
			return email, database.ErrNotFound
		}

		return email, err
	}

	return email, nil
}

func (m *mongoRepo) CreateEmail(ctx context.Context, email EMail) (err error) {
	collection := m.client.Database("pigeomail").Collection("email")
	if _, err = collection.InsertOne(ctx, email); err != nil {
		return err
	}

	return nil
}

func (m *mongoRepo) DeleteEmail(ctx context.Context, email EMail) (err error) {
	collection := m.client.Database("pigeomail").Collection("email")
	if _, err = collection.DeleteOne(ctx,
		bson.D{
			{"name", email.Name},
			{"chat_id", email.ChatID},
		},
	); err != nil {
		return err
	}

	return nil
}
