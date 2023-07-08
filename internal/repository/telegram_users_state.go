package repository

import (
	"context"
	"fmt"

	"github.com/looplab/fsm"

	"github.com/requiemofthesouls/pigeomail/internal/customerrors"
	"github.com/requiemofthesouls/pigeomail/internal/repository/entity"
	"github.com/requiemofthesouls/pigeomail/pkg/state"
)

type TelegramUsersWithState interface {
	PrepareCreate(ctx context.Context, chatID int64) error
	Create(ctx context.Context, email *entity.TelegramUser) error

	GetByEMail(ctx context.Context, email string) (*entity.TelegramUser, error)
	GetByChatID(ctx context.Context, chatID int64) (*entity.TelegramUser, error)
	GetStateByChatID(chatID int64) (userState string, ok bool)

	PrepareDelete(ctx context.Context, chatID int64) (*entity.TelegramUser, error)
	CancelDelete(ctx context.Context, chatID int64) error
	Delete(ctx context.Context, chatID int64) error
}

type tgUsersWithState struct {
	repo  TelegramUsers
	state *state.State
}

func NewUsersWithState(repo TelegramUsers, state *state.State) TelegramUsersWithState {
	return &tgUsersWithState{repo: repo, state: state}
}

func (s *tgUsersWithState) PrepareCreate(ctx context.Context, chatID int64) error {
	var (
		user *entity.TelegramUser
		err  error
	)
	if user, err = s.repo.GetByChatID(ctx, chatID); err != nil {
		return err
	}

	if !user.IsExist() {
		return customerrors.NewTelegramError("tgUsers already created: " + user.EMail)
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

func (s *tgUsersWithState) Create(ctx context.Context, user *entity.TelegramUser) error {
	if err := user.ValidateEMail(); err != nil {
		return err
	}

	var isExist bool
	var err error
	if isExist, err = s.repo.ExistsByEMail(ctx, user.EMail); err != nil {
		return err
	}

	if isExist {
		return customerrors.NewTelegramError("tgUsers <" + user.EMail + "> already exists, please choose a new one")
	}

	var sm *fsm.FSM
	var ok bool
	if sm, ok = s.state.Get(user.ChatID); !ok {
		return customerrors.NewTelegramError("create tgUsers not requested, use /create")
	}

	if err = sm.Event(ctx, entity.StateCreateEmail); err != nil {
		return customerrors.NewTelegramError(fmt.Errorf("create tgUsers: %w", err).Error())
	}

	if err = s.repo.Create(ctx, user); err != nil {
		return err
	}

	s.state.Delete(user.ChatID)

	return nil
}

func (s *tgUsersWithState) GetByEMail(ctx context.Context, email string) (*entity.TelegramUser, error) {
	return s.repo.GetByEMail(ctx, email)
}

func (s *tgUsersWithState) GetByChatID(ctx context.Context, chatID int64) (*entity.TelegramUser, error) {
	var (
		user *entity.TelegramUser
		err  error
	)
	if user, err = s.repo.GetByChatID(ctx, chatID); err != nil {
		return nil, err
	}

	if !user.IsExist() {
		return nil, customerrors.NewTelegramError("tgUsers not found, /create a new one")
	}

	return user, err
}

func (s *tgUsersWithState) GetStateByChatID(chatID int64) (userState string, ok bool) {
	var sm *fsm.FSM
	if sm, ok = s.state.Get(chatID); !ok {
		return "", false
	}

	return sm.Current(), true
}

func (s *tgUsersWithState) PrepareDelete(ctx context.Context, chatID int64) (*entity.TelegramUser, error) {
	var (
		user *entity.TelegramUser
		err  error
	)

	if user, err = s.repo.GetByChatID(ctx, chatID); err != nil {
		return nil, err
	}

	if !user.IsExist() {
		return nil, customerrors.NewTelegramError("there's no created tgUsers, use /create")
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

	return user, nil
}

func (s *tgUsersWithState) CancelDelete(ctx context.Context, chatID int64) error {
	var (
		sm *fsm.FSM
		ok bool
	)
	if sm, ok = s.state.Get(chatID); !ok {
		return customerrors.NewTelegramError("delete tgUsers not requested, use /delete")
	}

	if err := sm.Event(ctx, entity.StateCancelDeleteEmail); err != nil {
		return customerrors.NewTelegramError(fmt.Errorf("cancel delete tgUsers: %w", err).Error())
	}

	s.state.Delete(chatID)

	return nil
}

func (s *tgUsersWithState) Delete(ctx context.Context, chatID int64) error {
	var sm *fsm.FSM
	var ok bool
	if sm, ok = s.state.Get(chatID); !ok {
		return customerrors.NewTelegramError("delete tgUsers not requested, use /delete")
	}

	if err := sm.Event(ctx, entity.StateDeleteEmail); err != nil {
		return customerrors.NewTelegramError(fmt.Errorf("delete tgUsers: %w", err).Error())
	}

	if err := s.repo.DeleteByChatID(ctx, chatID); err != nil {
		return err
	}

	s.state.Delete(chatID)

	return nil
}
