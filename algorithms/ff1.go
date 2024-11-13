// algorithms/ff1.go

package algorithms

import (
	"fmt"
	"math"
)

// FF1Encrypt - Format-preserving encryption using FF1
func FF1Encrypt(K, T []byte, X string, radix, minlen, maxlen, maxTlen int) (string, error) {
	n := len(X)
	if n < minlen || n > maxlen {
		return "", fmt.Errorf("input length %d out of range", n)
	}

	u := n / 2
	v := n - u
	A := X[:u]
	B := X[u:]

	// P construction
	b := (v * int(math.Log2(float64(radix)))) / 8
	P := constructP(radix, n, u, len(T))

	// Feistel rounds
	for i := 0; i < 10; i++ {
		Q := constructQ(T, B, i, b)
		R, err := PRF(append(P, Q...), K)
		if err != nil {
			return "", err
		}
		y := NumBits(R)

		m := u
		if i%2 != 0 {
			m = v
		}

		aNum := NumRadix(A, radix)
		c := (aNum + y) % Power(uint64(radix), uint64(m))
		C := StrmRadix(c, radix, m)

		A, B = B, C
	}

	return A + B, nil
}
