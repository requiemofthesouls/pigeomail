package fsm

import (
	"context"
	"errors"
)

type (
	StateManager interface {
		GetState() State
		SendEvent(ctx context.Context, event Event) (State, error)
	}

	transitions map[State]map[EventName]State

	stateManager struct {
		state       State
		transitions transitions
	}
)

func NewStateManager(
	initState State,
	transitions transitions,
) StateManager {
	return &stateManager{
		state:       initState,
		transitions: transitions,
	}
}

func (sm *stateManager) GetState() State {
	return sm.state
}

func (sm *stateManager) SendEvent(ctx context.Context, event Event) (newState State, err error) {
	var stateFound bool
	if stateFound, newState = getNewStateByEvent(sm.transitions, sm.state, EventName(event.GetName())); !stateFound {
		return sm.state, errors.New("incorrect event")
	}

	if err = event.Process(ctx); err != nil {
		return sm.state, err
	}

	return newState, nil
}

func getNewStateByEvent(
	transitions transitions,
	currentState State,
	eventName EventName,
) (ok bool, newState State) {
	var transitionsForState map[EventName]State
	if transitionsForState, ok = transitions[currentState]; !ok {
		return false, ""
	}

	if newState, ok = transitionsForState[eventName]; !ok {
		return false, ""
	}

	return true, newState
}
