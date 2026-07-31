package transport

import "time"

type SaleResponse struct {
	ID          string    `json:"id"`
	ClientID    string    `json:"client_id"`
	PerfumeName string    `json:"perfume_name"`
	VolumeML    int       `json:"volume_ml"`
	Price       float64   `json:"price"`
	Comment     *string   `json:"comment,omitempty"`
	SaleDate    time.Time `json:"sale_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
