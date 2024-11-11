// algorithms/strm.go
package algorithms

import (
	"fmt"
	"strings"
)

// StrmRadix(X) function implements Algorithm 3: STR^m_radix(X) from NIST Special Publication 800-38G.
// Converts an integer to a numeral string of a given length and radix.
func StrmRadix(x uint64, radix int, m int) (string, error) {
	maxValue := power(uint64(radix), uint64(m))
	if x >= maxValue {
		return "", fmt.Errorf("input x (%d) exceeds radix^m (%d)", x, maxValue)
	}

	X := make([]int, m)
	for i := 0; i < m; i++ {
		X[m-1-i] = int(x % uint64(radix))
		x /= uint64(radix)
	}

	var result strings.Builder
	for _, digit := range X {
		if digit < 10 {
			result.WriteByte(byte(digit + '0'))
		} else {
			result.WriteByte(byte(digit - 10 + 'A'))
		}
	}
	return result.String(), nil
}
