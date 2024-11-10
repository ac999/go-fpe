// algorithms/num.go
package algorithms

import (
	"fmt"
	"strconv"
)

// NumRadix(X) function implements Algorithm 1: NUM_radix(X) from NIST Special Publication 800-38G.
// Converts a numeral string to a number based on the given radix.
func NumRadix(X string, radix int) (uint64, error) {
	var x uint64 = 0

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

		x = x*uint64(radix) + digit
	}

	return x, nil
}

// NumBits converts a binary string X, represented in bits, to an integer.
func NumBits(X string) (uint64, error) {
	var x uint64 = 0
	for i := 0; i < len(X); i++ {
		// Check for valid binary characters.
		if X[i] != '0' && X[i] != '1' {
			return 0, fmt.Errorf("invalid character '%c' in binary string", X[i])
		}

		// Convert '0' or '1' to integer 0 or 1 by subtracting '0'.
		bit := X[i] - '0'
		x = 2*x + uint64(bit)
	}

	return x, nil
}
