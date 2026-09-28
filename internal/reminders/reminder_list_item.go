package reminders

import (
    "time"

    "github.com/google/uuid"
)

type ReminderListItem struct {
    ID          uuid.UUID
    ClientID    uuid.UUID
    ClientName  string
    Phone       string
    PerfumeName string
    VolumeML    int
    Price       float64
    Comment     string
    SaleDate    time.Time
    ReminderAt  time.Time
    IsSent      bool
}