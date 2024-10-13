// algorithms/ff1.go
package algorithms

import (
	"crypto/aes"
	"crypto/cipher"
	"math"
)

// FF1 Encryption Algorithm
func FF1Encrypt(K, T []byte, X string, radix, minlen, maxlen, maxTlen int) string {
	n := len(X)
	if n < minlen || n > maxlen {
		panic("Input length out of range")
	}

	u := n / 2
	v := n - u
	A := X[:u]
	B := X[u:]

	b := int(math.Ceil(float64(v) * math.Log2(float64(radix)) / 8))
	d := 4*int(math.Ceil(float64(b)/4)) + 4

	P := append([]byte{1, 2, 1}, byte(radix>>16), byte(radix>>8), byte(radix), 10, byte(u%256), byte(n>>24), byte(n>>16), byte(n>>8), byte(n), byte(len(T)>>24), byte(len(T)>>16), byte(len(T)>>8), byte(len(T)))

	for i := 0; i < 10; i++ {
		Q := append(T, make([]byte, (-len(T)-b-1)%16)...)
		Q = append(Q, byte(i))
		Q = append(Q, NumRadix(B, radix))

		R := PRF(P, Q, K)
		S := R
		for j := 1; j < int(math.Ceil(float64(d)/16)); j++ {
			block, _ := aes.NewCipher(K)
			block.Encrypt(S, append(R, byte(j)))
		}

		y := NumRadix(string(S[:d]), 2)
		m := u
		if i%2 != 0 {
			m = v
		}

		c := (NumRadix(A, radix) + y) % uint(math.Pow(float64(radix), float64(m)))
		C := StrmRadix(c, radix, m)

		A = B
		B = C
	}

	return A + B
}

// FF1 Decryption Algorithm
func FF1Decrypt(K, T []byte, X string, radix, minlen, maxlen, maxTlen int) string {
	n := len(X)
	if n < minlen || n > maxlen {
		panic("Input length out of range")
	}

	u := n / 2
	v := n - u
	A := X[:u]
	B := X[u:]

	b := int(math.Ceil(float64(v) * math.Log2(float64(radix)) / 8))
	d := 4*int(math.Ceil(float64(b)/4)) + 4

	P := append([]byte{1, 2, 1}, byte(radix>>16), byte(radix>>8), byte(radix), 10, byte(u%256), byte(n>>24), byte(n>>16), byte(n>>8), byte(n), byte(len(T)>>24), byte(len(T)>>16), byte(len(T)>>8), byte(len(T)))

	for i := 9; i >= 0; i-- {
		Q := append(T, make([]byte, (-len(T)-b-1)%16)...)
		Q = append(Q, byte(i))
		Q = append(Q, NumRadix(A, radix))

		R := PRF(P, Q, K)
		S := R
		for j := 1; j < int(math.Ceil(float64(d)/16)); j++ {
			block, _ := aes.NewCipher(K)
			block.Encrypt(S, append(R, byte(j)))
		}

		y := NumRadix(string(S[:d]), 2)
		m := u
		if i%2 == 0 {
			m = v
		}

		c := (NumRadix(B, radix) - y + uint(math.Pow(float64(radix), float64(m)))) % uint(math.Pow(float64(radix), float64(m)))
		C := StrmRadix(c, radix, m)

		B = A
		A = C
	}

	return A + B
}
