package transport

type CreateSaleRequest struct {
	ClientID    string  `json:"client_id"`
	Phone       string  `json:"phone"`
	FullName    string  `json:"full_name"`
	Birthday    string  `json:"birthday"`
	PerfumeName string  `json:"perfume_name"`
	VolumeML    int     `json:"volume_ml"`
	Price       float64 `json:"price"`
	Comment     *string `json:"comment"`
}
