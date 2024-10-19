// algorithms/ff1.go
package algorithms

import (
	"crypto/aes"
	"errors"
	"fmt"
	"log"
	"math"
)

// FF1Encrypt performs format-preserving encryption using FF1
func FF1Encrypt(K, T []byte, X string, radix, minlen, maxlen, maxTlen int) (string, error) {
	// Validate input length
	n := len(X)
	if n < minlen || n > maxlen {
		return "", errors.New("input length out of range")
	}

	// Validate the key length as in crypto/aes/cipher.go
	switch len(K) {
	default:
		return "", errors.New("invalid AES key size; must be 16, 24, or 32 bytes")
	case 16, 24, 32:
		break
	}

	// Steps 1 & 2: Split the numeral string X into A and B
	u := n / 2
	v := n - u
	A := X[:u]
	B := X[u:]

	// Steps 3 & 4: Calculate the byte length b and block length d
	b := int(math.Ceil(float64(v) * math.Log2(float64(radix)) / 8))
	if b <= 0 {
		return "", errors.New("invalid block length b")
	}
	d := 4*int(math.Ceil(float64(b)/4)) + 4

	// Step 5: Construct P = [1, 2, 1, radix (3 bytes), 10, u mod 256, n (4 bytes), t (4 bytes)]
	P := append([]byte{1, 2, 1}, byte(radix>>16), byte(radix>>8), byte(radix))
	P = append(P, 10, byte(u%256))
	P = append(P, byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
	P = append(P, byte(len(T)>>24), byte(len(T)>>16), byte(len(T)>>8), byte(len(T)))

	// Step 6: Perform the Feistel rounds
	for i := 0; i < 10; i++ {
		// Step 6.i: Construct Q safely
		padding := (16 - (len(T)+b+1)%16) % 16
		Q := append(T, make([]byte, padding)...)
		Q = append(Q, byte(i)) // append round number i

		// Convert B to a numeric representation
		BNum, err := NumRadix(B, radix)
		if err != nil {
			return "", err
		}
		Q = append(Q, UintToBytes(BNum)...)

		paddingLen := 16 - (len(Q) % 16)
		Q = append(Q, make([]byte, paddingLen)...)

		//  Check size of Q to be multiple of 16 before calling PRF(P, Q)
		switch len(Q) % 16 {
		default:
			error_str := fmt.Sprint("invalid Q size: ", len(Q), "; must be 16, 24, or 32 bytes")
			return "", errors.New(error_str)
		case 0:
		}

		// Step 6.ii: Compute PRF(P, Q) (pass P and Q as separate arguments)
		R, err := PRF(P, Q) // Pass P and Q as separate arguments
		if err != nil {
			log.Println("P is: ", P, " of length: ", len(P))
			log.Println("Q is: ", Q, " of length: ", len(Q))
			return "", err
		}

		// Step 6.iii: Let S be the first d bytes of CIPHK(R ⊕ [j]16)
		if d <= 0 {
			return "", errors.New("invalid block length d")
		}
		S := make([]byte, 0, d)
		block, err := aes.NewCipher(K)
		if err != nil {
			return "", err
		}
		blockSize := aes.BlockSize // 16 bytes for AES
		for j := 0; j < int(math.Ceil(float64(d)/float64(blockSize))); j++ {
			Rj := make([]byte, aes.BlockSize)
			copy(Rj, append(R, byte(j)))
			encrypted := make([]byte, blockSize)
			block.Encrypt(encrypted, Rj)
			S = append(S, encrypted...)
		}
		S = S[:d] // truncate to d bytes
		// Step 6.iv: Let y = NUM(S)
		y, err := NumRadix(string(S), 2)
		if err != nil {
			return "", err
		}

		// Step 6.v: Determine m
		m := u
		if i%2 != 0 {
			m = v
		}

		// Step 6.vi: Compute c = (NUMradix(A) + y) mod radix^m
		aNum, err := NumRadix(A, radix)
		if err != nil {
			return "", err
		}
		c := (aNum + y) % uint(math.Pow(float64(radix), float64(m)))

		// Step 6.vii: Let C = STRmradix(c)
		C, err := StrmRadix(c, radix, m)
		if err != nil {
			return "", err
		}

		// Step 6.viii and 6.ix: Swap A and B for next round
		A = B
		B = C
	}

	// Step 7: Return A || B as the encrypted result
	return A + B, nil
}

// // FF1 Decryption Algorithm
// func FF1Decrypt(K, T []byte, X string, radix, minlen, maxlen, maxTlen int) string {
// 	n := len(X)
// 	if n < minlen || n > maxlen {
// 		panic("Input length out of range")
// 	}

// 	u := n / 2
// 	v := n - u
// 	A := X[:u]
// 	B := X[u:]

// 	b := int(math.Ceil(float64(v) * math.Log2(float64(radix)) / 8))
// 	d := 4*int(math.Ceil(float64(b)/4)) + 4

// 	P := append([]byte{1, 2, 1}, byte(radix>>16), byte(radix>>8), byte(radix), 10, byte(u%256), byte(n>>24), byte(n>>16), byte(n>>8), byte(n), byte(len(T)>>24), byte(len(T)>>16), byte(len(T)>>8), byte(len(T)))

// 	for i := 9; i >= 0; i-- {
// 		Q := append(T, make([]byte, (-len(T)-b-1)%16)...)
// 		Q = append(Q, byte(i))
// 		Q = append(Q, NumRadix(A, radix))

// 		R := PRF(P, Q, K)
// 		S := R
// 		for j := 1; j < int(math.Ceil(float64(d)/16)); j++ {
// 			block, _ := aes.NewCipher(K)
// 			block.Encrypt(S, append(R, byte(j)))
// 		}

// 		y := NumRadix(string(S[:d]), 2)
// 		m := u
// 		if i%2 == 0 {
// 			m = v
// 		}

// 		c := (NumRadix(B, radix) - y + uint(math.Pow(float64(radix), float64(m)))) % uint(math.Pow(float64(radix), float64(m)))
// 		C := StrmRadix(c, radix, m)

// 		B = A
// 		A = C
// 	}

// 	return A + B
// }
