package transport

import (
	"math"
	"strings"
	"unicode/utf8"
)

const maxVolumeML = 2147483647

func runeLength(value string) int {
	return utf8.RuneCountInString(value)
}

func validName(value string, maxLength int) bool {
	trimmed := strings.TrimSpace(value)
	return trimmed != "" && runeLength(trimmed) <= maxLength
}

func validPhone(value string) bool {
	value = strings.TrimSpace(value)
	if runeLength(value) < 7 || runeLength(value) > 20 {
		return false
	}

	digits := 0
	for index, character := range value {
		switch {
		case character >= '0' && character <= '9':
			digits++
		case character == '+' && index == 0:
		case character == ' ', character == '-', character == '(', character == ')', character == '.':
		default:
			return false
		}
	}

	return digits >= 7
}

func validSalePrice(price float64) bool {
	if math.IsNaN(price) || math.IsInf(price, 0) || price <= 0 || price > 9999999999.99 {
		return false
	}
	return math.Abs(price*100-math.Round(price*100)) < 0.000001
}
