// algorithms/num.go
package algorithms

import (
	"fmt"
	"strconv"
)

// NumRadix(X) function implements Algorithm 1: NUM_radix(X) from NIST Special Publication 800-38G.
// Converts a numeral string to a number based on the given radix.
func NumRadix(X string, radix int) (uint, error) {
	var x uint = 0

	// Validate that the radix is within a reasonable range
	if radix < 2 || radix > 36 {
		return 0, fmt.Errorf("invalid radix: %d. Must be between 2 and 36", radix)
	}

	for i := 0; i < len(X); i++ {
		// Get the value of the current character
		digit, err := strconv.ParseUint(string(X[i]), radix, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid character '%c' for radix %d", X[i], radix)
		}

		x = x*uint(radix) + uint(digit)
	}

	return x, nil
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
