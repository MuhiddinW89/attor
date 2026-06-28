package sales

import (
	"time"

	"github.com/google/uuid"
)

type Sale struct {
	ID          uuid.UUID
	ClientID    uuid.UUID
	PerfumeName string
	VolumeML    int
	Price       float64
	Comment     *string
	SaleDate    time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
