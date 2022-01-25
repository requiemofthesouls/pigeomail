package repository

import (
	"context"
)

type EMail struct {
	ID     string `bson:"_id,omitempty"`
	ChatID int64  `json:"chat_id" bson:"chat_id"`
	Name   string `json:"name" bson:"name"`
}

type IEmailRepository interface {
	GetEmailByChatID(ctx context.Context, chatID int64) (email EMail, err error)
	GetChatIDByEmail(ctx context.Context, email string) (chatID int64, err error)
	GetEmailByName(ctx context.Context, name string) (email EMail, err error)
	CreateEmail(ctx context.Context, email EMail) (err error)
	DeleteEmail(ctx context.Context, chatID int64) (err error)
}
