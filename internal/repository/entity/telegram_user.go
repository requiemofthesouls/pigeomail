package entity

import (
	"errors"
	"net/mail"
)

type TelegramUser struct {
	ID     int64
	ChatID int64
	EMail  string
}

func (u *TelegramUser) IsExist() bool {
	return u != nil
}

func (u *TelegramUser) ValidateEMail() (err error) {
	if _, err = mail.ParseAddress(u.EMail); err != nil {
		return errors.New("mail name isn't valid, please choose a new one")
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
