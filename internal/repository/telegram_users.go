package repository

import (
	"context"
	"errors"

	"github.com/requiemofthesouls/pigeomail/internal/repository/entity"
	"github.com/requiemofthesouls/postgres"
)

const (
	sqlSelectUsers = "SELECT id, chat_id, email FROM pigeomail.telegram_users "
)

type tgUsers struct {
	db postgres.Wrapper
}

type TelegramUsers interface {
	GetByChatID(ctx context.Context, chatID int64) (*entity.TelegramUser, error)
	GetByEMail(ctx context.Context, email string) (*entity.TelegramUser, error)
	ExistsByEMail(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user *entity.TelegramUser) error
	DeleteByChatID(ctx context.Context, chatID int64) error
}

func NewUsers(db postgres.Wrapper) TelegramUsers {
	return &tgUsers{db: db}
}

func (u *tgUsers) GetByChatID(ctx context.Context, chatID int64) (*entity.TelegramUser, error) {
	var user entity.TelegramUser

	if err := u.db.QueryRow(
		ctx,
		sqlSelectUsers+"WHERE chat_id = $1",
		chatID,
	).Scan(
		&user.ID,
		&user.ChatID,
		&user.EMail,
	); err != nil {
		if errors.Is(err, postgres.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (u *tgUsers) GetByEMail(ctx context.Context, email string) (*entity.TelegramUser, error) {
	var user entity.TelegramUser

	if err := u.db.QueryRow(
		ctx,
		sqlSelectUsers+"WHERE tgUsers = $1",
		email,
	).Scan(
		&user.ID,
		&user.ChatID,
		&user.EMail,
	); err != nil {
		if errors.Is(err, postgres.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (u *tgUsers) ExistsByEMail(ctx context.Context, email string) (bool, error) {
	var exists bool

	if err := u.db.QueryRow(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM pigeomail.telegram_users WHERE email = $1)",
		email,
	).Scan(
		&exists,
	); err != nil {
		return false, err
	}

	return exists, nil
}

func (u *tgUsers) Create(ctx context.Context, user *entity.TelegramUser) error {
	err := u.db.QueryRow(
		ctx,
		"INSERT INTO pigeomail.telegram_users (chat_id, email) VALUES ($1, $2) RETURNING id",
		user.ChatID,
		user.EMail,
	).Scan(
		&user.ID,
	)

	return err
}

func (u *tgUsers) DeleteByChatID(ctx context.Context, chatID int64) error {
	_, err := u.db.Exec(
		ctx,
		"DELETE FROM pigeomail.telegram_users WHERE chat_id = $1",
		chatID,
	)

	return err
}
