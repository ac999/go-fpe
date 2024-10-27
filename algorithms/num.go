// algorithms/num.go
package algorithms

import (
	"fmt"
)

// Converts a numeral string to an integer based on radix.
func NUM(X []int, radix int) int {
	result := 0
	for _, numeral := range X {
		result = result*radix + numeral
	}
	return result
}

// NumBits(X) function implements Algorithm 2: NUM(X) from NIST Special Publication 800-38G.
// Converts a byte string represented in bits to an integer.
func NumBits(X string) (uint, error) {
	var x uint = 0

	for i := 0; i < len(X); i++ {
		if X[i] != '0' && X[i] != '1' {
			return 0, fmt.Errorf("invalid character '%c' in binary string", X[i])
		}
		bit := X[i] - '0' // Convert character '0' or '1' to integer 0 or 1
		x = 2*x + uint(bit)
	}

	return x, nil
}
