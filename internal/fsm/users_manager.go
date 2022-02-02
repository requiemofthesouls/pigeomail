package fsm

import "sync"

// User states
const (
	Idle State = iota
	CreatingEmail
	DeletingEmail
)

// User events
const (
	StartCreatingEmail Event = iota
	FinishCreatingEmail
	StartDeletingEmail
	FinishDeletingEmail
	Cancel
)

var (
	// UsersFsm transitions of states for users
	UsersFsm = NewFSM(Transitions{
		Idle: {
			StartCreatingEmail: CreatingEmail,
			StartDeletingEmail: DeletingEmail,
		},
		CreatingEmail: {
			FinishCreatingEmail: Idle,
			Cancel:              Idle,
		},
		DeletingEmail: {
			FinishDeletingEmail: Idle,
			Cancel:              Idle,
		},
	})
)

// UsersManager keeps states of all users and manages it
type UsersManager struct {
	idleState     State // init state for every new user
	Fsm           FSM
	usersStates   map[int64]State
	usersStatesMu sync.RWMutex
}

func NewUserManager() *UsersManager {
	return &UsersManager{
		idleState:   Idle,
		Fsm:         UsersFsm,
		usersStates: make(map[int64]State),
	}
}

// GetState get state for user. If no user found - manager saves user and return idle state
func (um *UsersManager) GetState(userId int64) (state State) {
	um.usersStatesMu.RLock()
	defer um.usersStatesMu.RUnlock()

	if state, ok := um.usersStates[userId]; ok {
		return state
	}
	um.usersStates[userId] = um.idleState
	return um.idleState
}

// SendEventE send event for user. If no user found - user is considered to be in idle state
func (um *UsersManager) SendEventE(userId int64, event Event) (newState State, err error) {
	var currentState = um.GetState(userId)

	um.usersStatesMu.RLock()
	defer um.usersStatesMu.RUnlock()

	if newState, err = um.Fsm.SendEvent(currentState, event); err == nil {
		um.usersStates[userId] = newState
	}
	return newState, err
}

func (um *UsersManager) SendEvent(userId int64, event Event) (newState State) {
	newState, _ = um.SendEventE(userId, event)
	return newState
}
