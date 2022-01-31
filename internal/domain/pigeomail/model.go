package pigeomail

import (
	"errors"
	"net/mail"
)

type EMail struct {
	ID     string `bson:"_id,omitempty"`
	ChatID int64  `json:"chat_id" bson:"chat_id"`
	Name   string `json:"name" bson:"name"`
}

func (e *EMail) Validate() (err error) {
	if _, err = mail.ParseAddress(e.Name); err != nil {
		return errors.New("mail name isn't valid, please choose a new one")
	}

	return nil
}

type UserState struct {
	ID     string `bson:"_id,omitempty"`
	ChatID int64  `json:"chat_id" bson:"chat_id"`
	State  state  `json:"state" bson:"state"`
}

type state string

const StateCreateEmailStep1 state = "create_email_step_1"
const StateDeleteEmailStep1 state = "delete_email_step_1"
