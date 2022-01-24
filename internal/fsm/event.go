package fsm

import "context"

type (
	Event interface {
		GetName() string
		Process(ctx context.Context) error
	}

	EventName string
)
