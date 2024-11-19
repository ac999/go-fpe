// component.go
package algorithms

import (
	"crypto/aes"
	"fmt"
	"math/big"
)

// NUM - Bit string to uint64 - could be faster by using bit shifting
func NUM(X []byte) uint64 {
	x := uint64(0)
	for _, b := range X {
		x = (x << 8) | uint64(b) // Shift existing bits left and append the current byte
	}
	return x
}

// NUM - Bit string to big.Int (instead of uint64) to avoid overflow
func BigNUM(X []byte) *big.Int {
	// Start with big.Int initialized to 0
	result := big.NewInt(0)

	for _, b := range X {
		// Shift the existing value (multiply by 256) and add the current byte (as an integer)
		result.Lsh(result, 8)                    // Equivalent to result = result * 256 (shift left by 8 bits)
		result.Add(result, big.NewInt(int64(b))) // Add the byte value (converted to int64)
	}

	return result
}

// NUMradix - Numeral string to uint64
func NUMradix(X []byte, radix uint64) uint64 {
	x := uint64(0)
	for _, i := range X {
		x = x*radix + uint64(i)
	}
	return x
}

// BigNUMradix - Numeral string to big.Int, with the specified radix
func BigNUMradix(X []byte, radix uint64) *big.Int {
	x := big.NewInt(0) // Start with zero
	for _, i := range X {
		// Multiply by radix and add the current byte value
		x.Mul(x, big.NewInt(int64(radix)))
		x.Add(x, big.NewInt(int64(i)))
	}
	return x
}

// STRmRadix - Representation of a uint64 as a string of m numerals in base
func STRmRadix(x uint64, radix uint64, m int64) []byte {
	X := make([]byte, m)
	for i := int64(0); i < m; i++ {
		X[m-1-i] = byte(Mod(x, radix))
		x = x / radix
	}
	return X
}

// BigSTRmRadix - Representation of a big.Int as a string of m numerals in base `radix`
func BigSTRmRadix(x *big.Int, radix uint64, m int64) []byte {
	xCopy := new(big.Int).Set(x)  // Make a copy of x
	X := make([]byte, m)          // Create a byte slice to hold the result
	r := big.NewInt(int64(radix)) // radix as big.Int

	for i := int64(0); i < m; i++ {
		mod := new(big.Int)
		mod.Mod(xCopy, r)             // mod = xCopy % radix
		X[m-1-i] = byte(mod.Uint64()) // Store the modulus as byte in the slice
		xCopy.Div(xCopy, r)           // xCopy = xCopy / radix
	}

	return X
}

func PRF(K []byte, X []byte) ([]byte, error) {
	block, err := aes.NewCipher(K)
	if err != nil {
		return nil, err
	}

	blockSize := aes.BlockSize // 16 bytes = 128 bits
	m := len(X) / blockSize

	if len(X)%blockSize != 0 {
		return nil, fmt.Errorf("input length must be a multiple of %d bytes", blockSize)
	}

	// Step 1: Initialize Y[0] as 0^128 (16 zero bytes)
	Y := make([]byte, blockSize)

	// Step 2: Process each block
	for j := 0; j < m; j++ {
		start := j * blockSize
		end := start + blockSize
		blockX := X[start:end]

		// XOR Y[j-1] with X[j]
		xorValue, err := XORBytes(Y, blockX)
		if err != nil {
			return nil, fmt.Errorf("XORBytes failed: %w", err)
		}

		// Step 3: Encrypt the result with AES
		block.Encrypt(Y, xorValue)
	}

	// Step 4: Return the last block
	return Y, nil
}
