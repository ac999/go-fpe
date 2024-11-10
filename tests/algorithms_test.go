// tests/algorithms_test.go
package tests

import (
	"crypto/aes"
	"testing"

	"github.com/ac999/go-fpe/algorithms"
)

func TestNumRadix(t *testing.T) {
	x, err := algorithms.NumRadix("123", 10)
	if err != nil || x != 123 {
		t.Errorf("Expected 123, got %d", x)
	}

	x, err = algorithms.NumRadix("1A", 16)
	if err != nil || x != 26 {
		t.Errorf("Expected 26, got %d", x)
	}

	// Test case: conversion from nist
	x, err = algorithms.NumRadix("00011010", 5)
	if err != nil || x != 755 {
		t.Errorf("Expected 755, got %d", x)
	}

	_, err = algorithms.NumRadix("1A", 10)
	if err == nil {
		t.Errorf("Expected error for invalid input")
	}
}

func TestNumBits(t *testing.T) {
	x, err := algorithms.NumBits("10000000")
	if err != nil || x != 128 {
		t.Errorf("Expected 128, got %d", x)
	}

	_, err = algorithms.NumBits("11012")
	if err == nil {
		t.Errorf("Expected error for invalid input")
	}
}

func TestStrmRadix(t *testing.T) {
	// Test case: basic conversion
	result, err := algorithms.StrmRadix(26, 16, 2)
	if err != nil || result != "110" {
		t.Errorf("Expected 110, got %s", result)
	}

	// Test case: conversion from nist
	result, err = algorithms.StrmRadix(559, 12, 4)
	if err != nil || result != "03107" {
		t.Errorf("Expected 03107, got %s", result)
	}

	// Test case: radix 10, length 3
	result, err = algorithms.StrmRadix(123, 10, 3)
	if err != nil || result != "123" {
		t.Errorf("Expected 123, got %s", result)
	}

	// Test case: minimum value
	result, err = algorithms.StrmRadix(0, 10, 3)
	if err != nil || result != "000" {
		t.Errorf("Expected 000, got %s", result)
	}

	// Test case: x out of range
	_, err = algorithms.StrmRadix(1000, 10, 2)
	if err == nil {
		t.Errorf("Expected error for out-of-range value")
	}
}

func TestRev(t *testing.T) {
	result := algorithms.Rev("12345")
	if result != "54321" {
		t.Errorf("Expected 54321, got %s", result)
	}

	result = algorithms.Rev("ABC")
	if result != "CBA" {
		t.Errorf("Expected CBA, got %s", result)
	}
}

func TestRevB(t *testing.T) {
	input := []byte{0b10110010, 0b11001101}    // [179, 205]
	expected := []byte{0b10110011, 0b01001101} // bit-reversed and byte-reversed

	result := algorithms.RevB(input)

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected %08b, got %08b", v, result[i])
		}
	}
}

func TestPRF(t *testing.T) {
	key := make([]byte, 16) // 16-byte AES key
	X := make([]byte, 16*4)
	Y, err := algorithms.PRF(X, key)
	if err != nil {
		t.Fatalf("Failed to create AES cipher: %v", err)
	}
	t.Logf("PRF Output: %v", Y)
}

func TestMinimalFF1Encrypt(t *testing.T) {
	key := []byte("examplekey123456") // 16-byte AES key
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("Failed to create AES cipher: %v", err)
	}
	t.Logf("AES Cipher created successfully: %v", block)
}

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
