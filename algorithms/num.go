// algorithms/num.go
package algorithms

import (
	"fmt"
	"strconv"
)

// NumRadix(X) function implements Algorithm 1: NUM_radix(X) from NIST Special Publication 800-38G.
// Converts a numeral string to a number based on the given radix.
func NumRadix(X string, radix int) (uint64, error) {
	var x uint64
	if radix < 2 || radix > 36 {
		return 0, fmt.Errorf("invalid radix: %d. Must be between 2 and 36", radix)
	}
	for _, digit := range X {
		val, err := strconv.ParseUint(string(digit), radix, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid character '%c' for radix %d", digit, radix)
		}
		x = x*uint64(radix) + val
	}
	return x, nil
}

// NumBits(X) function implements Algorithm 2: NUM(X) from NIST Special Publication 800-38G.
// Converts a byte string represented in bits to an integer.
// NumBits now converts binary numeral strings properly to uint64
func NumBits(X string) (uint64, error) {
	var x uint64
	for _, bit := range X {
		if bit != '0' && bit != '1' {
			return 0, fmt.Errorf("invalid character '%c' in binary string", bit)
		}
		x = 2*x + uint64(bit-'0')
	}
	return x, nil
}
