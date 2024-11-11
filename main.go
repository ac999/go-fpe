// // main.go
package main

// import (
// 	"fmt"

// 	"github.com/ac999/go-fpe/algorithms"
// )

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"math"
	"math/big"
)

// Config for the alphabet and radix mapping.
type Alphabet struct {
	Chars  string
	Map    map[rune]int
	RevMap []rune
}

// NewAlphabet creates a new Alphabet based on the given characters.
func NewAlphabet(chars string) (*Alphabet, error) {
	radix := len(chars)
	if radix < 2 {
		return nil, errors.New("alphabet must contain at least two characters")
	}
	revMap := []rune(chars)
	charMap := make(map[rune]int)
	for i, char := range revMap {
		charMap[char] = i
	}
	return &Alphabet{
		Chars:  chars,
		Map:    charMap,
		RevMap: revMap,
	}, nil
}

// Encode converts a character string to numeral sequence based on the alphabet.
func (a *Alphabet) Encode(text string) ([]int, error) {
	result := make([]int, len(text))
	for i, char := range text {
		val, ok := a.Map[char]
		if !ok {
			return nil, fmt.Errorf("character %c not in alphabet", char)
		}
		result[i] = val
	}
	return result, nil
}

// Decode converts numeral sequence back to a character string based on the alphabet.
func (a *Alphabet) Decode(nums []int) (string, error) {
	var result bytes.Buffer
	for _, num := range nums {
		if num < 0 || num >= len(a.RevMap) {
			return "", fmt.Errorf("numeral %d out of alphabet bounds", num)
		}
		result.WriteRune(a.RevMap[num])
	}
	return result.String(), nil
}

// Helper Functions

// numRadix converts a numeral string X in a given radix to a number.
func numRadix(X []int, radix int) *big.Int {
	result := big.NewInt(0)
	r := big.NewInt(int64(radix))
	for _, x := range X {
		result.Mul(result, r)
		result.Add(result, big.NewInt(int64(x)))
	}
	return result
}

// strRadix converts a number x to a numeral string of a given length and radix.
func strRadix(x *big.Int, m int, radix int) ([]int, error) {
	result := make([]int, m)
	r := big.NewInt(int64(radix))
	for i := m - 1; i >= 0; i-- {
		mod := new(big.Int)
		x.DivMod(x, r, mod)
		result[i] = int(mod.Int64())
	}
	return result, nil
}

// PRF implements the pseudorandom function used in FF1 (AES-CBC).
func PRF(P, Q []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mac := make([]byte, aes.BlockSize)
	cbc := cipher.NewCBCEncrypter(block, mac)
	data := append(P, Q...)
	padding := bytes.Repeat([]byte{0}, aes.BlockSize-len(data)%aes.BlockSize)
	data = append(data, padding...)
	cbc.CryptBlocks(data, data)
	return data[len(data)-aes.BlockSize:], nil
}

// FF1Encrypt performs format-preserving encryption using FF1.
func FF1Encrypt(key, tweak []byte, X []int, radix, minLen, maxLen, maxTlen int) ([]int, error) {
	n := len(X)
	if n < minLen || n > maxLen || len(tweak) > maxTlen {
		return nil, errors.New("input length or tweak length out of bounds")
	}

	// Split X into two halves
	u := n / 2
	v := n - u
	A := X[:u]
	B := X[u:]

	// Initialize parameters for Feistel rounds
	b := int(math.Ceil(float64(v) * math.Log2(float64(radix)) / 8))
	if b < len(tweak)+1 { // Ensure b is not too small
		b = len(tweak) + 1
	}
	d := 4*(int(math.Ceil(float64(b)/4))) + 4

	P := append([]byte{0x01, 0x02, 0x01}, []byte{
		byte(radix >> 16), byte(radix >> 8), byte(radix),
		0x0a, byte(u % 256), byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n), byte(len(tweak) >> 8), byte(len(tweak)),
	}...)

	for i := 0; i < 10; i++ {
		// Round function
		Q := append(tweak, bytes.Repeat([]byte{0}, b-len(tweak)-1)...)
		Q = append(Q, byte(i))
		Bnum := numRadix(B, radix)
		Bbytes := Bnum.Bytes()
		if len(Bbytes) < b {
			Bbytes = append(bytes.Repeat([]byte{0}, b-len(Bbytes)), Bbytes...)
		}
		Q = append(Q, Bbytes...)

		R, err := PRF(P, Q, key)
		if err != nil {
			return nil, err
		}

		// Truncate or expand R as needed
		S := make([]byte, d)
		copy(S, R)
		for j := 1; j < int(math.Ceil(float64(d)/float64(aes.BlockSize))); j++ {
			block, _ := aes.NewCipher(key)
			rj := make([]byte, aes.BlockSize)
			copy(rj, R)
			for k := range rj {
				rj[k] ^= byte(j)
			}
			block.Encrypt(rj, rj)
			S = append(S, rj...)
		}
		S = S[:d]

		y := new(big.Int).SetBytes(S)

		// Update A and B
		m := u
		if i%2 != 0 {
			m = v
		}

		Aint := numRadix(A, radix)
		Aint.Add(Aint, y)
		Aint.Mod(Aint, new(big.Int).Exp(big.NewInt(int64(radix)), big.NewInt(int64(m)), nil))

		A, err = strRadix(Aint, m, radix)
		if err != nil {
			return nil, err
		}

		// Swap A and B
		if i < 9 {
			A, B = B, A
		}
	}

	return append(A, B...), nil
}

