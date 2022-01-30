package pigeomail

import (
	"context"

	"pigeomail/internal/errors"
)

type Service interface {
	PrepareCreateEmail(ctx context.Context, chatID int64) error
	CreateEmail(ctx context.Context, email EMail) error
	GetChatIDByEmail(ctx context.Context, email string) (chatID int64, err error)
	GetEmailByChatID(ctx context.Context, chatID int64) (email EMail, err error)
	GetUserState(ctx context.Context, chatID int64) (userState UserState, err error)

	PrepareDeleteEmail(ctx context.Context, chatID int64) (email EMail, err error)
	CancelDeleteEmail(ctx context.Context, chatID int64) error
	DeleteEmail(ctx context.Context, chatID int64) error
}

type service struct {
	storage Storage
}

func (s *service) GetEmailByChatID(ctx context.Context, chatID int64) (email EMail, err error) {
	email, err = s.storage.GetEmailByChatID(ctx, chatID)
	if err != nil && err == customerrors.ErrNotFound {
		return email, customerrors.NewTelegramError("email not found, /create a new one")
	}

	return email, err
}

func NewService(storage Storage) Service {
	return &service{storage: storage}
}

func (s *service) PrepareCreateEmail(ctx context.Context, chatID int64) (err error) {
	var email EMail
	email, err = s.storage.GetEmailByChatID(ctx, chatID)
	if err != nil && err != customerrors.ErrNotFound {
		return err
	}

	if err == nil {
		return customerrors.NewTelegramError("email already created: " + email.Name)
	}

	var userState = UserState{
		ChatID: chatID,
		State:  StateCreateEmailStep1,
	}

	if err = s.storage.CreateUserState(ctx, userState); err != nil {
		return err
	}

	return nil
}

func (s *service) PrepareDeleteEmail(ctx context.Context, chatID int64) (email EMail, err error) {
	if email, err = s.storage.GetEmailByChatID(ctx, chatID); err != nil {
		if err == customerrors.ErrNotFound {
			return email, customerrors.NewTelegramError("there's no created email, use /create")
		}

		return email, err
	}

	var userState = UserState{
		ChatID: chatID,
		State:  StateDeleteEmailStep1,
	}

	if err = s.storage.CreateUserState(ctx, userState); err != nil {
		return email, err
	}

	return email, nil
}

func (s *service) CancelDeleteEmail(ctx context.Context, chatID int64) (err error) {
	var userState = UserState{
		ChatID: chatID,
		State:  StateDeleteEmailStep1,
	}

	if err = s.storage.DeleteUserState(ctx, userState); err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteEmail(ctx context.Context, chatID int64) (err error) {
	if err = s.storage.DeleteEmail(ctx, chatID); err != nil {
		return err
	}

	var userState = UserState{
		ChatID: chatID,
		State:  StateDeleteEmailStep1,
	}

	if err = s.storage.DeleteUserState(ctx, userState); err != nil {
		return err
	}

	return nil
}

func (s *service) CreateEmail(ctx context.Context, email EMail) (err error) {
	if err = email.Validate(); err != nil {
		return err
	}

	_, err = s.storage.GetEmailByName(ctx, email.Name)
	if err != nil && err != customerrors.ErrNotFound {
		return err
	}

	if err == nil {
		return customerrors.NewTelegramError("email <" + email.Name + "> already exists, please choose a new one")
	}

	if err = s.storage.CreateEmail(ctx, email); err != nil {
		return err
	}

	var userState = UserState{
		ChatID: email.ChatID,
		State:  StateCreateEmailStep1,
	}

	if err = s.storage.DeleteUserState(ctx, userState); err != nil {
		return err
	}

	return nil
}

func (s *service) GetChatIDByEmail(ctx context.Context, email string) (chatID int64, err error) {
	return s.storage.GetChatIDByEmail(ctx, email)
}

func (s *service) GetUserState(ctx context.Context, chatID int64) (userState UserState, err error) {
	return s.storage.GetUserState(ctx, chatID)
}
