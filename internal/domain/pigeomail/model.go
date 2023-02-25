package pigeomail

import (
	"net/mail"

	customerrors "pigeomail/internal/errors"
)

type EMail struct {
	ID     string `bson:"_id,omitempty"`
	ChatID int64  `json:"chat_id" bson:"chat_id"`
	Name   string `json:"name" bson:"name"`
}

func (e *EMail) Validate() (err error) {
	if _, err = mail.ParseAddress(e.Name); err != nil {
		return customerrors.NewTelegramError("mail name isn't valid, please choose a new one")
	}

	return nil
}

const StateRequestedCreateEmail = "requested_create_email"
const StateCreateEmail = "create_email"
const StateEmailCreated = "email_created"

const StateRequestedDeleteEmail = "requested_delete_email"
const StateDeleteEmail = "delete_email"
const StateCancelDeleteEmail = "cancel_delete_email"
const StateDeleteEmailCancelled = "delete_email_cancelled"
const StateEmailDeleted = "email_deleted"
