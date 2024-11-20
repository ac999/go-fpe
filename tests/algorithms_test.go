// tests/algorithms_test.go
package tests

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
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

func TestBigFloorDiv(t *testing.T) {
	// Helper to easily convert integers to *big.Int
	toBigInt := func(x int64) *big.Int {
		return big.NewInt(x)
	}

	// Test case 1: General case, 7 / 5
	result := algorithms.BigFloorDiv(toBigInt(7), toBigInt(5))
	expected := toBigInt(1) // 7 / 5 = 1 (floor division)
	if result.Cmp(expected) != 0 {
		t.Errorf("expected %d, got %d", expected, result)
	}

	// Test case 2: Exact division, 3 / 3
	result = algorithms.BigFloorDiv(toBigInt(3), toBigInt(3))
	expected = toBigInt(1) // 3 / 3 = 1
	if result.Cmp(expected) != 0 {
		t.Errorf("expected %d, got %d", expected, result)
	}

	// Test case 3: Negative numbers, -7 / 5
	result = algorithms.BigFloorDiv(toBigInt(-7), toBigInt(5))
	expected = toBigInt(-2) // -7 / 5 = -2 (floor division)
	if result.Cmp(expected) != 0 {
		t.Errorf("expected %d, got %d", expected, result)
	}

	// Test case 4: Negative numbers, 7 / -5
	result = algorithms.BigFloorDiv(toBigInt(7), toBigInt(-5))
	expected = toBigInt(-2) // 7 / -5 = -2 (floor division)
	if result.Cmp(expected) != 0 {
		t.Errorf("expected %d, got %d", expected, result)
	}

	// Test case 5: Large numbers
	result = algorithms.BigFloorDiv(toBigInt(1234567890123456789), toBigInt(1234567890))
	expected = toBigInt(1000000000) // Expected result for large numbers
	if result.Cmp(expected) != 0 {
		t.Errorf("expected %d, got %d", expected, result)
	}

	// Test case 6: Buffer overflow-like check with very large numbers
	x := new(big.Int).Exp(big.NewInt(2), big.NewInt(500), nil) // 2^500
	y := big.NewInt(2)
	result = algorithms.BigFloorDiv(x, y)
	expected = new(big.Int).Div(x, y)
	if result.Cmp(expected) != 0 {
		t.Errorf("expected %d, got %d", expected, result)
	}

	// Test case 7: Edge case with zero divisor (expecting panic)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic on division by zero")
		}
	}()
	_ = algorithms.BigFloorDiv(toBigInt(10), toBigInt(0)) // should panic
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

func TestBigCeilingDiv(t *testing.T) {
	toBigInt := func(x int64) *big.Int {
		return big.NewInt(x)
	}

	tests := []struct {
		x, y, expected int64
	}{
		// General positive division
		{7, 5, 2},
		{3, 3, 1},
		{1234567890123456789, 1234567890, 1000000001},

		// Negative dividend
		{-7, 5, -1},
		{-15, 5, -3},
		{-7, -5, 2},
		{-15, -5, 3},

		// Mixed signs
		{7, -5, -1},
		{-7, 5, -1},
		{15, -5, -3},
		{-15, 5, -3},

		// Zero dividend
		{0, 1, 0},
	}

	for _, tt := range tests {
		x := toBigInt(tt.x)
		y := toBigInt(tt.y)
		expected := toBigInt(tt.expected)

		result := algorithms.BigCeilingDiv(x, y)
		if result.Cmp(expected) != 0 {
			t.Errorf("BigCeilingDiv(%d, %d) = %d, want %d", tt.x, tt.y, result, expected)
		}
	}

	// Edge case: Buffer overflow-like check
	x := new(big.Int).Exp(big.NewInt(2), big.NewInt(500), nil) // 2^500
	y := big.NewInt(2)
	result := algorithms.BigCeilingDiv(x, y)
	expected := new(big.Int).Div(x, y)
	if result.Cmp(expected) != 0 {
		t.Errorf("expected %d, got %d", expected, result)
	}

	// Edge case: Division by zero
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic on division by zero")
		}
	}()
	_ = algorithms.BigCeilingDiv(toBigInt(10), toBigInt(0)) // Should panic
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
	result = algorithms.ModInt(13, 7)
	if result != 6 {
		t.Errorf("expected %d, got %d", 6, result)
	}
	result = algorithms.ModInt(-6, 16)
	if result != 10 {
		t.Errorf("expected %d, got %d", 10, result)
	}
}

