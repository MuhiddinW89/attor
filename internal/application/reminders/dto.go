package reminders

import "time"

type ReminderListItem struct {
	ID          string    `json:"id"`
	ClientID    string    `json:"client_id"`
	ClientName  string    `json:"client_name"`
	Phone       string    `json:"phone"`
	PerfumeName string    `json:"perfume_name"`
	VolumeML    int       `json:"volume_ml"`
	Price       float64   `json:"price"`
	Comment     *string   `json:"comment"`
	SaleDate    time.Time `json:"sale_date"`
	ReminderAt  time.Time `json:"reminder_at"`
	IsSent      bool      `json:"is_sent"`
}