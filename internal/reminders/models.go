package reminders

import (
	"time"

	"github.com/google/uuid"
)

type Reminder struct {
	ID          uuid.UUID
	ClientID    uuid.UUID
	SaleID      uuid.UUID
	ReminderAt  time.Time
	IsSent	    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}