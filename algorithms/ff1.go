// algorithms/ff1.go
package algorithms

import "fmt"

func FF1Encrypt(key, tweak []byte, X []int, radix, minLen, maxLen int) ([]int, error) {
	u := len(X) / 2
	v := len(X) - u
	A, B := X[:u], X[u:]

	// Perform 10 Feistel rounds
	for i := 0; i <= 9; i++ {
		y := FeistelRound(key, tweak, B, i, radix)

		if i%2 == 0 {
			m := u
			A = updateNumeral(A, y, m, radix)
		} else {
			m := v
			B = updateNumeral(B, y, m, radix)
		}

		// Swap A and B
		A, B = B, A
	}

	return append(A, B...), nil
}

func FF1Decrypt(key, tweak []byte, X []int, radix, minLen, maxLen int) ([]int, error) {
	u := len(X) / 2
	v := len(X) - u
	A, B := X[:u], X[u:]

	// Perform the Feistel rounds in reverse order
	for i := 9; i >= 0; i-- {
		fmt.Printf("Decrypt Round %d - Start A: %v, B: %v\n", i, A, B)

		y := FeistelRound(key, tweak, A, i, radix)
		fmt.Printf("Decrypt Round %d - Computed y: %d\n", i, y)

		if i%2 == 0 {
			m := u
			B = updateNumeral(B, -y, m, radix)
		} else {
			m := v
			A = updateNumeral(A, -y, m, radix)
		}

		// Swap A and B
		A, B = B, A
		fmt.Printf("Decrypt Round %d - End A: %v, B: %v\n", i, A, B)
	}

	return append(A, B...), nil
}

// Update numeral string based on the calculation in Feistel round
func updateNumeral(A []int, y, m, radix int) []int {
	aNum := NUM(A, radix)
	maxNum := pow(radix, m)

	// Handle addition and subtraction properly for both encryption and decryption
	c := (aNum + y) % maxNum
	if c < 0 {
		c += maxNum // Ensure non-negative result if subtraction goes below 0
	}

	return STR(c, m, radix)
}

// Calculate radix^m
func pow(radix, m int) int {
	result := 1
	for i := 0; i < m; i++ {
		result *= radix
	}
	return result
}
