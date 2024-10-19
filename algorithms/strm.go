// algorithms/strm.go
package algorithms

import (
	"fmt"
	"strings"
)

// StrmRadix(X) function implements Algorithm 3: STR^m_radix(X) from NIST Special Publication 800-38G.
// Converts an integer to a numeral string of a given length and radix.
func StrmRadix(x uint, radix int, m int) (string, error) {
	// Validate that the input x is within the valid range
	maxValue := uint(1)
	for i := 0; i < m; i++ {
		maxValue *= uint(radix)
	}
	if x >= maxValue {
		return "", fmt.Errorf("invalid input: x (%d) must be less than radix^m (%d)", x, maxValue)
	}

	X := make([]int, m)
	for i := 0; i < m; i++ {
		X[m-1-i] = int(x % uint(radix)) // X[m+1–i] = x mod radix
		x = x / uint(radix)             // x = floor(x / radix)
	}

	// Efficiently build the result string
	var result strings.Builder
	for _, digit := range X {
		if digit < 10 {
			result.WriteByte(byte(digit + '0')) // 0-9 are represented as '0'-'9'
		} else {
			result.WriteByte(byte(digit - 10 + 'A')) // 10-35 are represented as 'A'-'Z'
		}
	}
	return result.String(), nil
}
