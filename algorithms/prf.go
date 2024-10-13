// algorithms/prf.go
package algorithms

import (
	"crypto/aes"
	"crypto/cipher"
)

func PRF(X []byte, K []byte) []byte {
	block, err := aes.NewCipher(K)
	if err != nil {
		panic(err)
	}

	m := len(X) / aes.BlockSize
	Y := make([]byte, aes.BlockSize)
	for j := 0; j < m; j++ {
		Xj := X[j*aes.BlockSize : (j+1)*aes.BlockSize]
		for i := range Y {
			Y[i] ^= Xj[i]
		}
		block.Encrypt(Y, Y)
	}
	return Y
}
