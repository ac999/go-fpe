// helpers.go
package algorithms

import (
	"errors"
	"math/big"
)

// Helper functions

func floor(x, y uint64) uint64 {
	return x / y
}

func ceiling(x, y int) int {
	if x%y == 0 {
		return x / y
	}
	return x/y + 1
}

func mod(x, m int) int {
	return ((x % m) + m) % m
}

func representCharacters(s string, radix int) ([]int, error) {
	characters := make([]int, len(s))
	for i, ch := range s {
		charVal := int(ch - '0')
		if charVal >= radix {
			return nil, errors.New("character out of range for radix")
		}
		characters[i] = charVal
	}
	return characters, nil
}

func numStringToInt(s string, radix int) *big.Int {
	n := big.NewInt(0)
	for _, ch := range s {
		n.Mul(n, big.NewInt(int64(radix)))
		n.Add(n, big.NewInt(int64(ch-'0')))
	}
	return n
}

func intToNumString(x *big.Int, m int, radix int) string {
	num := x
	radixBig := big.NewInt(int64(radix))
	digits := make([]byte, m)
	for i := m - 1; i >= 0; i-- {
		mod := new(big.Int)
		num.DivMod(num, radixBig, mod)
		digits[i] = byte(mod.Int64() + '0')
	}
	return string(digits)
}

// // Power - Computes x^y for uint64, avoids overflow
// func Power(x, y uint64) uint64 {
// 	result := uint64(1)
// 	for i := uint64(0); i < y; i++ {
// 		result *= x
// 	}
// 	return result
// }

// // NumRadix - Converts a numeral string `X` into an integer based on the radix
// func NumRadix(X string, radix int) (uint64, error) {
// 	var x uint64
// 	for _, char := range X {
// 		var digit uint64
// 		if char >= '0' && char <= '9' {
// 			digit = uint64(char - '0')
// 		} else if char >= 'A' && char <= 'Z' {
// 			digit = uint64(char-'A') + 10
// 		} else {
// 			return 0, fmt.Errorf("invalid character '%c' for radix %d", char, radix)
// 		}

// 		if digit >= uint64(radix) {
// 			return 0, fmt.Errorf("invalid digit '%c' for radix %d", char, radix)
// 		}
// 		x = x*uint64(radix) + digit
// 	}
// 	return x, nil
// }

// // StrmRadix - Converts integer `x` into a numeral string of length `m` in a given radix
// func StrmRadix(x uint64, radix int, m int) (string, error) {
// 	if x >= Power(uint64(radix), uint64(m)) {
// 		return "", fmt.Errorf("x (%d) out of bounds for radix^m", x)
// 	}

// 	digits := make([]int, m)
// 	for i := 0; i < m; i++ {
// 		digits[m-1-i] = int(x % uint64(radix))
// 		x /= uint64(radix)
// 	}

// 	var result strings.Builder
// 	for _, digit := range digits {
// 		if digit < 10 {
// 			result.WriteByte(byte(digit + '0'))
// 		} else {
// 			result.WriteByte(byte(digit - 10 + 'A'))
// 		}
// 	}
// 	return result.String(), nil
// }

// // NumBits converts a binary string X, represented in bits, to an integer.
// func NumBits(X string) (uint64, error) {
// 	var x uint64 = 0
// 	for i := 0; i < len(X); i++ {
// 		// Check for valid binary characters.
// 		if X[i] != '0' && X[i] != '1' {
// 			return 0, fmt.Errorf("invalid character '%c' in binary string", X[i])
// 		}

// 		// Convert '0' or '1' to integer 0 or 1 by subtracting '0'.
// 		bit := X[i] - '0'
// 		x = 2*x + uint64(bit)
// 	}

// 	return x, nil
// }

// func uint64ToNBytes(x uint64, n int) []byte {
// 	result := make([]byte, n)
// 	for i := n - 1; i >= 0; i-- {
// 		result[i] = byte(x % 10)
// 		x /= 10
// 	}
// 	return []byte(result)
// }
