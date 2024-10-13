// algorithms/num.go
package algorithms

import (
	"fmt"
	"strconv"
)

// Converts a numeral string to a number based on the given radix.
func NumRadix(X string, radix int) uint {
	var x uint = 0
	for i := 0; i < len(X); i++ {
		digit, err := strconv.Atoi(string(X[i]))
		if err != nil {
			fmt.Println("Error converting character to digit:", err)
			return 0
		}
		x = x*uint(radix) + uint(digit)
	}
	return x
}

// Converts a byte string represented in bits to an integer.
func NumBits(X string) uint {
	var x uint = 0
	for i := 0; i < len(X); i++ {
		bit := X[i] - '0' // Convert character '0' or '1' to integer 0 or 1
		x = 2*x + uint(bit)
	}
	return x
}
