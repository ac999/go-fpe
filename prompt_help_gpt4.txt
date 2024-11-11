I tried to implement FF1 encryption and decryption functions in GO. Please find my implementation below:

// algorithms/helpers.go:
package algorithms

import "encoding/binary"

func UintToBytes(num uint64) []byte {
	var bytes = make([]byte, 8)
	binary.PutUvarint(bytes, num)
	return bytes
}

func power(x uint64, y uint64) uint64 {
	result := uint64(1)
	for i := uint64(0); i < y; i++ {
		result *= x
	}
	return result
}

func intToNBytes(x int, n int) []byte {

	result := make([]byte, n)
	for i := n - 1; i >= 0; i-- {
		result[i] = byte(x % 10)
		x /= 10
	}
	return []byte(result)
}

func uint64ToNBytes(x uint64, n int) []byte {
	result := make([]byte, n)
	for i := n - 1; i >= 0; i-- {
		result[i] = byte(x % 10)
		x /= 10
	}
	return []byte(result)
}

// algorithms/num.go
package algorithms

import (
	"fmt"
	"strconv"
)

// NumRadix(X) function implements Algorithm 1: NUM_radix(X) from NIST Special Publication 800-38G.
// Converts a numeral string to a number based on the given radix.
func NumRadix(X string, radix int) (uint64, error) {
	var x uint64 = 0

	// Validate that the radix is within a reasonable range
	if radix < 2 || radix > 36 {
		return 0, fmt.Errorf("invalid radix: %d. Must be between 2 and 36", radix)
	}

	for i := 0; i < len(X); i++ {
		// Get the value of the current character
		digit, err := strconv.ParseUint(string(X[i]), radix, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid character '%c' for radix %d", X[i], radix)
		}

		x = x*uint64(radix) + digit
	}

	return x, nil
}

// NumBits(X) function implements Algorithm 2: NUM(X) from NIST Special Publication 800-38G.
// Converts a byte string represented in bits to an integer.
func NumBits(X string) (uint64, error) {
	var x uint64 = 0
	for i := 0; i < len(X); i++ {
		// if X[i] != '0' && X[i] != '1' {
		// 	return 0, fmt.Errorf("invalid character '%c' in binary string", X[i])
		// }

		bit := X[i] // - '0' // Convert character '0' or '1' to integer 0 or 1
		x = 2*x + uint64(bit)
	}

	return x, nil
}

// algorithms/strm.go
package algorithms

import (
	"fmt"
	"strings"
)

// StrmRadix(X) function implements Algorithm 3: STR^m_radix(X) from NIST Special Publication 800-38G.
// Converts an integer to a numeral string of a given length and radix.
func StrmRadix(x uint64, radix int, m int) (string, error) {
	// Validate that the input x is within the valid range
	maxValue := power(uint64(radix), uint64(m))

	if x >= maxValue {
		return "", fmt.Errorf("invalid input: x (%d) must be less than radix^m (%d)", x, maxValue)
	}

	X := make([]int, m)
	for i := 0; i < m; i++ {
		X[m-1-i] = int(x % uint64(radix)) // X[m+1–i] = x mod radix
		x = x / uint64(radix)             // x = floor(x / radix)
	}

	// Efficiently build the result string
	var result strings.Builder
	for _, digit := range X {
		if digit < 10 {
			result.WriteByte(byte(digit + '0')) // 0-9 are represented as '0'-'9'
		} else {
			result.WriteByte(byte(digit - 10 + 'A')) // 10-35 are represented as 'A'-'Z'
		}
	}
	return result.String(), nil
}

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
	Yaux := make([]byte, blockSize)

	// Iterate over each block Xj
	for j := 0; j < m; j++ {
		Xj := X[j*blockSize : (j+1)*blockSize] // Get block Xj

		// XOR Yj-1 with Xj
		for i := range Y {
			Y[i] ^= Xj[i]
		}

		// Encrypt the result to produce Yj
		block.Encrypt(Yaux, Y)
		copy(Y, Yaux)
	}

	// Return Ym
	return Y, nil
}

// algorithms/ff1.go
package algorithms

import (
	"crypto/aes"
	"errors"
	"fmt"
	"log"
	"math"
)

