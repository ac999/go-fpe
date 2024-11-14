// tests/algorithms_test.go
package tests

// import (
// 	"encoding/hex"
// 	"testing"

// 	"github.com/ac999/go-fpe/algorithms"
// )

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
