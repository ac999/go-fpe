// helpers.go
package algorithms

import (
	"errors"
	"fmt"
	"math/big"
)

// Helper functions

// Representation of Character Strings

// StringToNumeralSlice converts a character string to a slice of numerals (`[]uint64`) based on the specified alphabet.
// The radix is derived from the length of the alphabet. Returns an error if the input contains characters not in the alphabet.
func StringToNumeralSlice(input, alphabet string) ([]byte, error) {
	// Create a map for character to numeral values based on the alphabet
	charToNum := make(map[rune]byte)
	for i, char := range alphabet {
		charToNum[char] = byte(i)
	}

	// Convert each character to its numeral representation
	numerals := make([]byte, len(input))
	for i, ch := range input {
		num, exists := charToNum[ch]
		if !exists {
			return nil, errors.New("input contains characters not in the specified alphabet")
		}
		numerals[i] = num
	}
	return numerals, nil
}

// Basic Operations and Functions

func FloorDiv(x, y uint64) uint64 {
	return x / y
}

// BigFloorDiv returns the floor of x / y for big.Int.
func BigFloorDiv(x, y *big.Int) *big.Int {
	quot := new(big.Int).Div(x, y)
	remainder := new(big.Int).Mod(x, y)

	// Adjust quotient for floor division correctly
	if remainder.Sign() != 0 && (x.Sign() != y.Sign()) {
		if x.Cmp(y) > 0 {
			quot.Sub(quot, big.NewInt(1))
		}
	}

	return quot
}

func CeilingDiv(x, y uint64) uint64 {
	if x%y == 0 {
		return x / y
	}
	return x/y + 1
}

func BigCeilingDiv(x, y *big.Int) *big.Int {
	if y.Sign() == 0 {
		panic("division by zero")
	}

	// Calculate quotient and remainder using Euclidean division
	quot := new(big.Int).Div(x, y)
	mod := new(big.Int).Mod(x, y)

	// Adjust for ceiling behavior
	if mod.Sign() != 0 {
		// Rounding up when:
		// - Both x and y are positive, or
		// - x is negative and y is positive (rounding towards zero)
		if (x.Sign() > 0 && y.Sign() > 0) || (x.Sign() < 0 && y.Sign() > 0) {
			quot.Add(quot, big.NewInt(1))
		}
	}

	return quot
}

func Mod(x, m uint64) uint64 {
	return x - m*(x/m)
}

func ModInt(x int64, m int64) int64 {
	remainder := x % m
	if remainder < 0 {
		remainder += m
	}
	return remainder
}

// BigMod calculates x % m using Euclidean modulus with *big.Int
func BigMod(x, m *big.Int) *big.Int {
	// Perform the modulus operation
	result := new(big.Int)
	result.Mod(x, m)

	// Ensure the result is in the range [0, |m|) (Euclidean modulus)
	if result.Sign() < 0 {
		result.Add(result, m)
	}

	return result
}

// Power - Computes x^y for uint64
func Power(x, y uint64) uint64 {
	result := uint64(1)
	for i := uint64(0); i < y; i++ {
		result *= x
	}
	return result
}

func BigPower(x, y *big.Int) *big.Int {
	result := big.NewInt(1) // Start with 1 as the result
	for i := big.NewInt(0); i.Cmp(y) < 0; i.Add(i, big.NewInt(1)) {
		result.Mul(result, x) // Multiply result by x each iteration
	}
	return result
}

func ByteLen(x []byte) uint64 {
	return uint64(len(x) / 8)
}

func BigByteLen(x []byte) *big.Int {
	// Get the length of the byte slice
	byteLength := len(x)

	// Convert the byte length to a big.Int
	bigLen := new(big.Int).SetUint64(uint64(byteLength))

	// Divide by 8 (since we're converting byte length to 8-bit units)
	bigLen.Div(bigLen, big.NewInt(8))

	return bigLen
}

// BreakInBlocks splits a byte slice into blocks of a specified size.
func BreakInBlocks(X []byte, blockSize int) ([][]byte, error) {
	if len(X)%blockSize != 0 {
		return nil, fmt.Errorf("the length of the byte slice must be a multiple of the block size")
	}
	var blocks [][]byte
	for i := 0; i < len(X); i += blockSize {
		blocks = append(blocks, X[i:i+blockSize])
	}
	return blocks, nil
}

// Xor on byte slices
func XORBytes(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("byte slices must be of the same length")
	}

	result := make([]byte, len(a))
	for i := range a {
		result[i] = a[i] ^ b[i]
	}
	return result, nil
}