// FF1Encrypt performs format-preserving encryption using FF1
func FF1Encrypt(K, T []byte, X string, radix, minlen, maxlen, maxTlen int) (string, error) {
	// Validate input length
	n := len(X)
	if n < minlen || n > maxlen {
		return "", errors.New("input length out of range")
	}

	// Validate the key length as in crypto/aes/cipher.go
	switch len(K) {
	default:
		return "", errors.New("invalid AES key size; must be 16, 24, or 32 bytes")
	case 16, 24, 32:
		break
	}

	// Steps 1 & 2: Split the numeral string X into A and B
	u := n / 2
	v := n - u
	A := X[:u]
	B := X[u:]

	t := len(T)

	// Steps 3 & 4: Calculate the byte length b and block length d
	b := int(math.Ceil(float64(v) * math.Log2(float64(radix)) / 8))
	if b <= 0 {
		return "", errors.New("invalid block length b")
	}
	d := 4*int(math.Ceil(float64(b)/4)) + 4

	// Step 5: Construct P = [1, 2, 1, radix (3 bytes), 10, u mod 256, n (4 bytes), t (4 bytes)]
	radix3bytes := intToNBytes(radix, 3)
	nBytes := intToNBytes(n, 4)
	tBytes := intToNBytes(t, 4)
	P := []byte{1, 2, 1}
	P = append(P, radix3bytes...)
	P = append(P, 10, byte(u%256))
	P = append(P, nBytes...)
	P = append(P, tBytes...)

	// Step 6: Perform the Feistel rounds
	for i := 0; i < 10; i++ {
		// Step 6.i: Construct Q safely
		padding := 16 + (0-t-b-1)%16
		Q := append(T, make([]byte, padding)...)
		Q = append(Q, byte(i)) // append round number i

		// Convert B to a numeric representation
		BNum, err := NumRadix(B, radix)
		if err != nil {
			return "", err
		}
		bNumBytes := uint64ToNBytes(BNum, b)
		Q = append(Q, bNumBytes...)

		//  Check size of Q to be multiple of 16 before calling PRF(P, Q)
		switch len(Q) % 16 {
		default:
			error_str := fmt.Sprint("invalid Q size: ", len(Q), "; must be 16, 24, or 32 bytes")
			return "", errors.New(error_str)
		case 0:
		}

		// Step 6.ii: Compute R = PRF(P, Q)
		R, err := PRF(append(P, Q...), K)
		if err != nil {
			log.Println("P is: ", P, " of length: ", len(P))
			log.Println("Q is: ", Q, " of length: ", len(Q))
			return "", err
		}

		// Step 6.iii: Let S be the first d bytes of CIPHK(R ⊕ [j]16)
		if d <= 0 {
			return "", errors.New("invalid block length d")
		}

		block, err := aes.NewCipher(K)
		if err != nil {
			return "", err
		}

		blockSize := aes.BlockSize // 16 bytes for AES

		dpe16 := int(math.Ceil(float64(d) / float64(blockSize)))
		S := make([]byte, len(R))
		copy(S, R)
		for j := 1; j < dpe16; j++ {
			j16Bytes := intToNBytes(j, 16)
			res := make([]byte, len(R))
			// R xor j
			for k, r := range R {
				res[k] = r ^ j16Bytes[k]
			}
			aux := make([]byte, len(R))
			block.Encrypt(aux, res)

			S = append(S, aux...)
		}

		// for j := 0; j < int(math.Ceil(float64(d)/float64(blockSize))); j++ {
		// 	Rj := make([]byte, aes.BlockSize)
		// 	copy(Rj, append(R, byte(j)))
		// 	encrypted := make([]byte, blockSize)
		// 	block.Encrypt(encrypted, Rj)
		// 	S = append(S, encrypted...)
		// }
		// S = S[:d] // truncate to d bytes

		// Step 6.iv: Let y = NUM(S)
		y, err := NumBits(string(S))
		if err != nil {
			return "", err
		}

		// Step 6.v: Determine m
		m := u
		if i%2 != 0 {
			m = v
		}

		// Step 6.vi: Compute c = (NUMradix(A) + y) mod radix^m
		aNum, err := NumRadix(A, radix)
		if err != nil {
			return "", err
		}
		c := (aNum + y) % power(uint64(radix), uint64(m))

		// Step 6.vii: Let C = STRmradix(c)
		C, err := StrmRadix(c, radix, m)
		if err != nil {
			return "", err
		}

		// Step 6.viii and 6.ix: Swap A and B for next round
		A = B
		B = C
	}

	// Step 7: Return A || B as the encrypted result
	return A + B, nil
}

