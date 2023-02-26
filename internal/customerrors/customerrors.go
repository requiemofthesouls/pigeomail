package customerrors

import (
	"errors"
)

var ErrNotFound = errors.New("data not found")

// TelegramError stores error which will send to user
type TelegramError struct {
	Err error
}

// NewTelegramError returns new error that you want to send to user
func NewTelegramError(err string) error {
	return &TelegramError{Err: errors.New(err)}
}

func (e *TelegramError) Error() string {
	return e.Err.Error()
}
