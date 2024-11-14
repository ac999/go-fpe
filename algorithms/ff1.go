package algorithms

// import (
// 	"crypto/aes"
// 	"crypto/cipher"
// 	"math/big"
// )

// // Encrypt a single block using AES
// func encryptBlockAES(block cipher.Block, input []byte) []byte {
// 	output := make([]byte, block.BlockSize())
// 	block.Encrypt(output, input)
// 	return output
// }

// // FF1 Encryption Core with Feistel rounds
// func FF1Encrypt(key, tweak []byte, plaintext string, radix int) (string, error) {
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return "", err
// 	}

// 	n := len(plaintext)
// 	u := n / 2
// 	v := n - u

// 	A, B := plaintext[:u], plaintext[u:]
// 	charsA, err := representCharacters(A, radix)
// 	charsB, err := representCharacters(B, radix)

// 	for i := 0; i < 10; i++ {
// 		// Construct Q with tweak and round index
// 		Q := make([]byte, len(tweak)+1+len(charsA)+len(charsB))
// 		copy(Q, tweak)
// 		Q[len(tweak)] = byte(i)

// 		// Calculate S and y using AES encryption on Q
// 		y := new(big.Int)
// 		R := encryptBlockAES(block, Q)
// 		y.SetBytes(R[:len(R)/2])

// 		m := len(charsB)
// 		c := mod(int(y.Int64()), m)
// 		C := intToNumString(big.NewInt(int64(c)), m, radix)

// 		charsA, charsB = charsB, C
// 	}
// 	return A + B, nil
// }

// func FF1Decrypt(key, tweak []byte, ciphertext string, radix int) (string, error) {
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return "", err
// 	}

// 	n := len(ciphertext)
// 	u := n / 2
// 	v := n - u

// 	A, B := ciphertext[:u], ciphertext[u:]
// 	charsA, err := representCharacters(A, radix)
// 	charsB, err := representCharacters(B, radix)

// 	for i := 9; i >= 0; i-- {
// 		Q := make([]byte, len(tweak)+1+len(charsA)+len(charsB))
// 		copy(Q, tweak)
// 		Q[len(tweak)] = byte(i)

// 		y := new(big.Int)
// 		R := encryptBlockAES(block, Q)
// 		y.SetBytes(R[:len(R)/2])

// 		m := len(charsA)
// 		c := mod(int(y.Int64()), m)
// 		C := intToNumString(big.NewInt(int64(c)), m, radix)

// 		charsA, charsB = charsB, C
// 	}
// 	return A + B, nil
// }
