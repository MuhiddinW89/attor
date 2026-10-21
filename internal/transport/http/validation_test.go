package transport

import (
	"math"
	"testing"
)

func TestValidPhone(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		valid bool
	}{
		{name: "international", phone: "+1 (555) 123-4567", valid: true},
		{name: "too short", phone: "123456", valid: false},
		{name: "too long", phone: "+123456789012345678901", valid: false},
		{name: "letters", phone: "+1 555 CALL", valid: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := validPhone(test.phone); got != test.valid {
				t.Fatalf("validPhone(%q) = %t, want %t", test.phone, got, test.valid)
			}
		})
	}
}

func TestValidSalePrice(t *testing.T) {
	for _, test := range []struct {
		name  string
		price float64
		valid bool
	}{
		{name: "positive currency", price: 12.34, valid: true},
		{name: "zero", price: 0, valid: false},
		{name: "negative", price: -1, valid: false},
		{name: "too many decimals", price: 1.234, valid: false},
		{name: "nan", price: math.NaN(), valid: false},
		{name: "overflow", price: 10000000000, valid: false},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := validSalePrice(test.price); got != test.valid {
				t.Fatalf("validSalePrice(%v) = %t, want %t", test.price, got, test.valid)
			}
		})
	}
}
