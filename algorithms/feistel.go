// algorithms/feistel.go
package algorithms

import (
	"encoding/binary"
	"fmt"
)

func FeistelRound(key []byte, tweak []byte, B []int, roundIndex, radix int) int {
	// Prepare the tweak + round number data
	roundData := make([]byte, len(tweak)+1)
	copy(roundData, tweak)
	roundData[len(tweak)] = byte(roundIndex)

	// Convert B into a numeral string and encode it to bytes
	BNum := NUM(B, radix)
	BBytes := encodeToBytes(BNum)

	// Concatenate the tweak data and numeral bytes
	prfInput := append(roundData, BBytes...)

	// Compute PRF using the key
	prfResult, err := PRF(key, prfInput)
	if err != nil {
		panic(fmt.Sprintf("PRF error: %v", err))
	}

	// Convert PRF result to a number
	prfNumber := binary.BigEndian.Uint64(prfResult[:8]) // Use only the first 8 bytes

	// Return the resulting value
	return int(prfNumber % uint64(radix))
}

// Helper to encode an integer to a byte slice
func encodeToBytes(x int) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(x))
	return buf
}
