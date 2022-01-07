package repository

import (
	"context"
)

type EMail struct {
	ID     string `bson:"_id,omitempty"`
	ChatID int64  `json:"chat_id" bson:"chat_id"`
	Name   string `json:"name" bson:"name"`
}

type UserState struct {
	ID     string `bson:"_id,omitempty"`
	ChatID int64  `json:"chat_id" bson:"chat_id"`
	State  string `json:"state" bson:"state"`
}

const StateEmailCreationStep2 = "email_creation_step_2"

type IEmailRepository interface {
	GetEmailByChatID(ctx context.Context, chatID int64) (email EMail, err error)
	GetEmailByName(ctx context.Context, name string) (email EMail, err error)
	CreateEmail(ctx context.Context, email EMail) (err error)
	DeleteEmail(ctx context.Context, email EMail) (err error)

	GetUserStateByChatID(ctx context.Context, chatID int64) (state UserState, err error)
	CreateUserState(ctx context.Context, state UserState) (err error)
	DeleteUserState(ctx context.Context, state UserState) (err error)
}
