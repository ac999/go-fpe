// prf.go
package algorithms

import (
	"crypto/aes"
	"crypto/cipher"
)

// PRF - AES-CBC-based pseudorandom function for FF1
func PRF(K, X []byte) ([]byte, error) {
	block, err := aes.NewCipher(K)
	if err != nil {
		return nil, err
	}

	// Use an all-zero IV for CBC mode
	iv := make([]byte, aes.BlockSize)
	cbc := cipher.NewCBCEncrypter(block, iv)

	// Pad X to be a multiple of the block size
	padding := aes.BlockSize - (len(X) % aes.BlockSize)
	paddedX := append(X, make([]byte, padding)...)

	Y := make([]byte, len(paddedX))
	cbc.CryptBlocks(Y, paddedX)

	// Return only the last block of Y
	return Y[len(Y)-aes.BlockSize:], nil
}
