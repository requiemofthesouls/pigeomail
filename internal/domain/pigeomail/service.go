package pigeomail

import (
	"context"
	"fmt"
	"sync"

	customerrors "pigeomail/internal/errors"

	"github.com/looplab/fsm"
)

type Service interface {
	PrepareCreateEmail(ctx context.Context, chatID int64) error
	CreateEmail(ctx context.Context, email EMail) error
	GetChatIDByEmail(ctx context.Context, email string) (chatID int64, err error)
	GetEmailByChatID(ctx context.Context, chatID int64) (email EMail, err error)
	GetUserState(ctx context.Context, chatID int64) (userState string, ok bool)

	PrepareDeleteEmail(ctx context.Context, chatID int64) (email EMail, err error)
	CancelDeleteEmail(ctx context.Context, chatID int64) error
	DeleteEmail(ctx context.Context, chatID int64) error
}

// Создаем структуру, которая объединяет карту и мьютекс
type FSMMap struct {
	mu sync.Mutex
	m  map[int64]*fsm.FSM
}

func NewFSMMap() *FSMMap {
	return &FSMMap{
		m: make(map[int64]*fsm.FSM),
	}
}

func (fsmMap *FSMMap) Add(key int64, fsm *fsm.FSM) {
	fsmMap.mu.Lock()
	defer fsmMap.mu.Unlock()
	fsmMap.m[key] = fsm
}

func (fsmMap *FSMMap) Get(key int64) (*fsm.FSM, bool) {
	fsmMap.mu.Lock()
	defer fsmMap.mu.Unlock()
	val, ok := fsmMap.m[key]
	return val, ok
}

func (fsmMap *FSMMap) Delete(key int64) {
	fsmMap.mu.Lock()
	defer fsmMap.mu.Unlock()
	delete(fsmMap.m, key)
}

type service struct {
	storage Storage
	state   *FSMMap
}

func NewService(storage Storage) Service {
	return &service{storage: storage, state: NewFSMMap()}
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

	s.state.Add(chatID,
		fsm.NewFSM(
			StateRequestedCreateEmail,
			fsm.Events{
				{
					Name: StateCreateEmail,
					Src:  []string{StateRequestedCreateEmail},
					Dst:  StateEmailCreated,
				},
			},
			fsm.Callbacks{},
		))

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

	var state *fsm.FSM
	var ok bool
	if state, ok = s.state.Get(email.ChatID); !ok {
		return customerrors.NewTelegramError("create email not requested, use /create")
	}

	if err = state.Event(ctx, StateCreateEmail); err != nil {
		return customerrors.NewTelegramError(fmt.Errorf("create email: %w", err).Error())
	}

	if err = s.storage.CreateEmail(ctx, email); err != nil {
		return err
	}

	s.state.Delete(email.ChatID)

	return nil
}

func (s *service) GetChatIDByEmail(ctx context.Context, email string) (chatID int64, err error) {
	return s.storage.GetChatIDByEmail(ctx, email)
}

func (s *service) GetEmailByChatID(ctx context.Context, chatID int64) (email EMail, err error) {
	email, err = s.storage.GetEmailByChatID(ctx, chatID)
	if err != nil && err == customerrors.ErrNotFound {
		return email, customerrors.NewTelegramError("email not found, /create a new one")
	}

	return email, err
}

func (s *service) GetUserState(ctx context.Context, chatID int64) (userState string, ok bool) {
	var state *fsm.FSM
	if state, ok = s.state.Get(chatID); !ok {
		return "", false
	}

	return state.Current(), true
}

func (s *service) PrepareDeleteEmail(ctx context.Context, chatID int64) (email EMail, err error) {
	if email, err = s.storage.GetEmailByChatID(ctx, chatID); err != nil {
		if err == customerrors.ErrNotFound {
			return email, customerrors.NewTelegramError("there's no created email, use /create")
		}

		return email, err
	}

	s.state.Add(chatID,
		fsm.NewFSM(
			StateRequestedDeleteEmail,
			fsm.Events{
				{
					Name: StateDeleteEmail,
					Src:  []string{StateRequestedDeleteEmail},
					Dst:  StateEmailDeleted,
				},
				{
					Name: StateCancelDeleteEmail,
					Src:  []string{StateRequestedDeleteEmail},
					Dst:  StateDeleteEmailCancelled,
				},
			},
			fsm.Callbacks{},
		))

	return email, nil
}

func (s *service) CancelDeleteEmail(ctx context.Context, chatID int64) (err error) {
	var state *fsm.FSM
	var ok bool
	if state, ok = s.state.Get(chatID); !ok {
		return customerrors.NewTelegramError("delete email not requested, use /delete")
	}

	if err = state.Event(ctx, StateCancelDeleteEmail); err != nil {
		return customerrors.NewTelegramError(fmt.Errorf("cancel delete email: %w", err).Error())
	}

	s.state.Delete(chatID)

	return nil
}

func (s *service) DeleteEmail(ctx context.Context, chatID int64) (err error) {
	var state *fsm.FSM
	var ok bool
	if state, ok = s.state.Get(chatID); !ok {
		return customerrors.NewTelegramError("delete email not requested, use /delete")
	}

	if err = state.Event(ctx, StateDeleteEmail); err != nil {
		return customerrors.NewTelegramError(fmt.Errorf("delete email: %w", err).Error())
	}

	if err = s.storage.DeleteEmail(ctx, chatID); err != nil {
		return err
	}

	s.state.Delete(chatID)

	return nil
}
