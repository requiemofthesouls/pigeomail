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
	State  state  `json:"state" bson:"state"`
}

type state string

const StateCreateEmailStep1 state = "create_email_step_1"
const StateDeleteEmailStep1 state = "delete_email_step_1"

type IEmailRepository interface {
	GetEmailByChatID(ctx context.Context, chatID int64) (email EMail, err error)
	GetEmailByName(ctx context.Context, name string) (email EMail, err error)
	CreateEmail(ctx context.Context, email EMail) (err error)
	DeleteEmail(ctx context.Context, chatID int64) (err error)

	GetUserState(ctx context.Context, chatID int64) (state UserState, err error)
	CreateUserState(ctx context.Context, state UserState) (err error)
	DeleteUserState(ctx context.Context, state UserState) (err error)
}
