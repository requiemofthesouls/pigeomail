package fsm

type UsersManager struct {
	idleState   State
	Fsm         FSM
	usersStates map[int64]State
}

func NewUserManager(idleState State, fsm FSM) *UsersManager {
	return &UsersManager{
		idleState:   idleState,
		Fsm:         fsm,
		usersStates: make(map[int64]State),
	}
}

func (um *UsersManager) GetState(userId int64) (state State) {
	if state, ok := um.usersStates[userId]; ok {
		return state
	}
	um.usersStates[userId] = um.idleState
	return um.idleState
}

func (um *UsersManager) SendEvent(userId int64, event Event) (newState State, err error) {
	var currentState = um.GetState(userId)
	if newState, err = um.Fsm.SendEvent(currentState, event); err == nil {
		um.usersStates[userId] = newState
	}
	return newState, err
}
