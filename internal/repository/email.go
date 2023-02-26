package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	mdb "go.mongodb.org/mongo-driver/mongo"

	"github.com/requiemofthesouls/pigeomail/internal/customerrors"
	"github.com/requiemofthesouls/pigeomail/internal/repository/entity"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/mongo"
)

type email struct {
	db mongo.Wrapper
}

type Email interface {
	GetEmailByChatID(ctx context.Context, chatID int64) (email entity.EMail, err error)
	GetChatIDByEmail(ctx context.Context, email string) (chatID int64, err error)
	GetEmailByName(ctx context.Context, name string) (email entity.EMail, err error)
	CreateEmail(ctx context.Context, email entity.EMail) (err error)
	DeleteEmail(ctx context.Context, chatID int64) (err error)
}

func NewEmail(db mongo.Wrapper) Email {
	return &email{db: db}
}

func (m *email) GetChatIDByEmail(ctx context.Context, email string) (_ int64, err error) {
	var result entity.EMail
	var collection = m.db.Collection("email")
	if err = collection.FindOne(ctx, bson.M{"name": email}).Decode(&result); err != nil {
		if err == mdb.ErrNoDocuments {
			return 0, customerrors.ErrNotFound
		}

		return 0, err
	}

	return result.ChatID, nil
}

func (m *email) GetEmailByChatID(ctx context.Context, chatID int64) (email entity.EMail, err error) {
	collection := m.db.Collection("email")
	if err = collection.FindOne(ctx, bson.M{"chat_id": chatID}).Decode(&email); err != nil {
		if err == mdb.ErrNoDocuments {
			return email, customerrors.ErrNotFound
		}

		return email, err
	}

	return email, nil
}

func (m *email) GetEmailByName(ctx context.Context, name string) (email entity.EMail, err error) {
	collection := m.db.Collection("email")
	if err = collection.FindOne(ctx, bson.M{"name": name}).Decode(&email); err != nil {
		if err == mdb.ErrNoDocuments {
			return email, customerrors.ErrNotFound
		}

		return email, err
	}

	return email, nil
}

func (m *email) CreateEmail(ctx context.Context, email entity.EMail) (err error) {
	collection := m.db.Collection("email")
	if _, err = collection.InsertOne(ctx, email); err != nil {
		return err
	}

	return nil
}

func (m *email) DeleteEmail(ctx context.Context, chatID int64) (err error) {
	collection := m.db.Collection("email")
	if _, err = collection.DeleteOne(ctx,
		bson.M{"chat_id": chatID},
	); err != nil {
		return err
	}

	return nil
}