// FF1 Decryption Algorithm
func FF1Decrypt(K, T []byte, X string, radix, minlen, maxlen, maxTlen int) (string, error) {
	// Validate input length
	n := len(X)
	if n < minlen || n > maxlen {
		return "", errors.New("input length out of range")
	}

	// Validate the key length as in crypto/aes/cipher.go
	switch len(K) {
	default:
		return "", errors.New("invalid AES key size; must be 16, 24, or 32 bytes")
	case 16, 24, 32:
		break
	}

	// Steps 1 & 2: Split the numeral string X into A and B
	u := n / 2
	v := n - u
	A := X[:u]
	B := X[u:]

	t := len(T)

	// Steps 3 & 4: Calculate the byte length b and block length d
	b := int(math.Ceil(float64(v) * math.Log2(float64(radix)) / 8))
	if b <= 0 {
		return "", errors.New("invalid block length b")
	}
	d := 4*int(math.Ceil(float64(b)/4)) + 4

	// Step 5: Construct P = [1, 2, 1, radix (3 bytes), 10, u mod 256, n (4 bytes), t (4 bytes)]
	radix3bytes := intToNBytes(radix, 3)
	nBytes := intToNBytes(n, 4)
	tBytes := intToNBytes(t, 4)
	P := []byte{1, 2, 1}
	P = append(P, radix3bytes...)
	P = append(P, 10, byte(u%256))
	P = append(P, nBytes...)
	P = append(P, tBytes...)

	// Step 6: Perform the Feistel rounds
	for i := 9; i >= 0; i-- {
		// Step 6.i: Construct Q safely
		padding := 16 + (0-t-b-1)%16
		Q := append(T, make([]byte, padding)...)
		Q = append(Q, byte(i)) // append round number i

		// Convert B to a numeric representation
		ANum, err := NumRadix(A, radix)
		if err != nil {
			return "", err
		}
		bNumBytes := uint64ToNBytes(ANum, b)
		Q = append(Q, bNumBytes...)

		//  Check size of Q to be multiple of 16 before calling PRF(P, Q)
		switch len(Q) % 16 {
		default:
			error_str := fmt.Sprint("invalid Q size: ", len(Q), "; must be 16, 24, or 32 bytes")
			return "", errors.New(error_str)
		case 0:
		}

		// Step 6.ii: Compute R = PRF(P, Q)
		R, err := PRF(append(P, Q...), K)
		if err != nil {
			log.Println("P is: ", P, " of length: ", len(P))
			log.Println("Q is: ", Q, " of length: ", len(Q))
			return "", err
		}

		// Step 6.iii: Let S be the first d bytes of CIPHK(R ⊕ [j]16)
		if d <= 0 {
			return "", errors.New("invalid block length d")
		}

		block, err := aes.NewCipher(K)
		if err != nil {
			return "", err
		}

		blockSize := aes.BlockSize // 16 bytes for AES

		dpe16 := int(math.Ceil(float64(d) / float64(blockSize)))
		S := make([]byte, len(R))
		copy(S, R)
		for j := 1; j < dpe16; j++ {
			j16Bytes := intToNBytes(j, 16)
			res := make([]byte, len(R))
			// R xor j
			for k, r := range R {
				res[k] = r ^ j16Bytes[k]
			}
			aux := make([]byte, len(R))
			block.Decrypt(aux, res)

			S = append(S, aux...)
		}

		// for j := 0; j < int(math.Ceil(float64(d)/float64(blockSize))); j++ {
		// 	Rj := make([]byte, aes.BlockSize)
		// 	copy(Rj, append(R, byte(j)))
		// 	encrypted := make([]byte, blockSize)
		// 	block.Encrypt(encrypted, Rj)
		// 	S = append(S, encrypted...)
		// }
		// S = S[:d] // truncate to d bytes

		// Step 6.iv: Let y = NUM(S)
		y, err := NumBits(string(S))
		if err != nil {
			return "", err
		}

		// Step 6.v: Determine m
		m := u
		if i%2 != 0 {
			m = v
		}

		// Step 6.vi: Compute c = (NUMradix(A) + y) mod radix^m
		bNum, err := NumRadix(B, radix)
		if err != nil {
			return "", err
		}
		c := (bNum - y) % power(uint64(radix), uint64(m))

		// Step 6.vii: Let C = STRmradix(c)
		C, err := StrmRadix(c, radix, m)
		if err != nil {
			return "", err
		}

		// Step 6.viii and 6.ix: Swap A and B for next round
		B = A
		A = C
	}

	// Step 7: Return A || B as the encrypted result
	return A + B, nil
}

// tests/algorithms_test.go
package tests