func TestBigMod(t *testing.T) {
	// Helper to convert int64 to *big.Int
	toBigInt := func(x int64) *big.Int {
		return big.NewInt(x)
	}

	// Test cases for signed modulus with *big.Int
	tests := []struct {
		x, m, expected int64
	}{
		{-3, 7, 4},   // -3 % 7 = 4
		{-4, 16, 12}, // -4 % 16 = 12
		{10, 3, 1},   // 10 % 3 = 1
		{-15, 6, 3},  // -15 % 6 = 3
		{13, 7, 6},   // 13 % 7 = 6
		{100, 6, 4},  // 100 % 6 = 4
		{7, 5, 2},    // 7 % 5 = 2
		{1234567890123456789, 1234567890, 123456789}, // Large case
	}

	for _, tt := range tests {
		x := toBigInt(tt.x)
		m := toBigInt(tt.m)
		expected := toBigInt(tt.expected)

		result := algorithms.BigMod(x, m)
		if result.Cmp(expected) != 0 {
			t.Errorf("BigMod(%d, %d) = %d, want %d", tt.x, tt.m, result, expected)
		}
	}
}

func TestBigPower(t *testing.T) {
	// Helper to convert int64 to *big.Int
	toBigInt := func(x int64) *big.Int {
		return big.NewInt(x)
	}

	tests := []struct {
		x, y int64
		want string
	}{
		// Simple cases
		{2, 0, "1"},
		{2, 1, "2"},
		{2, 2, "4"},
		{3, 3, "27"},
		{5, 5, "3125"},
		{10, 10, "10000000000"},
		{7, 8, "5764801"},
		// Large number case
		{12345, 6, "3539537889086624823140625"},
		{9876, 5, "93951865167752549376"},
	}

	for _, tt := range tests {
		t.Run("Testing BigPower", func(t *testing.T) {
			got := algorithms.BigPower(toBigInt(tt.x), toBigInt(tt.y)).String()
			if got != tt.want {
				t.Errorf("BigPower(%d, %d) = %v; want %v", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

func TestByteLen(t *testing.T) {
	x := []byte{1, 0, 1, 1, 1, 0, 0, 1, 1, 0, 1, 0, 1, 1, 0, 0}
	result := algorithms.ByteLen(x)
	if result != 2 {
		t.Errorf("expected %d, got %d", 2, result)
	}
}

func TestBigByteLen(t *testing.T) {
	tests := []struct {
		x    []byte
		want string
	}{
		{[]byte{1, 2, 3, 4}, "0"},                                            // 4 bytes => 4 / 8 = 0 8-bit units
		{[]byte{1, 2, 3, 4, 5, 6, 7, 8}, "1"},                                // 8 bytes => 8 / 8 = 1 8-bit unit
		{[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, "1"},                         // 10 bytes => 10 / 8 = 1 8-bit unit
		{[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, "2"}, // 16 bytes => 16 / 8 = 2 8-bit units
	}

	for _, tt := range tests {
		t.Run("Testing ByteLen", func(t *testing.T) {
			got := algorithms.BigByteLen(tt.x).String()
			if got != tt.want {
				t.Errorf("ByteLen(%v) = %v; want %v", tt.x, got, tt.want)
			}
		})
	}
}

func TestNUM(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected uint64
	}{
		{
			name:     "Single byte",
			input:    []byte{128}, // Binary: 10000000
			expected: 128,
		},
		{
			name:     "Multiple bytes",
			input:    []byte{195, 184, 41, 161, 232, 100, 43, 120}, // Hex: c3b829a1e8642b78
			expected: 14103068008476060536,
		},
		{
			name:     "All zero bytes",
			input:    []byte{0, 0, 0, 0},
			expected: 0,
		},
		{
			name:     "All one bytes",
			input:    []byte{255, 255}, // Binary: 11111111 11111111
			expected: 65535,
		},
		{
			name:     "Sequential bytes",
			input:    []byte{1, 2, 3, 4}, // Hex: 0x01020304
			expected: 16909060,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := algorithms.NUM(tc.input)
			if result != tc.expected {
				t.Errorf("NUM(%v) = %d, expected %d", tc.input, result, tc.expected)
			}
		})
	}
}

func TestBigNUM(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected *big.Int
	}{
		{
			name:     "Single byte",
			input:    []byte{128}, // Binary: 10000000
			expected: big.NewInt(128),
		},
		{
			name:     "Multiple bytes",
			input:    []byte{195, 184, 41, 161, 232, 100, 43, 120}, // Hex: c3b829a1e8642b78
			expected: new(big.Int).SetBytes([]byte{195, 184, 41, 161, 232, 100, 43, 120}),
		},
		{
			name:     "All zero bytes",
			input:    []byte{0, 0, 0, 0},
			expected: big.NewInt(0),
		},
		{
			name:     "All one bytes",
			input:    []byte{255, 255}, // Binary: 11111111 11111111
			expected: big.NewInt(65535),
		},
		{
			name:     "Sequential bytes",
			input:    []byte{1, 2, 3, 4}, // Hex: 0x01020304
			expected: big.NewInt(16909060),
		},

		// Larger test cases with more bytes
		{
			name:     "Large number 1",
			input:    []byte{255, 255, 255, 255, 255, 255, 255, 255},                        // Large hex number
			expected: new(big.Int).SetBytes([]byte{255, 255, 255, 255, 255, 255, 255, 255}), // 0xFFFFFFFFFFFFFFFF
		},
		{
			name:     "Large number 2",
			input:    []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, // Even larger hex number
			expected: new(big.Int).SetBytes([]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}),
		},
		{
			name:     "Very large number",
			input:    []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, // Larger hex number
			expected: new(big.Int).SetBytes([]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}),
		},
		{
			name:     "Large number with mixed bytes",
			input:    []byte{1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0}, // Mixed byte sequence
			expected: new(big.Int).SetBytes([]byte{1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0}),
		},
		{
			name:     "Huge number",
			input:    []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			expected: new(big.Int).SetBytes([]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := algorithms.BigNUM(tc.input)
			// Compare the result with the expected value using Cmp (for big.Int comparison)
			if result.Cmp(tc.expected) != 0 {
				t.Errorf("NUM(%v) = %s, expected %s", tc.input, result, tc.expected)
			}
		})
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
		{[]byte{1, 2, 3}, 10, 123},
		{[]byte{0, 0, 1, 2, 3}, 10, 123},
		{[]byte{7}, 10, 7},
		{[]byte{}, 10, 0},
		{[]byte{3, 2, 1}, 4, 57},
		{[]byte{1, 0, 1, 1}, 2, 11},
		{[]byte{15, 15, 15}, 16, 4095}, // 0xFFF
		{[]byte{3, 5, 7}, 100, 30507},
		{[]byte{255, 255, 255}, 256, 16777215}, // 0xFFFFFF
	}

	for _, test := range tests {
		result := algorithms.NUMradix(test.X, test.radix)
		if result != test.expected {
			t.Errorf("For input X=%v and radix=%d, expected %d, but got %d", test.X, test.radix, test.expected, result)
		}
	}
}

func TestBigNUMradix(t *testing.T) {
	tests := []struct {
		X        []byte
		radix    uint64
		expected *big.Int
	}{
		{
			X:        []byte{0, 0, 0, 1, 1, 0, 1, 0},
			radix:    5,
			expected: big.NewInt(755),
		},
		{
			X:        []byte{1, 0, 1, 1},
			radix:    2,
			expected: big.NewInt(11),
		},
		{
			X:        []byte{3, 2, 1},
			radix:    4,
			expected: big.NewInt(57),
		},
		{
			X:        []byte{1, 2, 3},
			radix:    10,
			expected: big.NewInt(123),
		},
		{
			X:        []byte{0, 0, 1, 2, 3},
			radix:    10,
			expected: big.NewInt(123),
		},
		{
			X:        []byte{7},
			radix:    10,
			expected: big.NewInt(7),
		},
		{
			X:        []byte{},
			radix:    10,
			expected: big.NewInt(0),
		},
		{
			X:        []byte{15, 15, 15},
			radix:    16,
			expected: big.NewInt(4095), // 0xFFF
		},
		{
			X:        []byte{3, 5, 7},
			radix:    100,
			expected: big.NewInt(30507),
		},
		{
			X:        []byte{255, 255, 255},
			radix:    256,
			expected: big.NewInt(16777215), // 0xFFFFFF
		},
		// Test larger numbers
		{
			X:        []byte{255, 255, 255, 255, 255, 255, 255, 255},
			radix:    256,
			expected: new(big.Int).SetBytes([]byte{255, 255, 255, 255, 255, 255, 255, 255}),
		},
		{
			X:        []byte{1, 2, 3, 4, 5, 6, 7, 8},
			radix:    10,
			expected: big.NewInt(12345678),
		},
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {
			result := algorithms.BigNUMradix(tc.X, tc.radix)
			// Compare the result with the expected value using Cmp (for big.Int comparison)
			if result.Cmp(tc.expected) != 0 {
				t.Errorf("BigNUMradix(%v, %d) = %s, expected %s", tc.X, tc.radix, result, tc.expected)
			}
		})
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
		{0, 256, 10, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
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

func TestBigSTRmRadix(t *testing.T) {
	tests := []struct {
		x        *big.Int
		radix    uint64
		m        int64
		expected []byte
	}{
		{
			x:        big.NewInt(559),
			radix:    12,
			m:        4,
			expected: []byte{0, 3, 10, 7},
		},
		{
			x:        big.NewInt(1),
			radix:    2,
			m:        8,
			expected: []byte{0, 0, 0, 0, 0, 0, 0, 1},
		},
		{
			x:        big.NewInt(255),
			radix:    16,
			m:        2,
			expected: []byte{15, 15},
		},
		{
			x:        big.NewInt(1024),
			radix:    2,
			m:        11,
			expected: []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			x:        big.NewInt(123),
			radix:    10,
			m:        3,
			expected: []byte{1, 2, 3},
		},
		{
			x:        big.NewInt(0),
			radix:    256,
			m:        10,
			expected: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			result := algorithms.BigSTRmRadix(test.x, test.radix, test.m)
			if len(result) != len(test.expected) {
				t.Errorf("BigSTRmRadix() length = %v, expected %v", len(result), len(test.expected))
			} else {
				for i := range result {
					if result[i] != test.expected[i] {
						t.Errorf("BigSTRmRadix() result at index %v = %v, expected %v", i, result[i], test.expected[i])
					}
				}
			}
		})
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
			// Example 1 from FF1samples.pd
			name:         "FF1-AES128-Sample1",
			keyHex:       "2B7E151628AED2A6ABF7158809CF4F3C",
			tweak:        []byte{},
			plaintextStr: "0123456789",
			expectedEnc:  "2433477484",
			radix:        10,
		},
		{
			// In example 2 of FF1samples.pdf the tweak is wrong. Adjusted in this test.
			name:         "FF1-AES128-Sample2",
			keyHex:       "2B7E151628AED2A6ABF7158809CF4F3C",
			tweak:        []byte{57, 56, 55, 54, 53, 52, 51, 50, 49, 48},
			plaintextStr: "0123456789",
			expectedEnc:  "6124200773",
			radix:        10,
		}, {
			// Example 6 from FF1samples.pdf
			name:         "FF1-AES128-Sample6",
			keyHex:       "2B7E151628AED2A6ABF7158809CF4F3CEF4359D8D580AA4F",
			tweak:        []byte{55, 55, 55, 55, 112, 113, 114, 115, 55, 55, 55},
			plaintextStr: "0123456789abcdefghi",
			expectedEnc:  "xbj3kv35jrawxv32ysr",
			radix:        36,
		},
		{
			// Credit Card Number example
			name:         "FF1-AES128-CreditCardNumber",
			keyHex:       "637265646974636172646E756D626572",
			tweak:        []byte{},
			plaintextStr: "4557534296728436",
			expectedEnc:  "5055966384254029",
			radix:        10,
		},
		{
			// Vehicle Number Plate example
			name:         "FF1-AES128-VehicleNumberPlate",
			keyHex:       "706F6C6974696174726563656E696E6F",
			tweak:        []byte{},
			plaintextStr: "is01mai",
			expectedEnc:  "y8u7ykm",
			radix:        36,
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
			var alphabet string
			switch tc.radix {
			case 10:
				alphabet = alphabets["base10"]
			case 26:
				alphabet = alphabets["base26"]
			case 36:
				alphabet = alphabets["base36"]

			}
			// Step 2: Convert plaintext and expected ciphertext to numeral slices
			plaintext, err := algorithms.StringToNumeralSlice(tc.plaintextStr, alphabet)
			if err != nil {
				t.Fatalf("Failed to convert plaintext string to numeral slice: %v", err)
			}

			expectedEnc, err := algorithms.StringToNumeralSlice(tc.expectedEnc, alphabet)
			if err != nil {
				t.Fatalf("Failed to convert expected ciphertext string to numeral slice: %v", err)
			}

			// Debugging: Print the converted plaintext
			t.Logf("Converted plaintext = %v", plaintext)

			// Step 3: Perform encryption
			fmt.Printf("##################\n%s\n##################\n", tc.name)
			fmt.Printf("Encrypt():\n")
			ciphertext, err := algorithms.Encrypt(key, tc.tweak, plaintext, tc.radix)
			if err != nil {
				t.Fatalf("Encryption failed: %v", err)
			}

			// Debugging: Print encryption result
			t.Logf("Ciphertext = %v (expected: %v)", ciphertext, expectedEnc)

			expectedCiphertextStr, err := algorithms.NumeralSliceToString(ciphertext, alphabet)
			if err != nil {
				t.Fatalf("Failed to convert expected ciphertext string to numeral slice: %v", err)
			}
			t.Logf("Decrypted plaintext as string = %v", expectedCiphertextStr)

			// Step 4: Validate encryption result
			if !reflect.DeepEqual(ciphertext, expectedEnc) {
				t.Errorf("Encrypt() result = %v, expected %v", ciphertext, expectedEnc)
			}

			// Step 5: Perform decryption
			fmt.Printf("Decrypt():\n")
			decryptedText, err := algorithms.Decrypt(key, tc.tweak, ciphertext, tc.radix)
			if err != nil {
				t.Fatalf("Decryption failed: %v", err)
			}

			// Debugging: Print decryption result
			t.Logf("Decrypted plaintext = %v (expected: %v)", decryptedText, plaintext)

			expectedDecStr, err := algorithms.NumeralSliceToString(decryptedText, alphabet)
			if err != nil {
				t.Fatalf("Failed to convert expected ciphertext string to numeral slice: %v", err)
			}
			t.Logf("Decrypted plaintext as string = %v", expectedDecStr)

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
