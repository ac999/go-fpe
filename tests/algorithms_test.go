// tests/algorithms_test.go
package tests

import (
	"crypto/aes"
	"testing"

	"github.com/ac999/go-fpe/algorithms"
)

func TestFF1EncryptDecrypt(t *testing.T) {
	key := make([]byte, 16) // 128-bit key (zeroed for simplicity)
	tweak := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9} // Input numeral string
	radix := 10

	// Encrypt
	encrypted, err := algorithms.FF1Encrypt(key, tweak, input, radix, 2, 32)
	if err != nil {
		t.Fatalf("FF1 Encrypt failed: %v", err)
	}

	// Decrypt
	decrypted, err := algorithms.FF1Decrypt(key, tweak, encrypted, radix, 2, 32)
	if err != nil {
		t.Fatalf("FF1 Decrypt failed: %v", err)
	}

	if !equalSlices(input, decrypted) {
		t.Errorf("FF1 Decryption failed: got %v, want %v", decrypted, input)
	}
}

func TestFF1EncryptDecryptEdgeCases(t *testing.T) {
	key := make([]byte, 16) // 128-bit key (zeroed for simplicity)

	// Edge Case 1: Minimum radix
	radix := 2
	tweak := []byte{}             // Empty tweak
	input := []int{1, 0, 1, 1, 0} // Binary input
	encrypted, err := algorithms.FF1Encrypt(key, tweak, input, radix, 2, 32)
	if err != nil {
		t.Fatalf("FF1 Encrypt failed for minimum radix: %v", err)
	}
	decrypted, err := algorithms.FF1Decrypt(key, tweak, encrypted, radix, 2, 32)
	if err != nil {
		t.Fatalf("FF1 Decrypt failed for minimum radix: %v", err)
	}
	if !equalSlices(input, decrypted) {
		t.Errorf("FF1 Decrypt failed for minimum radix: got %v, want %v", decrypted, input)
	}

	// Edge Case 2: Maximum radix (2^16)
	radix = 65536
	tweak = []byte{0x01, 0x02, 0x03, 0x04} // 4-byte tweak
	input = []int{500, 1000, 1500, 2000}   // Input with high numeral values
	encrypted, err = algorithms.FF1Encrypt(key, tweak, input, radix, 2, 32)
	if err != nil {
		t.Fatalf("FF1 Encrypt failed for maximum radix: %v", err)
	}
	decrypted, err = algorithms.FF1Decrypt(key, tweak, encrypted, radix, 2, 32)
	if err != nil {
		t.Fatalf("FF1 Decrypt failed for maximum radix: %v", err)
	}
	if !equalSlices(input, decrypted) {
		t.Errorf("FF1 Decrypt failed for maximum radix: got %v, want %v", decrypted, input)
	}

	// Edge Case 3: Odd-length numeral string
	radix = 10
	tweak = []byte{0x05, 0x06, 0x07, 0x08, 0x09} // 5-byte tweak
	input = []int{1, 2, 3, 4, 5}                 // Odd length
	encrypted, err = algorithms.FF1Encrypt(key, tweak, input, radix, 2, 32)
	if err != nil {
		t.Fatalf("FF1 Encrypt failed for odd-length numeral: %v", err)
	}
	decrypted, err = algorithms.FF1Decrypt(key, tweak, encrypted, radix, 2, 32)
	if err != nil {
		t.Fatalf("FF1 Decrypt failed for odd-length numeral: %v", err)
	}
	if !equalSlices(input, decrypted) {
		t.Errorf("FF1 Decrypt failed for odd-length numeral: got %v, want %v", decrypted, input)
	}

	// Edge Case 4: Large numeral input
	tweak = []byte{0x0A, 0x0B, 0x0C, 0x0D}
	input = []int{999, 1234, 5678, 9876, 5432, 1000} // Large values
	encrypted, err = algorithms.FF1Encrypt(key, tweak, input, radix, 2, 32)
	if err != nil {
		t.Fatalf("FF1 Encrypt failed for large numeral values: %v", err)
	}
	decrypted, err = algorithms.FF1Decrypt(key, tweak, encrypted, radix, 2, 32)
	if err != nil {
		t.Fatalf("FF1 Decrypt failed for large numeral values: %v", err)
	}
	if !equalSlices(input, decrypted) {
		t.Errorf("FF1 Decrypt failed for large numeral values: got %v, want %v", decrypted, input)
	}
}

