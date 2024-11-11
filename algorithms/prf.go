// algorithms/prf.go
package algorithms

import (
	"crypto/aes"
	"errors"
)

// PRF function adjusted to handle padding and block processing for AES-CBC PRF
func PRF(X, K []byte) ([]byte, error) {
	block, err := aes.NewCipher(K)
	if err != nil {
		return nil, err
	}
	blockSize := aes.BlockSize
	if len(X)%blockSize != 0 {
		return nil, errors.New("input length must be a multiple of the AES block size")
	}

	Y := make([]byte, blockSize)
	Yaux := make([]byte, blockSize)
	for j := 0; j < len(X)/blockSize; j++ {
		Xj := X[j*blockSize : (j+1)*blockSize]
		for i := range Y {
			Y[i] ^= Xj[i]
		}
		block.Encrypt(Yaux, Y)
		copy(Y, Yaux)
	}
	return Y, nil
}
