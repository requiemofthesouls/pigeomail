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