// FF1Decrypt performs format-preserving decryption using FF1.
func FF1Decrypt(key, tweak []byte, X []int, radix, minLen, maxLen, maxTlen int) ([]int, error) {
	n := len(X)
	if n < minLen || n > maxLen || len(tweak) > maxTlen {
		return nil, errors.New("input length or tweak length out of bounds")
	}

	// Split X into two halves
	u := n / 2
	v := n - u
	A := X[:u]
	B := X[u:]

	// Initialize parameters for Feistel rounds
	b := int(math.Ceil(float64(v) * math.Log2(float64(radix)) / 8))
	if b < len(tweak)+1 { // Ensure b is not too small
		b = len(tweak) + 1
	}
	d := 4*(int(math.Ceil(float64(b)/4))) + 4

	P := append([]byte{0x01, 0x02, 0x01}, []byte{
		byte(radix >> 16), byte(radix >> 8), byte(radix),
		0x0a, byte(u % 256), byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n), byte(len(tweak) >> 8), byte(len(tweak)),
	}...)

	// Reverse Feistel rounds from 9 to 0
	for i := 9; i >= 0; i-- {
		// Round function
		Q := append(tweak, bytes.Repeat([]byte{0}, b-len(tweak)-1)...)
		Q = append(Q, byte(i))
		Anum := numRadix(A, radix)
		Abytes := Anum.Bytes()
		if len(Abytes) < b {
			Abytes = append(bytes.Repeat([]byte{0}, b-len(Abytes)), Abytes...)
		}
		Q = append(Q, Abytes...)

		R, err := PRF(P, Q, key)
		if err != nil {
			return nil, err
		}

		// Truncate or expand R as needed
		S := make([]byte, d)
		copy(S, R)
		for j := 1; j < int(math.Ceil(float64(d)/float64(aes.BlockSize))); j++ {
			block, _ := aes.NewCipher(key)
			rj := make([]byte, aes.BlockSize)
			copy(rj, R)
			for k := range rj {
				rj[k] ^= byte(j)
			}
			block.Encrypt(rj, rj)
			S = append(S, rj...)
		}
		S = S[:d]

		y := new(big.Int).SetBytes(S)

		// Update A and B
		m := u
		if i%2 != 0 {
			m = v
		}

		Bint := numRadix(B, radix)
		Bint.Sub(Bint, y) // Use modular subtraction for decryption
		Bint.Mod(Bint, new(big.Int).Exp(big.NewInt(int64(radix)), big.NewInt(int64(m)), nil))

		B, err = strRadix(Bint, m, radix)
		if err != nil {
			return nil, err
		}

		// Swap A and B for the next round
		if i > 0 {
			A, B = B, A
		}
	}

	return append(A, B...), nil
}

func main() {
	// Define the alphabet and sample text
	alphabet, _ := NewAlphabet("abcdefghijklmnopqrstuvwxyz")
	originalText := "hello"

	// Encode the text to a numeral sequence using the defined alphabet
	X, _ := alphabet.Encode(originalText)
	fmt.Println("Original Numeral Sequence:", X)

	// Encryption setup
	key := []byte("examplekey123456") // 16-byte key for AES-128
	tweak := []byte("tweak")
	radix := 26 // Base-26 for lowercase English alphabet
	minLen, maxLen, maxTlen := 2, 10, 8

	// Encrypt the numeral sequence
	ciphertext, err := FF1Encrypt(key, tweak, X, radix, minLen, maxLen, maxTlen)
	if err != nil {
		fmt.Println("Encryption Error:", err)
		return
	}
	fmt.Println("Encrypted Numeral Sequence:", ciphertext)

	// Decrypt the encrypted numeral sequence
	decryptedNumeralSeq, err := FF1Decrypt(key, tweak, ciphertext, radix, minLen, maxLen, maxTlen)
	if err != nil {
		fmt.Println("Decryption Error:", err)
		return
	}
	fmt.Println("Decrypted Numeral Sequence:", decryptedNumeralSeq)

	// Decode the numeral sequence back to text
	decryptedText, _ := alphabet.Decode(decryptedNumeralSeq)
	fmt.Println("Decrypted Text:", decryptedText)

	// Verify that the decrypted text matches the original text
	if decryptedText == originalText {
		fmt.Println("Decryption successful, original text matches decrypted text.")
	} else {
		fmt.Println("Decryption failed, original text does not match decrypted text.")
	}
}
