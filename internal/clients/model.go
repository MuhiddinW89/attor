package clients

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID        uuid.UUID  `json:"id"`
	FullName  string     `json:"full_name"`
	Phone     string     `json:"phone"`
	Instagram *string    `json:"instagram,omitempty"`
	BirthDate *time.Time `json:"birthDate,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}
