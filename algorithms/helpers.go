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
	// Calculate the quotient and remainder
	quot := new(big.Int).Div(x, y)
	mod := new(big.Int).Mod(x, y)

	fmt.Printf("quot = %v, mod = %v\n", quot, mod)

	// If the remainder is non-zero, check if we need to round up
	if mod.Sign() != 0 {
		fmt.Printf("mod sign != 0\n")
		// Round up if both numbers have the same sign (positive/positive or negative/negative)
		if (x.Sign() > 0 && y.Sign() > 0) || (x.Sign() < 0 && y.Sign() < 0) {
			fmt.Printf("numbers same sign\n")
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

func ModBigInt(x, m *big.Int) *big.Int {
	// Create a new big.Int to hold the result
	result := new(big.Int).Mod(x, m)
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

func ByteLen(x []byte) uint64 {
	return uint64(len(x) / 8)
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

func Pad(x, m uint64) []byte {
	result := make([]byte, m)
	result[m-1] = byte(x)
	return result
}

func BytesToUint64Array(data []byte) []uint64 {
	result := make([]uint64, len(data))
	for i, v := range data {
		result[i] = uint64(v)
	}
	return result
}
