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

	t := len(T)

	// Steps 3 & 4: Calculate the byte length b and block length d
	b := int(math.Ceil(float64(v) * math.Log2(float64(radix)) / 8))
	if b <= 0 {
		return "", errors.New("invalid block length b")
	}
	d := 4*int(math.Ceil(float64(b)/4)) + 4

	// Step 5: Construct P = [1, 2, 1, radix (3 bytes), 10, u mod 256, n (4 bytes), t (4 bytes)]
	radix3bytes := intToNBytes(radix, 3)
	nBytes := intToNBytes(n, 4)
	tBytes := intToNBytes(t, 4)
	P := []byte{1, 2, 1}
	P = append(P, radix3bytes...)
	P = append(P, 10, byte(u%256))
	P = append(P, nBytes...)
	P = append(P, tBytes...)

	// Step 6: Perform the Feistel rounds
	for i := 0; i < 10; i++ {
		// Step 6.i: Construct Q safely
		padding := 16 + (0-t-b-1)%16
		Q := append(T, make([]byte, padding)...)
		Q = append(Q, byte(i)) // append round number i

		// Convert B to a numeric representation
		BNum, err := NumRadix(B, radix)
		if err != nil {
			return "", err
		}
		bNumBytes := uint64ToNBytes(BNum, b)
		Q = append(Q, bNumBytes...)

		//  Check size of Q to be multiple of 16 before calling PRF(P, Q)
		switch len(Q) % 16 {
		default:
			error_str := fmt.Sprint("invalid Q size: ", len(Q), "; must be 16, 24, or 32 bytes")
			return "", errors.New(error_str)
		case 0:
		}

		// Step 6.ii: Compute R = PRF(P, Q)
		R, err := PRF(append(P, Q...), K)
		if err != nil {
			log.Println("P is: ", P, " of length: ", len(P))
			log.Println("Q is: ", Q, " of length: ", len(Q))
			return "", err
		}

		// Step 6.iii: Let S be the first d bytes of CIPHK(R ⊕ [j]16)
		if d <= 0 {
			return "", errors.New("invalid block length d")
		}

		block, err := aes.NewCipher(K)
		if err != nil {
			return "", err
		}

		blockSize := aes.BlockSize // 16 bytes for AES

		dpe16 := int(math.Ceil(float64(d) / float64(blockSize)))
		S := make([]byte, len(R))
		copy(S, R)
		for j := 1; j < dpe16; j++ {
			j16Bytes := intToNBytes(j, 16)
			res := make([]byte, len(R))
			// R xor j
			for k, r := range R {
				res[k] = r ^ j16Bytes[k]
			}
			aux := make([]byte, len(R))
			block.Encrypt(aux, res)

			S = append(S, aux...)
		}

		// for j := 0; j < int(math.Ceil(float64(d)/float64(blockSize))); j++ {
		// 	Rj := make([]byte, aes.BlockSize)
		// 	copy(Rj, append(R, byte(j)))
		// 	encrypted := make([]byte, blockSize)
		// 	block.Encrypt(encrypted, Rj)
		// 	S = append(S, encrypted...)
		// }
		// S = S[:d] // truncate to d bytes

		// Step 6.iv: Let y = NUM(S)
		y, err := NumBits(string(S))
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
		c := (aNum + y) % power(uint64(radix), uint64(m))

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

// FF1 Decryption Algorithm
func FF1Decrypt(K, T []byte, X string, radix, minlen, maxlen, maxTlen int) (string, error) {
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

	t := len(T)

	// Steps 3 & 4: Calculate the byte length b and block length d
	b := int(math.Ceil(float64(v) * math.Log2(float64(radix)) / 8))
	if b <= 0 {
		return "", errors.New("invalid block length b")
	}
	d := 4*int(math.Ceil(float64(b)/4)) + 4

	// Step 5: Construct P = [1, 2, 1, radix (3 bytes), 10, u mod 256, n (4 bytes), t (4 bytes)]
	radix3bytes := intToNBytes(radix, 3)
	nBytes := intToNBytes(n, 4)
	tBytes := intToNBytes(t, 4)
	P := []byte{1, 2, 1}
	P = append(P, radix3bytes...)
	P = append(P, 10, byte(u%256))
	P = append(P, nBytes...)
	P = append(P, tBytes...)

	// Step 6: Perform the Feistel rounds
	for i := 9; i >= 0; i-- {
		// Step 6.i: Construct Q safely
		padding := 16 + (0-t-b-1)%16
		Q := append(T, make([]byte, padding)...)
		Q = append(Q, byte(i)) // append round number i

		// Convert B to a numeric representation
		ANum, err := NumRadix(A, radix)
		if err != nil {
			return "", err
		}
		bNumBytes := uint64ToNBytes(ANum, b)
		Q = append(Q, bNumBytes...)

		//  Check size of Q to be multiple of 16 before calling PRF(P, Q)
		switch len(Q) % 16 {
		default:
			error_str := fmt.Sprint("invalid Q size: ", len(Q), "; must be 16, 24, or 32 bytes")
			return "", errors.New(error_str)
		case 0:
		}

		// Step 6.ii: Compute R = PRF(P, Q)
		R, err := PRF(append(P, Q...), K)
		if err != nil {
			log.Println("P is: ", P, " of length: ", len(P))
			log.Println("Q is: ", Q, " of length: ", len(Q))
			return "", err
		}

		// Step 6.iii: Let S be the first d bytes of CIPHK(R ⊕ [j]16)
		if d <= 0 {
			return "", errors.New("invalid block length d")
		}

		block, err := aes.NewCipher(K)
		if err != nil {
			return "", err
		}

		blockSize := aes.BlockSize // 16 bytes for AES

		dpe16 := int(math.Ceil(float64(d) / float64(blockSize)))
		S := make([]byte, len(R))
		copy(S, R)
		for j := 1; j < dpe16; j++ {
			j16Bytes := intToNBytes(j, 16)
			res := make([]byte, len(R))
			// R xor j
			for k, r := range R {
				res[k] = r ^ j16Bytes[k]
			}
			aux := make([]byte, len(R))
			block.Decrypt(aux, res)

			S = append(S, aux...)
		}

		// for j := 0; j < int(math.Ceil(float64(d)/float64(blockSize))); j++ {
		// 	Rj := make([]byte, aes.BlockSize)
		// 	copy(Rj, append(R, byte(j)))
		// 	encrypted := make([]byte, blockSize)
		// 	block.Encrypt(encrypted, Rj)
		// 	S = append(S, encrypted...)
		// }
		// S = S[:d] // truncate to d bytes

		// Step 6.iv: Let y = NUM(S)
		y, err := NumBits(string(S))
		if err != nil {
			return "", err
		}

		// Step 6.v: Determine m
		m := u
		if i%2 != 0 {
			m = v
		}

		// Step 6.vi: Compute c = (NUMradix(A) + y) mod radix^m
		bNum, err := NumRadix(B, radix)
		if err != nil {
			return "", err
		}
		c := (bNum - y) % power(uint64(radix), uint64(m))

		// Step 6.vii: Let C = STRmradix(c)
		C, err := StrmRadix(c, radix, m)
		if err != nil {
			return "", err
		}

		// Step 6.viii and 6.ix: Swap A and B for next round
		B = A
		A = C
	}

	// Step 7: Return A || B as the encrypted result
	return A + B, nil
}