import (
	"crypto/aes"
	"testing"
	"go-fpe/algorithms"
)
func TestFF1Encrypt(t *testing.T) {
	key := []byte("examplekey123456") // 16-byte AES key
	tweak := []byte{0, 1, 2, 3, 4, 5} // example tweak
	input := "0123456789"             // numeral string to encrypt
	radix := 10
	minlen := 6
	maxlen := 12
	maxTlen := 8

	encrypted, err := algorithms.FF1Encrypt(key, tweak, input, radix, minlen, maxlen, maxTlen)
	if err != nil {
		t.Fatalf("FF1Encrypt failed: %v", err)
	}
	t.Logf("Encrypted Output: %s", encrypted)
	if len(encrypted) != len(input) {
		t.Fatalf("Expected length: %v but got: %v", len(input), len(encrypted))
	}

	decrypted, err := algorithms.FF1Decrypt(key, tweak, encrypted, radix, minlen, maxlen, maxTlen)
	if err != nil {
		t.Fatalf("FF1Decrypt failed: %v", err)
	}
	t.Logf("Decrypted Output: %s", decrypted)

	if input != encrypted {
		t.Fatalf("Expected decrypted text: %v but got: %v", input, decrypted)
	}
}

Correct my implementation, taking into consideration all sections from the NIST recommendation for format-preserving encryption (including preliminaries, etc.)





------------------------

https://chatgpt.com/c/67325e4e-8af8-800f-998f-5fb9969b3d6d

For the following input to the encrypt function:
```
key := []byte("examplekey123456") // 16-byte AES key 
	tweak := []byte{0, 1, 2, 3, 4, 5} // example tweak
	input := "0123456789"             // numeral string to encrypt
	radix := 10
	minlen := 6
	maxlen := 12
	maxTlen := 8

	encrypted, err := algorithms.FF1Encrypt(key, tweak, input, radix, minlen, maxlen, maxTlen)
``` i get the following error:
FF1Encrypt failed: invalid character '�' in binary string.

Shouldn't the code also include the following implementation of the preliminaries from NIST?
```
The data inputs and outputs for FF1 and FF3 are sequences of numbers that can represent both numeric and non-numeric data, as discussed below.
A finite set of two or more symbols is called an alphabet. The symbols in an alphabet are called the characters of the alphabet. The number of characters in an alphabet is called the base, denoted by radix; thus, radix ≥ 2.
A character string is a finite sequence of characters from an alphabet; individual characters may repeat in the string. In this publication, character strings (and bit strings) are presented in the Courier Newfont.
Thus, for the alphabet of lower-case English letters,
{a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p, q, r, s, t, u, v, w, x, y, z},
hello and cannot are character strings, but Hello and can’t are not, because the symbols “H” and “ ’ ” are not in the alphabet.
SSNs or CCNs can be regarded as character strings in the alphabet of base ten numerals, namely, {0, 1, 2, 3, 4, 5, 6, 7, 8, 9}. The notion of numerals is generalized to any given base as follows: the set of base radix numerals is
{0, 1, …, radix-1}.
The data inputs and outputs to the FF1 and FF3 encryption and decryption functions must be finite sequences of numerals, i.e., numeral strings. If the data to be encrypted is formatted in an alphabet that is not already the set of base radix numerals, then each character must be represented by a distinct numeral in order to apply FF1 or FF3.
For example, the natural representation of lower-case English letters with base 26 numerals is
a→0, b→1, c→2, … x→23, y→24, z→25.
6
NIST SP 800-38G METHODS FOR FORMAT-PRESERVING ENCRYPTION
The character string hello would then be represented by the numeral string 7 4 11 11 14. Other representations are possible.
The choice and implementation of a one-to-one correspondence between a given alphabet and the set of base radix numerals that represents the alphabet is outside the scope of this publication.
In this publication, individual numerals are themselves represented in base ten. In order to display numeral sequences unambiguously when the base is greater than ten, a delimiter between the numerals is required, such as a space (as in the base 26 example above) or a comma.
FF1 and FF3 use different conventions for interpreting numeral strings as numbers. For FF1, numbers are represented by strings of numerals with decreasing order of significance; for FF3, numbers are represented by strings of numerals in the reverse order, i.e., with increasing order of significance. Algorithms for the functions that convert numeral strings to numbers and vice versa are given in Sec. 4.6.
```
Also are all the examples from section 4.2 Examples of Basic Operations and Functions from NIST validated for this Go-lang implementation?

Make sure all this are taken into consideration and correct the error.