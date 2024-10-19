// algorithms/prf.go
package algorithms

import (
	"crypto/aes"
	"errors"
)

// PRF function implements Algorithm 6: PRF(X) from NIST Special Publication 800-38G.
// It takes an input block string X and a key K, and returns the result Y.
func PRF(X []byte, K []byte) ([]byte, error) {
	// Ensure the key length is valid (AES-128 requires 16-byte keys).
	block, err := aes.NewCipher(K)
	if err != nil {
		return nil, err
	}

	blockSize := aes.BlockSize // AES block size is always 16 bytes (128 bits)

	// Ensure X is a multiple of the block size.
	if len(X)%blockSize != 0 {
		return nil, errors.New("input length must be a multiple of the AES block size")
	}

	// Number of blocks (m)
	m := len(X) / blockSize

	// Y0 = 0128, i.e., a block of 16 zero bytes.
	Y := make([]byte, blockSize)

	// Iterate over each block Xj
	for j := 0; j < m; j++ {
		Xj := X[j*blockSize : (j+1)*blockSize] // Get block Xj

		// XOR Yj-1 with Xj
		for i := range Y {
			Y[i] ^= Xj[i]
		}

		// Encrypt the result to produce Yj
		block.Encrypt(Y, Y)
	}

	// Return Ym
	return Y, nil
}
