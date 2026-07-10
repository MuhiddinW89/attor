package sales

import "errors"

var (
	ErrSaleNotFound        = errors.New("sale not found")
	ErrInvalidSalePrice    = errors.New("price must be greater than zero")
	ErrInvalidVolume       = errors.New("volume must be greater than zero")
	ErrPerfumeNameRequired = errors.New("perfume name is required")
)
