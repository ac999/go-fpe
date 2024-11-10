// algorithms/strm.go
package algorithms

import (
	"fmt"
)

// StrmRadix converts an integer to a numeral string of a given length and radix,
// as per the NIST specification.

func StrmRadix(x uint64, radix int, m int) (string, error) {
	// Validate that the input x is within the valid range.
	maxValue := power(uint64(radix), uint64(m))
	if x >= maxValue {
		return "", fmt.Errorf("invalid input: x (%d) must be less than radix^m (%d)", x, maxValue)
	}

	// Generate the numeral string in reverse order.
	X := make([]int, m)
	for i := 0; i < m; i++ {
		X[m-1-i] = int(x % uint64(radix)) // X[m+1–i] = x mod radix
		x = x / uint64(radix)             // x = floor(x / radix)
	}

	// Build the result string.
	result := ""
	for _, digit := range X {
		result += fmt.Sprintf("%d", digit)
	}
	return result, nil
}
