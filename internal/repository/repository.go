package repository

import (
	"context"
)

type EMail struct {
	ChatID int64  `json:"chat_id"`
	Name   string `json:"name"`
}

type IEmailRepository interface {
	GetEmailByChatID(ctx context.Context, chatID int64) (email EMail, err error)
	GetEmailByName(ctx context.Context, name string) (email EMail, err error)
	CreateEmail(ctx context.Context, email EMail) (err error)
	DeleteEmail(ctx context.Context, email EMail) (err error)
}
