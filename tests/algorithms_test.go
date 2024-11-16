// tests/algorithms_test.go
package tests

import (
	"bytes"
	"encoding/hex"
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
		expected  []byte
		shouldErr bool
	}{
		{"hello", alphabets["base26"], []byte{7, 4, 11, 11, 14}, false},
		{"01234", alphabets["base10"], []byte{0, 1, 2, 3, 4}, false},
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

func TestModInt(t *testing.T) {
	result := algorithms.ModInt(-3, 7)
	if result != 4 {
		t.Errorf("expected %d, got %d", 4, result)
	}
	result = algorithms.ModInt(-4, 16)
	if result != 12 {
		t.Errorf("expected %d, got %d", 12, result)
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
	tests := []struct {
		X        []byte
		radix    uint64
		expected uint64
	}{
		{[]byte{0, 0, 0, 1, 1, 0, 1, 0}, 5, 755},
		{[]byte{1, 0, 1, 1}, 2, 11},
		{[]byte{3, 2, 1}, 4, 57},
		{[]byte{1, 2, 3}, 10, 123},
		{[]byte{0, 0, 1, 2, 3}, 10, 123},
		{[]byte{7}, 10, 7},
		{[]byte{}, 10, 0},
	}

	for _, test := range tests {
		result := algorithms.NUMradix(test.X, test.radix)
		if result != test.expected {
			t.Errorf("For input X=%v and radix=%d, expected %d, but got %d", test.X, test.radix, test.expected, result)
		}
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
	tests := []struct {
		x        uint64
		radix    uint64
		m        int64
		expected []byte
	}{
		{559, 12, 4, []byte{0, 3, 10, 7}},
		{1, 2, 8, []byte{0, 0, 0, 0, 0, 0, 0, 1}},
		{255, 16, 2, []byte{15, 15}},
		{1024, 2, 11, []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{123, 10, 3, []byte{1, 2, 3}},
	}

	for _, test := range tests {
		result := algorithms.STRmRadix(test.x, test.radix, test.m)
		if len(result) != len(test.expected) {
			t.Errorf("STRmRadix() length = %v, expected %v", len(result), len(test.expected))
		} else {
			for i := range result {
				if result[i] != test.expected[i] {
					t.Errorf("STRmRadix() result at index %v = %v, expected %v", i, result[i], test.expected[i])
				}
			}
		}
	}
}

// TestAesEncrypt validates AES encryption against FIPS 197 test vectors.
func TestAesEncrypt(t *testing.T) {
	tests := []struct {
		name              string
		keyHex            string
		plaintextHex      string
		expectedCipherHex string
	}{
		{
			name:              "FIPS 197 Example: AES-128",
			keyHex:            "000102030405060708090A0B0C0D0E0F",
			plaintextHex:      "00112233445566778899AABBCCDDEEFF",
			expectedCipherHex: "69C4E0D86A7B0430D8CDB78070B4C55A",
		},
		{
			name:              "All Zeros Key and Plaintext",
			keyHex:            "00000000000000000000000000000000",
			plaintextHex:      "00000000000000000000000000000000",
			expectedCipherHex: "66E94BD4EF8A2C3B884CFA59CA342B2E",
		},
		{
			name:              "Incrementing Bytes Key and Plaintext",
			keyHex:            "101112131415161718191A1B1C1D1E1F",
			plaintextHex:      "202122232425262728292A2B2C2D2E2F",
			expectedCipherHex: "D31DD57E62812CDDABD1CCAA3C47979B",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			key, err := hex.DecodeString(test.keyHex)
			if err != nil {
				t.Fatalf("Invalid key hex: %s", test.keyHex)
			}

			plaintext, err := hex.DecodeString(test.plaintextHex)
			if err != nil {
				t.Fatalf("Invalid plaintext hex: %s", test.plaintextHex)
			}

			expectedCipher, err := hex.DecodeString(test.expectedCipherHex)
			if err != nil {
				t.Fatalf("Invalid expected cipher hex: %s", test.expectedCipherHex)
			}

			cipher := algorithms.AesEncrypt(plaintext, key)

			if !bytes.Equal(cipher, expectedCipher) {
				t.Errorf("Test %s failed.\nGot:      %X\nExpected: %X", test.name, cipher, expectedCipher)
			}
		})
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
	X := []byte{1, 2, 1, 0, 0, 10, 10, 5, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 221, 213}

	expected := []byte{195, 184, 41, 161, 232, 100, 43, 120, 204, 41, 148, 123, 59, 147, 219, 99}

	result, err := algorithms.PRF(K, X)
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

func TestFF1EncryptDecrypt(t *testing.T) {
	// Test cases from FF1samples.pdf
	testCases := []struct {
		name         string
		keyHex       string
		tweak        []byte
		plaintextStr string
		expectedEnc  string
		radix        uint64
	}{
		{
			name:         "FF1-AES128-Sample1",
			keyHex:       "2B7E151628AED2A6ABF7158809CF4F3C",
			tweak:        []byte{},
			plaintextStr: "0123456789",
			expectedEnc:  "2433477484",
			radix:        10,
		},
		{
			name:         "FF1-AES128-Sample2",
			keyHex:       "2B7E151628AED2A6ABF7158809CF4F3C",
			tweak:        []byte{},
			plaintextStr: "9876543210",
			expectedEnc:  "5868123250",
			radix:        10,
		},
		// Add more test cases here from FF1samples.pdf
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Step 1: Parse the key
			key, err := hex.DecodeString(tc.keyHex)
			if err != nil {
				t.Fatalf("Failed to decode key: %v", err)
			}

			// Step 2: Convert plaintext and expected ciphertext to numeral slices
			plaintext, err := algorithms.StringToNumeralSlice(tc.plaintextStr, alphabets["base10"])
			if err != nil {
				t.Fatalf("Failed to convert plaintext string to numeral slice: %v", err)
			}

			expectedEnc, err := algorithms.StringToNumeralSlice(tc.expectedEnc, alphabets["base10"])
			if err != nil {
				t.Fatalf("Failed to convert expected ciphertext string to numeral slice: %v", err)
			}

			// Debugging: Print the converted plaintext
			t.Logf("Converted plaintext = %v", plaintext)

			// Step 3: Perform encryption
			ciphertext, err := algorithms.Encrypt(key, tc.tweak, plaintext, tc.radix)
			if err != nil {
				t.Fatalf("Encryption failed: %v", err)
			}

			// Debugging: Print encryption result
			t.Logf("Ciphertext = %v (expected: %v)", ciphertext, expectedEnc)

			// Step 4: Validate encryption result
			if !reflect.DeepEqual(ciphertext, expectedEnc) {
				t.Errorf("Encrypt() result = %v, expected %v", ciphertext, expectedEnc)
			}

			// Step 5: Perform decryption
			decryptedText, err := algorithms.Decrypt(key, tc.tweak, ciphertext, tc.radix)
			if err != nil {
				t.Fatalf("Decryption failed: %v", err)
			}

			// Debugging: Print decryption result
			t.Logf("Decrypted plaintext = %v (expected: %v)", decryptedText, plaintext)

			// Step 6: Validate decryption result
			if !reflect.DeepEqual(decryptedText, plaintext) {
				t.Errorf("Decrypt() result = %v, expected %v", decryptedText, plaintext)
			}

			// Step 7: Ensure encryption-decryption cycle is consistent
			if len(decryptedText) != len(plaintext) {
				t.Errorf("Decrypt() length = %v, expected %v", len(decryptedText), len(plaintext))
			} else {
				for i := range decryptedText {
					if decryptedText[i] != plaintext[i] {
						t.Errorf("Decrypt() mismatch at index %v: got %v, expected %v", i, decryptedText[i], plaintext[i])
					}
				}
			}
		})
	}
}

// func TestFF1EncryptDecrypt(t *testing.T) {
// 	key, _ := hex.DecodeString("2B7E151628AED2A6ABF7158809CF4F3C")
// 	tweak := []byte{}

// 	plaintext, err := algorithms.StringToNumeralSlice("0123456789", alphabets["base10"])
// 	expected_enc, err := algorithms.StringToNumeralSlice("2433477484", alphabets["base10"])
// 	if err != nil {
// 		t.Fatalf("conversion of plaintext failed: %v", err)
// 	}

// 	// Debugging
// 	fmt.Println("Converted plaintext = ")
// 	fmt.Println(plaintext)

// 	ciphertext, err := algorithms.Encrypt(key, tweak, plaintext, 10)
// 	if err != nil {
// 		t.Fatalf("encryption failed: %v", err)
// 	}

// 	decryptedText, err := algorithms.Decrypt(key, tweak, ciphertext, 10)
// 	if err != nil {
// 		t.Fatalf("decryption failed: %v", err)
// 	}

// 	if len(decryptedText) != len(plaintext) {
// 		t.Errorf("Decrypt() length = %v, expected %v", len(decryptedText), len(plaintext))
// 	} else {
// 		for i := range len(decryptedText) {
// 			if decryptedText[i] != plaintext[i] {
// 				t.Errorf("Decrypt() result at index %v = %v, expected %v", i, decryptedText[i], plaintext[i])
// 			}
// 		}
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
