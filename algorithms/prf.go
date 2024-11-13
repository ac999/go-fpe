// algorithms/prf.go
package algorithms

import (
	"crypto/aes"
	"crypto/cipher"
)

// PRF - AES-CBC-based pseudorandom function, used in FF1
func PRF(K, X []byte) ([]byte, error) {
	block, err := aes.NewCipher(K)
	if err != nil {
		return nil, err
	}

	// Initialize with all-zero IV
	iv := make([]byte, aes.BlockSize)
	cbc := cipher.NewCBCEncrypter(block, iv)

	// Pad input X to multiple of block size
	if len(X)%aes.BlockSize != 0 {
		padding := aes.BlockSize - (len(X) % aes.BlockSize)
		X = append(X, make([]byte, padding)...)
	}

	Y := make([]byte, len(X))
	cbc.CryptBlocks(Y, X)

	return Y[len(Y)-aes.BlockSize:], nil
}
