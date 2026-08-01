package reminders

import "time"

type ReminderListItem struct {
	ID string
	ClientID string
	ClientName string
	Phone string
	PerfumeName string
	VolumeML int
	Comment string
	ReminderAt time.Time
}