func TestFF3EncryptDecrypt(t *testing.T) {
	key := make([]byte, 16) // 128-bit key (zeroed for simplicity)
	tweak := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9} // Input numeral string
	radix := 10

	// Encrypt
	encrypted, err := algorithms.FF3Encrypt(key, tweak, input, radix)
	if err != nil {
		t.Fatalf("FF3 Encrypt failed: %v", err)
	}

	// Decrypt
	decrypted, err := algorithms.FF3Decrypt(key, tweak, encrypted, radix)
	if err != nil {
		t.Fatalf("FF3 Decrypt failed: %v", err)
	}

	if !equalSlices(input, decrypted) {
		t.Errorf("FF3 Decryption failed: got %v, want %v", decrypted, input)
	}
}

func TestFF3EncryptDecryptEdgeCases(t *testing.T) {
	key := make([]byte, 16) // 128-bit key (zeroed for simplicity)

	// Edge Case 1: Empty tweak (should fail for FF3)
	radix := 10
	tweak := []byte{} // Empty tweak, should cause an error
	input := []int{3, 4, 5, 6, 7}
	_, err := algorithms.FF3Encrypt(key, tweak, input, radix)
	if err == nil {
		t.Fatalf("FF3 Encrypt should have failed for empty tweak, but it didn't.")
	}

	// Edge Case 2: Maximum 64-bit tweak for FF3
	tweak = []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77}
	input = []int{1, 2, 3, 4, 5, 6, 7, 8} // Even length
	encrypted, err := algorithms.FF3Encrypt(key, tweak, input, radix)
	if err != nil {
		t.Fatalf("FF3 Encrypt failed for maximum 64-bit tweak: %v", err)
	}
	decrypted, err := algorithms.FF3Decrypt(key, tweak, encrypted, radix)
	if err != nil {
		t.Fatalf("FF3 Decrypt failed for maximum 64-bit tweak: %v", err)
	}
	if !equalSlices(input, decrypted) {
		t.Errorf("FF3 Decrypt failed for maximum 64-bit tweak: got %v, want %v", decrypted, input)
	}

	// Edge Case 3: All zeros numeral string
	tweak = []byte{0x01, 0x02, 0x03, 0x04}
	input = []int{0, 0, 0, 0, 0} // All zeros
	encrypted, err = algorithms.FF3Encrypt(key, tweak, input, radix)
	if err != nil {
		t.Fatalf("FF3 Encrypt failed for all zeros numeral: %v", err)
	}
	decrypted, err = algorithms.FF3Decrypt(key, tweak, encrypted, radix)
	if err != nil {
		t.Fatalf("FF3 Decrypt failed for all zeros numeral: %v", err)
	}
	if !equalSlices(input, decrypted) {
		t.Errorf("FF3 Decrypt failed for all zeros numeral: got %v, want %v", decrypted, input)
	}
}

// Helper function to compare two slices
func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestPRF(t *testing.T) {
	key := make([]byte, 16)
	data := []byte{0x00, 0x01, 0x02, 0x03}

	result, err := algorithms.PRF(key, data)
	if err != nil {
		t.Fatalf("PRF failed: %v", err)
	}

	if len(result) != aes.BlockSize {
		t.Errorf("PRF result length is incorrect: %d", len(result))
	}
}

func TestPRFEdgeCases(t *testing.T) {
	key := make([]byte, 16) // 128-bit key (zeroed for simplicity)

	// Edge Case 1: Short input data
	data := []byte{0x01}
	result, err := algorithms.PRF(key, data)
	if err != nil {
		t.Fatalf("PRF failed for short input data: %v", err)
	}
	if len(result) != 16 {
		t.Errorf("PRF result length is incorrect for short input data: got %d, want 16", len(result))
	}

	// Edge Case 2: Large input data
	data = make([]byte, 1000) // Large input data
	_, err = algorithms.PRF(key, data)
	if err != nil {
		t.Fatalf("PRF failed for large input data: %v", err)
	}

	// Edge Case 3: Maximum length key for AES-256
	key = make([]byte, 32) // 256-bit key
	data = []byte{0x10, 0x20, 0x30, 0x40}
	result, err = algorithms.PRF(key, data)
	if err != nil {
		t.Fatalf("PRF failed for maximum length key: %v", err)
	}
	if len(result) != 16 {
		t.Errorf("PRF result length is incorrect for maximum length key: got %d, want 16", len(result))
	}
}
