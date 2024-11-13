// algorithms/helpers.go
package algorithms

import (
	"fmt"
	"strings"
)

// Power - Computes x^y for uint64, avoids overflow
func Power(x, y uint64) uint64 {
	result := uint64(1)
	for i := uint64(0); i < y; i++ {
		result *= x
	}
	return result
}

// NumRadix - Converts a numeral string `X` into an integer based on the radix
func NumRadix(X string, radix int) (uint64, error) {
	var x uint64
	for _, char := range X {
		var digit uint64
		if char >= '0' && char <= '9' {
			digit = uint64(char - '0')
		} else if char >= 'A' && char <= 'Z' {
			digit = uint64(char-'A') + 10
		} else {
			return 0, fmt.Errorf("invalid character '%c' for radix %d", char, radix)
		}

		if digit >= uint64(radix) {
			return 0, fmt.Errorf("invalid digit '%c' for radix %d", char, radix)
		}
		x = x*uint64(radix) + digit
	}
	return x, nil
}

// StrmRadix - Converts integer `x` into a numeral string of length `m` in a given radix
func StrmRadix(x uint64, radix int, m int) (string, error) {
	if x >= Power(uint64(radix), uint64(m)) {
		return "", fmt.Errorf("x (%d) out of bounds for radix^m", x)
	}

	digits := make([]int, m)
	for i := 0; i < m; i++ {
		digits[m-1-i] = int(x % uint64(radix))
		x /= uint64(radix)
	}

	var result strings.Builder
	for _, digit := range digits {
		if digit < 10 {
			result.WriteByte(byte(digit + '0'))
		} else {
			result.WriteByte(byte(digit - 10 + 'A'))
		}
	}
	return result.String(), nil
}
