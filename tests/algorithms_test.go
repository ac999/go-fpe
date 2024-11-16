// tests/algorithms_test.go
package tests

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ac999/go-fpe/algorithms"
)

var alphabets = map[string]string{
	"base10": "0123456789",
	"base26": "abcdefghijklmnopqrstuvwxyz",
	"base36": "0123456789abcdefghijklmnopqrstuvwxyz",
}

func TestStringToNumeralSlice(t *testing.T) {
	tests := []struct {
		input     string
		alphabet  string
		expected  []uint64
		shouldErr bool
	}{
		{"hello", alphabets["base26"], []uint64{7, 4, 11, 11, 14}, false},
		{"01234", alphabets["base10"], []uint64{0, 1, 2, 3, 4}, false},
		{"hello1", alphabets["base26"], nil, true}, // Invalid character '1' for base26
	}

	for _, test := range tests {
		result, err := algorithms.StringToNumeralSlice(test.input, test.alphabet)
		if test.shouldErr && err == nil {
			t.Errorf("expected error but got none for input %s", test.input)
		} else if !test.shouldErr {
			for i, v := range result {
				if v != test.expected[i] {
					t.Errorf("expected %v at index %d, got %v", test.expected[i], i, v)
				}
			}
		}
	}
}

func TestNumeralSliceToInt(t *testing.T) {
	numeralSlice := []uint64{7, 4, 11, 11, 14}
	radix := uint64(26)
	expected := uint64(7*26*26*26*26 + 4*26*26*26 + 11*26*26 + 11*26 + 14)

	result := algorithms.NumeralSliceToInt(numeralSlice, radix)
	if *result != expected {
		t.Errorf("expected %d, got %d", expected, *result)
	}
}

func TestCeilingDiv(t *testing.T) {
	result := algorithms.CeilingDiv(7, 5)
	if result != 2 {
		t.Errorf("expected %d, got %d", 2, result)
	}
	result = algorithms.CeilingDiv(3, 3)
	if result != 1 {
		t.Errorf("expected %d, got %d", 1, result)
	}
}

func TestMod(t *testing.T) {
	result := algorithms.Mod(13, 7)
	if result != 6 {
		t.Errorf("expected %d, got %d", 6, result)
	}
}

func TestByteLen(t *testing.T) {
	x := []byte{1, 0, 1, 1, 1, 0, 0, 1, 1, 0, 1, 0, 1, 1, 0, 0}
	result := algorithms.ByteLen(x)
	if result != 2 {
		t.Errorf("expected %d, got %d", 2, result)
	}
}

func TestNUM(t *testing.T) {
	X := []byte{1, 0, 0, 0, 0, 0, 0, 0}
	result := algorithms.NUM(X)
	if result != 128 {
		t.Errorf("expected %d, got %d", 128, result)
	}
}

func TestNUMradix(t *testing.T) {
	X := []byte{0, 0, 0, 1, 1, 0, 1, 0}
	result := algorithms.NUMradix(X, 5)
	if result != 755 {
		t.Errorf("expected %d, got %d", 755, result)
	}
}

