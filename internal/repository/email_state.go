package repository

import (
	"context"
	"fmt"

	"github.com/looplab/fsm"

	"github.com/requiemofthesouls/pigeomail/internal/customerrors"
	"github.com/requiemofthesouls/pigeomail/internal/repository/entity"
	"github.com/requiemofthesouls/pigeomail/pkg/state"
)

type EmailState interface {
	PrepareCreateEmail(ctx context.Context, chatID int64) error
	CreateEmail(ctx context.Context, email entity.EMail) error
	GetChatIDByEmail(ctx context.Context, email string) (chatID int64, err error)
	GetEmailByChatID(ctx context.Context, chatID int64) (email entity.EMail, err error)
	GetUserState(ctx context.Context, chatID int64) (userState string, ok bool)

	PrepareDeleteEmail(ctx context.Context, chatID int64) (email entity.EMail, err error)
	CancelDeleteEmail(ctx context.Context, chatID int64) error
	DeleteEmail(ctx context.Context, chatID int64) error
}

type emailState struct {
	repo  Email
	state *state.State
}

func NewEmailState(repo Email, state *state.State) EmailState {
	return &emailState{repo: repo, state: state}
}

func (s *emailState) PrepareCreateEmail(ctx context.Context, chatID int64) (err error) {
	var e entity.EMail
	e, err = s.repo.GetEmailByChatID(ctx, chatID)
	if err != nil && err != customerrors.ErrNotFound {
		return err
	}

	if err == nil {
		return customerrors.NewTelegramError("email already created: " + e.Name)
	}

	s.state.Add(chatID,
		fsm.NewFSM(
			entity.StateRequestedCreateEmail,
			fsm.Events{
				{
					Name: entity.StateCreateEmail,
					Src:  []string{entity.StateRequestedCreateEmail},
					Dst:  entity.StateEmailCreated,
				},
			},
			fsm.Callbacks{},
		))

	return nil
}

func (s *emailState) CreateEmail(ctx context.Context, email entity.EMail) (err error) {
	if err = email.Validate(); err != nil {
		return err
	}

	_, err = s.repo.GetEmailByName(ctx, email.Name)
	if err != nil && err != customerrors.ErrNotFound {
		return err
	}

	if err == nil {
		return customerrors.NewTelegramError("email <" + email.Name + "> already exists, please choose a new one")
	}

	var sm *fsm.FSM
	var ok bool
	if sm, ok = s.state.Get(email.ChatID); !ok {
		return customerrors.NewTelegramError("create email not requested, use /create")
	}

	if err = sm.Event(ctx, entity.StateCreateEmail); err != nil {
		return customerrors.NewTelegramError(fmt.Errorf("create email: %w", err).Error())
	}

	if err = s.repo.CreateEmail(ctx, email); err != nil {
		return err
	}

	s.state.Delete(email.ChatID)

	return nil
}

func (s *emailState) GetChatIDByEmail(ctx context.Context, email string) (chatID int64, err error) {
	return s.repo.GetChatIDByEmail(ctx, email)
}

func (s *emailState) GetEmailByChatID(ctx context.Context, chatID int64) (email entity.EMail, err error) {
	email, err = s.repo.GetEmailByChatID(ctx, chatID)
	if err != nil && err == customerrors.ErrNotFound {
		return email, customerrors.NewTelegramError("email not found, /create a new one")
	}

	return email, err
}

func (s *emailState) GetUserState(ctx context.Context, chatID int64) (userState string, ok bool) {
	var sm *fsm.FSM
	if sm, ok = s.state.Get(chatID); !ok {
		return "", false
	}

	return sm.Current(), true
}

func (s *emailState) PrepareDeleteEmail(ctx context.Context, chatID int64) (email entity.EMail, err error) {
	if email, err = s.repo.GetEmailByChatID(ctx, chatID); err != nil {
		if err == customerrors.ErrNotFound {
			return email, customerrors.NewTelegramError("there's no created email, use /create")
		}

		return email, err
	}

	s.state.Add(chatID,
		fsm.NewFSM(
			entity.StateRequestedDeleteEmail,
			fsm.Events{
				{
					Name: entity.StateDeleteEmail,
					Src:  []string{entity.StateRequestedDeleteEmail},
					Dst:  entity.StateEmailDeleted,
				},
				{
					Name: entity.StateCancelDeleteEmail,
					Src:  []string{entity.StateRequestedDeleteEmail},
					Dst:  entity.StateDeleteEmailCancelled,
				},
			},
			fsm.Callbacks{},
		))

	return email, nil
}

func (s *emailState) CancelDeleteEmail(ctx context.Context, chatID int64) (err error) {
	var sm *fsm.FSM
	var ok bool
	if sm, ok = s.state.Get(chatID); !ok {
		return customerrors.NewTelegramError("delete email not requested, use /delete")
	}

	if err = sm.Event(ctx, entity.StateCancelDeleteEmail); err != nil {
		return customerrors.NewTelegramError(fmt.Errorf("cancel delete email: %w", err).Error())
	}

	s.state.Delete(chatID)

	return nil
}

func (s *emailState) DeleteEmail(ctx context.Context, chatID int64) (err error) {
	var sm *fsm.FSM
	var ok bool
	if sm, ok = s.state.Get(chatID); !ok {
		return customerrors.NewTelegramError("delete email not requested, use /delete")
	}

	if err = sm.Event(ctx, entity.StateDeleteEmail); err != nil {
		return customerrors.NewTelegramError(fmt.Errorf("delete email: %w", err).Error())
	}

	if err = s.repo.DeleteEmail(ctx, chatID); err != nil {
		return err
	}

	s.state.Delete(chatID)

	return nil
}
