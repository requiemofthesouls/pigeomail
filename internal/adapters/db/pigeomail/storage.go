package pigeomail

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"pigeomail/internal/domain/pigeomail"
	customerrors "pigeomail/pkg/errors"
)

type pigeomailStorage struct {
	db *mongo.Database
}

func NewStorage(db *mongo.Database) pigeomail.Storage {
	return &pigeomailStorage{db: db}
}

func (m *pigeomailStorage) GetChatIDByEmail(ctx context.Context, email string) (_ int64, err error) {
	var result pigeomail.EMail
	var collection = m.db.Collection("email")
	if err = collection.FindOne(ctx, bson.M{"name": email}).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, customerrors.ErrNotFound
		}

		return 0, err
	}

	return result.ChatID, nil
}

func (m *pigeomailStorage) GetEmailByChatID(ctx context.Context, chatID int64) (email pigeomail.EMail, err error) {
	collection := m.db.Collection("email")
	if err = collection.FindOne(ctx, bson.M{"chat_id": chatID}).Decode(&email); err != nil {
		if err == mongo.ErrNoDocuments {
			return email, customerrors.ErrNotFound
		}

		return email, err
	}

	return email, nil
}

func (m *pigeomailStorage) GetEmailByName(ctx context.Context, name string) (email pigeomail.EMail, err error) {
	collection := m.db.Collection("email")
	if err = collection.FindOne(ctx, bson.M{"name": name}).Decode(&email); err != nil {
		if err == mongo.ErrNoDocuments {
			return email, customerrors.ErrNotFound
		}

		return email, err
	}

	return email, nil
}

func (m *pigeomailStorage) CreateEmail(ctx context.Context, email pigeomail.EMail) (err error) {
	collection := m.db.Collection("email")
	if _, err = collection.InsertOne(ctx, email); err != nil {
		return err
	}

	return nil
}

func (m *pigeomailStorage) DeleteEmail(ctx context.Context, chatID int64) (err error) {
	collection := m.db.Collection("email")
	if _, err = collection.DeleteOne(ctx,
		bson.M{"chat_id": chatID},
	); err != nil {
		return err
	}

	return nil
}

func (m *pigeomailStorage) GetUserState(ctx context.Context, chatID int64) (state pigeomail.UserState, err error) {
	collection := m.db.Collection("state")
	if err = collection.FindOne(ctx, bson.M{"chat_id": chatID}).Decode(&state); err != nil {
		if err == mongo.ErrNoDocuments {
			return state, customerrors.ErrNotFound
		}

		return state, err
	}

	return state, nil
}

func (m *pigeomailStorage) CreateUserState(ctx context.Context, state pigeomail.UserState) (err error) {
	var collection = m.db.Collection("state")

	var oldState pigeomail.UserState
	if oldState, err = m.GetUserState(ctx, state.ChatID); err != nil {
		if err != customerrors.ErrNotFound {
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

func (m *pigeomailStorage) DeleteUserState(ctx context.Context, state pigeomail.UserState) (err error) {
	collection := m.db.Collection("state")
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
