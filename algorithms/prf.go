// algorithms/prf.go
package algorithms

import (
	"crypto/aes"
	"crypto/cipher"
)

// PRF computes a pseudorandom value based on AES using CBC mode.
func PRF(K []byte, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(K)
	if err != nil {
		return nil, err
	}

	// Initialize a zero IV (Initialization Vector)
	iv := make([]byte, aes.BlockSize)

	// Use CBC mode with the zero IV
	mode := cipher.NewCBCEncrypter(block, iv)

	// Pad the input data to be a multiple of block size
	paddedData := padToBlockSize(data, aes.BlockSize)
	ciphertext := make([]byte, len(paddedData))

	// Encrypt using AES-CBC
	mode.CryptBlocks(ciphertext, paddedData)

	// Return the last block of ciphertext as the PRF result
	lastBlock := ciphertext[len(ciphertext)-aes.BlockSize:]
	return lastBlock, nil
}

// Helper to pad input data to the block size
func padToBlockSize(data []byte, blockSize int) []byte {
	padLen := blockSize - (len(data) % blockSize)
	return append(data, make([]byte, padLen)...)
}
