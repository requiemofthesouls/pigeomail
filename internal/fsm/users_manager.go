package fsm

// User states
const (
	Idle State = iota
	ChoosingEmail
	DeletingEmail
)

// User events
const (
	CreateEmail Event = iota
	ChooseEmail
	DeleteEmail
	ConfirmDeletion
	Cancel
)

var (
	// UsersFsm transitions of states for users
	UsersFsm = NewFSM(Transitions{
		Idle: {
			CreateEmail: ChoosingEmail,
			DeleteEmail: DeletingEmail,
		},
		ChoosingEmail: {
			ChooseEmail: Idle,
			Cancel:      Idle,
		},
		DeletingEmail: {
			ConfirmDeletion: Idle,
			Cancel:          Idle,
		},
	})
)

// UsersManager keeps states of all users and manages it
type UsersManager struct {
	idleState   State // init state for every new user
	Fsm         FSM
	usersStates map[int64]State
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
	if state, ok := um.usersStates[userId]; ok {
		return state
	}
	um.usersStates[userId] = um.idleState
	return um.idleState
}

// SendEventE send event for user. If no user found - user is considered to be in idle state
func (um *UsersManager) SendEventE(userId int64, event Event) (newState State, err error) {
	var currentState = um.GetState(userId)
	if newState, err = um.Fsm.SendEvent(currentState, event); err == nil {
		um.usersStates[userId] = newState
	}
	return newState, err
}

func (um *UsersManager) SendEvent(userId int64, event Event) (newState State) {
	newState, _ = um.SendEventE(userId, event)
	return newState
}
