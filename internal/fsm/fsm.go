package fsm

import "errors"

type (
	State int
	Event int

	FSM interface {
		SendEvent(currentState State, event Event) (newState State, err error)
	}

	Transitions map[State]map[Event]State

	fsm struct {
		transitions Transitions
	}
)

func NewFSM(transitions Transitions) FSM {
	return &fsm{transitions: transitions}
}

func (m *fsm) SendEvent(currentState State, event Event) (newState State, err error) {
	var ok bool

	var eventsFromInitState map[Event]State
	if eventsFromInitState, ok = m.transitions[currentState]; !ok {
		return 0, errors.New("current state not found")
	}

	if newState, ok = eventsFromInitState[event]; !ok {
		return 0, errors.New("event for current state not found")
	}

	return newState, nil
}
