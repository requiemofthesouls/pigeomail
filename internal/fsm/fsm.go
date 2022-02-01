package fsm

import "errors"

type (
	State int
	Event int

	// FSM simple fsm, doesn't keep state, instead of it receives current state and event to get new state
	FSM interface {
		SendEvent(currentState State, event Event) (newState State, err error)
	}

	// Transitions map of state transitions: current state -> event -> new state
	Transitions map[State]map[Event]State

	fsm struct {
		transitions Transitions
	}
)

var (
	ErrCurrentStateNotFound         = errors.New("current state not found")
	ErrEventForCurrentStateNotFound = errors.New("event for current state not found")
)

func NewFSM(transitions Transitions) FSM {
	return &fsm{transitions: transitions}
}

func (m *fsm) SendEvent(currentState State, event Event) (newState State, err error) {
	var ok bool

	var eventsFromCurrentState map[Event]State
	if eventsFromCurrentState, ok = m.transitions[currentState]; !ok {
		return 0, ErrCurrentStateNotFound
	}

	if newState, ok = eventsFromCurrentState[event]; !ok {
		return 0, ErrEventForCurrentStateNotFound
	}

	return newState, nil
}
