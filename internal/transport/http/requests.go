package transport

type CreateSaleRequest struct {
	Phone       string  `json:"phone"`
	FullName    string  `json:"full_name"`
	PerfumeName string  `json:"perfume_name"`
	VolumeML    int     `json:"volume_ml"`
	Price       float64 `json:"price"`
	Comment     *string `json:"comment"`
}
