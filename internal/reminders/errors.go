package reminders

import "errors"

var (
	ErrReminderNotFound  = errors.New("reminder not found")
	ErrInvalidReminderID = errors.New("invalid reminder id")
)
