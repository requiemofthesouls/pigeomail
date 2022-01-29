package repository

import (
	"context"

	"pigeomail/database"
	"pigeomail/pkg/client/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepo struct {
	client *mongo.Database
}

func (m *mongoRepo) GetChatIDByEmail(ctx context.Context, email string) (_ int64, err error) {
	var result EMail
	var collection = m.client.Collection("email")
	if err = collection.FindOne(ctx, bson.M{"name": email}).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, database.ErrNotFound
		}

		return 0, err
	}

	return result.ChatID, nil
}

func NewMongoRepository(host, port, username, password, database, authSource string) (_ IEmailRepository, err error) {
	var ctx = context.Background()

	var client *mongo.Database
	if client, err = mongodb.NewClient(
		ctx,
		host,
		port,
		username,
		password,
		database,
		authSource,
	); err != nil {
		return nil, err
	}

	return &mongoRepo{client: client}, nil
}

func (m *mongoRepo) GetEmailByChatID(ctx context.Context, chatID int64) (email EMail, err error) {
	collection := m.client.Collection("email")
	if err = collection.FindOne(ctx, bson.M{"chat_id": chatID}).Decode(&email); err != nil {
		if err == mongo.ErrNoDocuments {
			return email, database.ErrNotFound
		}

		return email, err
	}

	return email, nil
}

func (m *mongoRepo) GetEmailByName(ctx context.Context, name string) (email EMail, err error) {
	collection := m.client.Collection("email")
	if err = collection.FindOne(ctx, bson.M{"name": name}).Decode(&email); err != nil {
		if err == mongo.ErrNoDocuments {
			return email, database.ErrNotFound
		}

		return email, err
	}

	return email, nil
}

func (m *mongoRepo) CreateEmail(ctx context.Context, email EMail) (err error) {
	collection := m.client.Collection("email")
	if _, err = collection.InsertOne(ctx, email); err != nil {
		return err
	}

	return nil
}

func (m *mongoRepo) DeleteEmail(ctx context.Context, chatID int64) (err error) {
	collection := m.client.Collection("email")
	if _, err = collection.DeleteOne(ctx,
		bson.M{"chat_id": chatID},
	); err != nil {
		return err
	}

	return nil
}

func (m *mongoRepo) GetUserState(ctx context.Context, chatID int64) (state UserState, err error) {
	collection := m.client.Collection("state")
	if err = collection.FindOne(ctx, bson.M{"chat_id": chatID}).Decode(&state); err != nil {
		if err == mongo.ErrNoDocuments {
			return state, database.ErrNotFound
		}

		return state, err
	}

	return state, nil
}

func (m *mongoRepo) CreateUserState(ctx context.Context, state UserState) (err error) {
	var collection = m.client.Collection("state")

	var oldState UserState
	if oldState, err = m.GetUserState(ctx, state.ChatID); err != nil {
		if err != database.ErrNotFound {
			return err
		}

		if _, err = collection.InsertOne(ctx, state); err != nil {
			return err
		}

		return nil
	}

	if oldState.State == state.State {
		return nil
	}

	if _, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": oldState.ID},
		bson.M{"$set": bson.M{"state": state.State}},
	); err != nil {
		return err
	}

	return nil
}

func (m *mongoRepo) DeleteUserState(ctx context.Context, state UserState) (err error) {
	collection := m.client.Collection("state")
	if _, err = collection.DeleteMany(ctx,
		bson.M{
			"chat_id": state.ChatID,
			"state":   state.State,
		},
	); err != nil {
		return err
	}

	return nil
}
