// algorithms/ff3.go
package algorithms

import "fmt"

func FF3Encrypt(key, tweak []byte, X []int, radix int) ([]int, error) {
	if len(tweak) != 8 {
		return nil, fmt.Errorf("FF3 requires a 64-bit (8-byte) tweak")
	}

	u := (len(X) + 1) / 2
	v := len(X) - u
	A, B := X[:u], X[u:]

	TL := tweak[:4] // Left tweak
	TR := tweak[4:] // Right tweak

	for i := 0; i < 8; i++ {
		W := TL
		if i%2 == 0 {
			W = TR
		}

		y := FeistelRound(key, W, B, i, radix)

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

func FF3Decrypt(key, tweak []byte, X []int, radix int) ([]int, error) {
	if len(tweak) != 8 {
		return nil, fmt.Errorf("FF3 requires a 64-bit (8-byte) tweak")
	}

	u := (len(X) + 1) / 2
	v := len(X) - u
	A, B := X[:u], X[u:]

	TL := tweak[:4] // Left tweak
	TR := tweak[4:] // Right tweak

	// Perform the Feistel rounds in reverse order
	for i := 7; i >= 0; i-- {
		W := TL
		if i%2 == 0 {
			W = TR
		}

		y := FeistelRound(key, W, A, i, radix)

		if i%2 == 0 {
			m := u
			B = updateNumeral(B, -y, m, radix)
		} else {
			m := v
			A = updateNumeral(A, -y, m, radix)
		}

		// Swap A and B
		A, B = B, A
	}

	return append(A, B...), nil
}
