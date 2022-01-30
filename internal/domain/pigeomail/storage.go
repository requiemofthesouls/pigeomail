package pigeomail

import (
	"context"
)

type Storage interface {
	GetEmailByChatID(ctx context.Context, chatID int64) (email EMail, err error)
	GetChatIDByEmail(ctx context.Context, email string) (chatID int64, err error)
	GetEmailByName(ctx context.Context, name string) (email EMail, err error)
	CreateEmail(ctx context.Context, email EMail) (err error)
	DeleteEmail(ctx context.Context, chatID int64) (err error)

	GetUserState(ctx context.Context, chatID int64) (state UserState, err error)
	CreateUserState(ctx context.Context, state UserState) (err error)
	DeleteUserState(ctx context.Context, state UserState) (err error)
}
