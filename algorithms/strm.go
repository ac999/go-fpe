// algorithms/strm.go
package algorithms

// Converts an integer to a numeral string based on radix.
func STR(x, m, radix int) []int {
	X := make([]int, m)
	for i := m - 1; i >= 0; i-- {
		X[i] = x % radix
		x = x / radix
	}
	return X
}
