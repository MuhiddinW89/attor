package application

import "time"

type ReminderListItem struct {
	ID          string    `json:"id"`
	ClientID    string    `json:"client_id"`
	ClientName  string    `json:"client_name"`
	Phone       string    `json:"phone"`
	PerfumeName string    `json:"perfume_name"`
	VolumeML    int       `json:"volume_ml"`
	ReminderAt  time.Time `json:"reminder_at"`
}
