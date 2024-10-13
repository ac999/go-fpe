// algorithms/strm.go
package algorithms

import (
	"strconv"
)

// Converts an integer to a numeral string of a given length and radix.
func StrmRadix(x uint, radix int, m int) string {
	X := make([]int, m)
	for i := 0; i < m; i++ {
		X[m-1-i] = int(x % uint(radix))
		x = x / uint(radix)
	}
	result := ""
	for _, digit := range X {
		result += strconv.Itoa(digit)
	}
	return result
}
