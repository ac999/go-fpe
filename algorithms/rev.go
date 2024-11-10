// algorithms/rev.go
package algorithms

// Rev(X) function implements Algorithm 4: REV(X) from NIST Special Publication 800-38G.
// Reverses a numeral string (handles multi-byte characters as well).
func Rev(X string) string {
	Y := []byte(X) // Convert the string to a rune slice to handle Unicode characters.
	for i := range X {
		Y[i] = X[len(X)-1-i]
	}
	return string(Y)
}

// RevB function implements Algorithm 5: REVB(X) from NIST Special Publication 800-38G.
// Reverses a byte string represented in bits.
func RevB(X []byte) []byte {
	Y := make([]byte, len(X))

	// Reverse the bits in each byte first, then reverse the byte order.
	for i := 0; i < len(X); i++ {
		Y[len(X)-1-i] = reverseBits(X[i])
	}

	return Y
}

// Helper function for RevB
// Reverses the bits in a single byte.
func reverseBits(b byte) byte {
	var rev byte = 0
	for i := 0; i < 8; i++ {
		rev = (rev << 1) | (b & 1)
		b >>= 1
	}
	return rev
}