func TestBreakInBlocks(t *testing.T) {
	x := []byte{0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	blockSize := 128
	expectedBlocks := [][]byte{
		{0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}

	blocks, _ := algorithms.BreakInBlocks(x, blockSize)
	if !reflect.DeepEqual(blocks, expectedBlocks) {
		t.Errorf("Expected %v, but got %v", expectedBlocks, blocks)
	}

}

func TestXORBytes(t *testing.T) {
	tests := []struct {
		name     string
		a        []byte
		b        []byte
		expected []byte
		wantErr  bool
	}{
		{
			name:     "Basic XOR",
			a:        []byte{0b10101010, 0b11001100},
			b:        []byte{0b01010101, 0b00110011},
			expected: []byte{0b11111111, 0b11111111},
			wantErr:  false,
		},
		{
			name:     "All zeroes",
			a:        []byte{0x00, 0x00},
			b:        []byte{0x00, 0x00},
			expected: []byte{0x00, 0x00},
			wantErr:  false,
		},
		{
			name:     "All ones",
			a:        []byte{0xFF, 0xFF},
			b:        []byte{0xFF, 0xFF},
			expected: []byte{0x00, 0x00},
			wantErr:  false,
		},
		{
			name:     "Mismatched Lengths",
			a:        []byte{0x01, 0x02},
			b:        []byte{0x01},
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "Single Byte",
			a:        []byte{0b10101010},
			b:        []byte{0b01010101},
			expected: []byte{0b11111111},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := algorithms.XORBytes(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("XORBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("XORBytes() = %08b, expected %08b", result, tt.expected)
			}
		})
	}
}

func TestSTRmRadix(t *testing.T) {
	expected1 := []byte{0, 3, 10, 7}
	expected2 := []byte{0, 0, 0, 0, 0, 0, 0, 1}

	result1 := algorithms.STRmRadix(559, 12, 4)
	result2 := algorithms.STRmRadix(1, 2, 8)
	if len(result1) != len(expected1) {
		t.Errorf("STRmRadix() length = %v, expected %v", len(result1), len(expected1))
	} else {
		for i := range len(result1) {
			if result1[i] != expected1[i] {
				fmt.Printf("expected result %d, got %d", expected1, result1)
				t.Errorf("STRmRadix() result at index %v = %v, expected %v", i, result1[i], expected1[i])
			}
		}
	}
	if len(result2) != len(expected2) {
		t.Errorf("STRmRadix() length = %v, expected %v", len(result2), len(expected2))
	} else {
		for i := range len(result1) {
			if result2[i] != expected2[i] {
				fmt.Printf("expected result %d, got %d", expected2, result2)
				t.Errorf("STRmRadix() result at index %v = %v, expected %v", i, result2[i], expected2[i])
			}
		}
	}
}

// func TestPRF(t *testing.T) {
// 	K := []byte{0x2b, 0x7e, 0x15, 0x16, 0x28, 0xae, 0xd2, 0xa6, 0xab, 0xf7, 0x15, 0x88, 0x09, 0xcf, 0x4f, 0x3c}
// 	X := []uint64{1, 2, 1, 0, 0, 10, 10, 5, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 221, 213}

// 	var Xbitstring []byte

// 	for i := range len(X) {
// 		Xbitstring = append(Xbitstring, algorithms.STRmRadix(X[i], 2, 8)...)
// 	}

// 	expected := []byte{195, 184, 41, 161, 232, 100, 43, 120, 204, 41, 148, 123, 59, 147, 219, 99}
// 	expectedLen := len(expected)
// 	result, err := algorithms.PRF(K, Xbitstring)
// 	if err != nil {
// 		t.Errorf("PRF() error = %v", err)
// 		return
// 	}

// 	if len(result) != expectedLen {
// 		t.Errorf("expected length %d, got %d", expectedLen, len(result))
// 	} else {
// 		for i := range expectedLen {
// 			if result[i] != expected[i] {
// 				t.Errorf("expected result at index %d: %d, got %d", i, expected[i], result[i])
// 			}
// 		}
// 	}

// }

func TestPRF(t *testing.T) {
	K := []byte{0x2b, 0x7e, 0x15, 0x16, 0x28, 0xae, 0xd2, 0xa6, 0xab, 0xf7, 0x15, 0x88, 0x09, 0xcf, 0x4f, 0x3c}
	X := []uint64{1, 2, 1, 0, 0, 10, 10, 5, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 221, 213}

	var Xbitstring []byte
	for _, x := range X {
		Xbitstring = append(Xbitstring, algorithms.STRmRadix(x, 2, 8)...)
	}

	expected := []byte{195, 184, 41, 161, 232, 100, 43, 120, 204, 41, 148, 123, 59, 147, 219, 99}

	result, err := algorithms.PRF(K, Xbitstring)
	if err != nil {
		t.Errorf("PRF() error: %v", err)
		return
	}

	for i, resultByte := range result {
		t.Logf("Result[%d]: %d (expected: %d)", i, resultByte, expected[i])
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

// func TestNumStringToIntAndBack(t *testing.T) {
// 	plaintext := "0123456789"
// 	radix := 10
// 	n := algorithms.numStringToInt(plaintext, radix)
// 	result := algorithms.intToNumString(n, len(plaintext), radix)

// 	if result != plaintext {
// 		t.Errorf("expected %v, got %v", plaintext, result)
// 	}
// }

// func TestMod(t *testing.T) {
// 	if res := algorithms.mod(-3, 7); res != 4 {
// 		t.Errorf("expected mod(-3, 7) to be 4, got %d", res)
// 	}
// 	if res := algorithms.mod(13, 7); res != 6 {
// 		t.Errorf("expected mod(13, 7) to be 6, got %d", res)
// 	}
// }

// func TestFF1EncryptDecrypt(t *testing.T) {
// 	key, _ := hex.DecodeString("2B7E151628AED2A6ABF7158809CF4F3C")
// 	tweak := []byte{}

// 	plaintext := "0123456789"
// 	ciphertext, err := algorithms.FF1Encrypt(key, tweak, plaintext, 10)
// 	if err != nil {
// 		t.Fatalf("encryption failed: %v", err)
// 	}

// 	decryptedText, err := algorithms.FF1Decrypt(key, tweak, ciphertext, 10)
// 	if err != nil {
// 		t.Fatalf("decryption failed: %v", err)
// 	}

// 	if decryptedText != plaintext {
// 		t.Errorf("expected %v, got %v", plaintext, decryptedText)
// 	}
// }

// func TestRepresentCharacters(t *testing.T) {
// 	result, err := algorithms.representCharacters("01234", 10)
// 	if err != nil {
// 		t.Fatalf("representation failed: %v", err)
// 	}
// 	expected := []int{0, 1, 2, 3, 4}
// 	for i, v := range result {
// 		if v != expected[i] {
// 			t.Errorf("expected %v at index %d, got %v", expected[i], i, v)
// 		}
// 	}
// }
