// component.go
package algorithms

import (
	"crypto/aes"
	"crypto/cipher"
)

// NUM - Bit string to uint64 - could be faster by using bit shifting
func NUM(X []byte) uint64 {
	x := uint64(0)
	for _, i := range X {
		x = 2*x + uint64(i)
	}
	return x
}

func NUMradix(X []byte, radix uint64) uint64 {
	x := uint64(0)
	for _, i := range X {
		x = x*radix + uint64(i)
	}
	return x
}

func STRmRadix(x uint64, radix uint64, m int) []byte {
	X := make([]byte, m)
	for i := range m {
		X[m-1-i] = byte(Mod(uint64(x), radix))
		x = x / radix
	}
	return X
}

// PRF - AES-CBC-based pseudorandom function for FF1
func PRF(K, X []byte) ([]byte, error) {

	// CIPHk initialization
	block, err := aes.NewCipher(K)
	if err != nil {
		return nil, err
	}

	// Use an all-zero IV for CBC mode
	iv := make([]byte, aes.BlockSize)
	cbc := cipher.NewCBCEncrypter(block, iv)

	m := len(X) / 128

	XBlocks, err := BreakInBlocks(X, 128)
	if err != nil {
		return nil, err
	}

	Y := make([][]byte, m)
	// Initialize each inner slice with a length of 128
	for i := range Y {
		Y[i] = make([]byte, 128)
	}

	// Create storage for Y[j-1] xor X[j] becayse XORBytes could return error
	var xorValue []byte

	for j := 1; j < m; j++ {
		xorValue, err = XORBytes(Y[j-1], XBlocks[j])
		if err != nil {
			return nil, err
		}
		cbc.CryptBlocks(Y[j], xorValue)
	}

	// Return only the last block of Y
	return Y[m-1], nil
}
