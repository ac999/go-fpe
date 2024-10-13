// algorithms/rev.go
package algorithms

// Reverses a numeral string.
func Rev(X string) string {
	Y := make([]byte, len(X))
	for i := 0; i < len(X); i++ {
		Y[i] = X[len(X)-1-i]
	}
	return string(Y)
}

// Reverses a byte string represented in bits.
func RevB(X []byte) []byte {
	Y := make([]byte, len(X))
	for i := 0; i < len(X); i++ {
		Y[i] = X[len(X)-1-i]
	}
	return Y
}
